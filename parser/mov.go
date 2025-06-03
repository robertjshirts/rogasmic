package parser

import (
	"fmt"

	"github.com/robertjshirts/rogasmic/types"
	"github.com/robertjshirts/rogasmic/utils"
)

type InstructionMOV struct {
	Mnemonic     types.MnemonicType
	Condition    types.ConditionType
	DestRegister uint32
	Immediate    uint32
}

func (p *Parser) parseMOV() (types.Instruction, error) {
	// Mnemonic
	mnemonic := types.TokenToMnemonic[p.current().Type]
	category, ok := types.MnemonicToCategory[mnemonic]
	if !ok || category != types.MnemonicCategoryMOV {
		return nil, fmt.Errorf("wrong instruction type! expected MOV instruction, got %s", p.current().Literal)
	}

	// Condition
	condition, err := utils.ParseMOVSuffixes(p.current().Literal)
	if err != nil {
		return nil, fmt.Errorf("error parsing MOV suffixes: %w", err)
	}
	p.consume() // consume MOVT or MOVW token

	// Register
	if p.current().Type != types.TokenRegister {
		return nil, fmt.Errorf("expected register after MOV mnemonic, got %s at line %d, col %d", p.current().Literal, p.current().Line, p.current().Col)
	}
	reg, err := utils.ParseRegister(p.current().Literal)
	if err != nil {
		return nil, fmt.Errorf("error parsing register: %w", err)
	}
	p.consume() // consume register token

	// Comma
	if p.current().Type != types.TokenComma {
		return nil, fmt.Errorf("expected comma after register, got %s at line %d, col %d", p.current().Literal, p.current().Line, p.current().Col)
	}
	p.consume() // consume comma token

	// Immediate
	if p.current().Type != types.TokenImmediate {
		return nil, fmt.Errorf("expected immediate value after comma, got %s at line %d, col %d", p.current().Literal, p.current().Line, p.current().Col)
	}
	immediate, err := utils.ParseImmediate(p.current().Literal)
	if err != nil {
		return nil, fmt.Errorf("error parsing immediate value: %w", err)
	}
	p.consume() // consume immediate token

	instruction := &InstructionMOV{
		Mnemonic:     mnemonic,
		Condition:    condition,
		Immediate:    immediate,
		DestRegister: reg,
	}

	return instruction, nil
}

func (i *InstructionMOV) ToMachineCode(labels map[string]uint32) ([]byte, error) {
	var binary uint32
	binary |= types.ConditionToBits[i.Condition] << 28
	binary |= types.MnemonicToBits[i.Mnemonic] << 20
	binary |= (i.Immediate >> 12) & 0xF << 16 // Top 4 bits of immediate
	binary |= i.DestRegister << 12
	binary |= i.Immediate & 0xFFF // Bottom 12 bits of immediate

	return utils.BitsToBytes(binary), nil
}
