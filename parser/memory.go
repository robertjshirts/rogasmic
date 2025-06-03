package parser

import (
	"fmt"

	"github.com/robertjshirts/rogasmic/types"
	"github.com/robertjshirts/rogasmic/utils"
)

type InstructionMemory struct {
	Mnemonic     types.MnemonicType
	Condition    types.ConditionType
	DestRegister uint32
	BaseRegister uint32
	IBit         uint32
	PBit         uint32
	UBit         uint32
	BBit         uint32
	WBit         uint32
}

func (p *Parser) parseMemory() (types.Instruction, error) {
	// Mnemonic
	mnemonic := types.TokenToMnemonic[p.current().Type]
	category, ok := types.MnemonicToCategory[mnemonic]
	if !ok || category != types.MnemonicCategoryLoadStore {
		return nil, fmt.Errorf("wrong instruction type! expected LDR or STR, got %s", p.current().Literal)
	}

	// Get condition code and suffixes
	condition, iBit, pBit, uBit, bBit, wBit, err := utils.ParseMemorySuffixes(p.current().Literal)
	if err != nil {
		return nil, fmt.Errorf("error parsing memory suffixes: %w", err)
	}
	p.consume() // consume LDR or STR token

	// Destination Register
	if p.current().Type != types.TokenRegister {
		return nil, fmt.Errorf("expected register after LDR/STR mnemonic, got %s at line %d, col %d", p.current().Literal, p.current().Line, p.current().Col)
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

	// Base Register
	if p.current().Type != types.TokenLBracket {
		return nil, fmt.Errorf("expected '[' for base register, got %s at line %d, col %d", p.current().Literal, p.current().Line, p.current().Col)
	}
	p.consume() // consume LBracket token
	if p.current().Type != types.TokenRegister {
		return nil, fmt.Errorf("expected base register after '[', got %s at line %d, col %d", p.current().Literal, p.current().Line, p.current().Col)
	}
	baseReg, err := utils.ParseRegister(p.current().Literal)
	if err != nil {
		return nil, fmt.Errorf("error parsing base register: %w", err)
	}
	p.consume() // consume base register token
	if p.current().Type != types.TokenRBracket {
		return nil, fmt.Errorf("expected ']' after base register, got %s at line %d, col %d", p.peek().Literal, p.peek().Line, p.peek().Col)
	}
	p.consume() // consume RBracket token

	instruction := &InstructionMemory{
		Mnemonic:     mnemonic,
		Condition:    condition,
		DestRegister: destReg,
		BaseRegister: baseReg,
		IBit:         iBit,
		PBit:         pBit,
		UBit:         uBit,
		BBit:         bBit,
		WBit:         wBit,
	}

	return instruction, nil
}

func (i *InstructionMemory) ToMachineCode(labels map[string]uint32) ([]byte, error) {
	var binary uint32
	binary |= types.ConditionToBits[i.Condition] << 28 // Condition code
	binary |= 1 << 26                                  // Data loading instruction
	binary |= i.IBit << 25                             // I bit, is offset immediate or from a register
	binary |= i.PBit << 24                             // P bit, pre or post index
	binary |= i.UBit << 23                             // U bit, add or subtract offset
	binary |= i.BBit << 22                             // B bit, byte or word
	binary |= i.WBit << 21                             // W bit, write back
	binary |= types.MnemonicToBits[i.Mnemonic] << 20   // Opcode bit (1 bit)
	binary |= i.BaseRegister << 16                     // Base register
	binary |= i.DestRegister << 12                     // Destination register
	binary |= 0 << 11                                  // Offset

	return utils.BitsToBytes(binary), nil
}
