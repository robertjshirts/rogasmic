package parser

import (
	"fmt"

	"github.com/robertjshirts/rogasmic/types"
	"github.com/robertjshirts/rogasmic/utils"
)

type InstructionArithmetic struct {
	Mnemonic     types.MnemonicType
	Condition    types.ConditionType
	DestRegister uint32
	BaseRegister uint32
	Immediate    uint32
	SBit         uint32
}

func (p *Parser) parseArithmetic() (types.Instruction, error) {
	// Mnemonic
	mnemonic := types.TokenToMnemonic[p.current().Type]
	category, ok := types.MnemonicToCategory[mnemonic]
	if !ok || category != types.MnemonicCategoryArithmetic {
		return nil, fmt.Errorf("wrong instruction type! expected arithmetic mnemonic, got %s", p.current().Literal)
	}

	// Get condition code and S suffix
	condition, sBit, err := utils.ParseArithmeticSuffixes(p.current().Literal)
	if err != nil {
		return nil, fmt.Errorf("error parsing arithmetic suffixes: %w", err)
	}
	p.consume() // consume arithmetic mnemonic token

	// Destination Register
	if p.current().Type != types.TokenRegister {
		return nil, fmt.Errorf("expected register after arithmetic mnemonic, got %s at line %d, col %d", p.current().Literal, p.current().Line, p.current().Col)
	}
	destReg, err := utils.ParseRegister(p.current().Literal)
	if err != nil {
		return nil, fmt.Errorf("error parsing destination register: %w", err)
	}
	p.consume() // consume destination register token

	// Comma
	if p.current().Type != types.TokenComma {
		return nil, fmt.Errorf("expected comma after destination register, got %s at line %d, col %d", p.current().Literal, p.current().Line, p.current().Col)
	}
	p.consume() // consume comma token

	// Base register
	if p.current().Type != types.TokenRegister {
		return nil, fmt.Errorf("expected base register after comma, got %s at line %d, col %d", p.current().Literal, p.current().Line, p.current().Col)
	}
	baseReg, err := utils.ParseRegister(p.current().Literal)
	if err != nil {
		return nil, fmt.Errorf("error parsing base register: %w", err)
	}
	p.consume() // consume base register token

	// Comma
	if p.current().Type != types.TokenComma {
		return nil, fmt.Errorf("expected comma after base register, got %s at line %d, col %d", p.current().Literal, p.current().Line, p.current().Col)
	}
	p.consume() // consume comma token

	// Immediate
	if p.current().Type != types.TokenImmediate {
		return nil, fmt.Errorf("expected immediate value after base register, got %s at line %d, col %d", p.current().Literal, p.current().Line, p.current().Col)
	}
	immediate, err := utils.ParseImmediate(p.current().Literal)
	if err != nil {
		return nil, fmt.Errorf("error parsing immediate value: %w", err)
	}
	p.consume() // consume immediate token

	instruction := &InstructionArithmetic{
		Mnemonic:     mnemonic,
		Condition:    condition,
		DestRegister: baseReg,
		BaseRegister: destReg,
		Immediate:    immediate,
		SBit:         sBit,
	}

	return instruction, nil
}

func (i *InstructionArithmetic) ToMachineCode(labels map[string]uint32) ([]byte, error) {
	var binary uint32
	binary |= types.ConditionToBits[i.Condition] << 28 // Condition code
	binary |= 0 << 26                                  // Always 0 for arithmetic instructions
	binary |= 1 << 25                                  // I Bit, we always gonna use immediate values
	binary |= types.MnemonicToBits[i.Mnemonic] << 21   // Mnemonic bits
	binary |= i.SBit << 20                             // S Bit
	binary |= i.DestRegister << 16                     // Destination register bits
	binary |= i.BaseRegister << 12                     // Base register bits
	binary |= i.Immediate & 0xFFF                      // Immediate value bits (lower 12 bits, top 4 are used for something called rotate)

	return utils.BitsToBytes(binary), nil
}
