package parser

import (
	"fmt"

	"github.com/robertjshirts/rogasmic/types"
	"github.com/robertjshirts/rogasmic/utils"
)

type InstructionBranch struct {
	Mnemonic      types.MnemonicType
	Condition     types.ConditionType
	LBit          uint32
	Offset        uint32
	Label         string
	InstructionNo uint32 // Instruction number for relative addressing
}

type InstructionBranchExchange struct {
	Mnemonic     types.MnemonicType
	Condition    types.ConditionType
	BaseRegister uint32
}

func (p *Parser) parseBranch() (types.Instruction, error) {
	// Mnemonic
	mnemonic := types.TokenToMnemonic[p.current().Type]
	category, ok := types.MnemonicToCategory[mnemonic]
	if !ok || category != types.MnemonicCategoryBranch {
		return nil, fmt.Errorf("wrong instruction type! expected branch mnemonic, got %s", p.current().Literal)
	}

	// Get condition code and L suffix
	condition, lBit, err := utils.ParseBranchSuffixes(p.current().Literal)
	if err != nil {
		return nil, fmt.Errorf("error parsing branch suffixes: %w", err)
	}
	p.consume() // consume branch mnemonic token

	// Offset/Label
	var offset uint32
	var label string
	var instNo uint32
	if p.current().Type == types.TokenImmediate {
		offset, err = utils.ParseImmediate(p.current().Literal)
		if err != nil {
			return nil, fmt.Errorf("error parsing immediate value: %w", err)
		}
		p.consume() // consume immediate token
	} else if p.current().Type == types.TokenIdentifier {
		label = p.current().Literal
		instNo = uint32(len(p.instructions))
		p.consume() // consume label identifier token
	} else {
		return nil, fmt.Errorf("expected immediate value or label identifier after branch mnemonic, got %s at line %d, col %d", p.current().Literal, p.current().Line, p.current().Col)
	}

	instruction := &InstructionBranch{
		Mnemonic:      mnemonic,
		Condition:     condition,
		LBit:          lBit,
		Offset:        offset,
		Label:         label,
		InstructionNo: instNo,
	}

	return instruction, nil
}

func (i *InstructionBranch) ToMachineCode(labels map[string]uint32) ([]byte, error) {
	// Calculate offset if label is provided
	if i.Label != "" {
		labelAddress, ok := labels[i.Label]
		if !ok {
			return nil, fmt.Errorf("label %s not found", i.Label)
		}

		i.Offset = labelAddress - i.InstructionNo - 2 // -2 for arm pre-fetching and processing
	}

	var binary uint32
	binary |= types.ConditionToBits[i.Condition] << 28 // Set condition bits
	binary |= types.MnemonicToBits[i.Mnemonic] << 25
	binary |= i.LBit << 24        // Set L bit
	binary |= i.Offset & 0xFFFFFF // Set offset bits (24 bits)

	return utils.BitsToBytes(binary), nil
}

func (p *Parser) parseBranchExchange() (types.Instruction, error) {
	// Mnemonic
	mnemonic := types.TokenToMnemonic[p.current().Type]
	category, ok := types.MnemonicToCategory[mnemonic]
	if !ok || category != types.MnemonicCategoryBranchExchange {
		return nil, fmt.Errorf("wrong instruction type! expected branch exchange mnemonic, got %s", p.current().Literal)
	}

	// Get condition code
	condition, err := utils.ParseBranchExchangeSuffixes(p.current().Literal)
	if err != nil {
		return nil, fmt.Errorf("error parsing branch exchange suffixes: %w", err)
	}
	p.consume() // consume branch exchange mnemonic token

	// Base Register
	if p.current().Type != types.TokenRegister {
		return nil, fmt.Errorf("expected register after branch exchange mnemonic, got %s at line %d, col %d", p.current().Literal, p.current().Line, p.current().Col)
	}
	baseReg, err := utils.ParseRegister(p.current().Literal)
	if err != nil {
		return nil, fmt.Errorf("error parsing base register: %w", err)
	}
	p.consume() // consume base register token

	instruction := &InstructionBranchExchange{
		Mnemonic:     mnemonic,
		Condition:    condition,
		BaseRegister: baseReg,
	}

	return instruction, nil
}

func (i *InstructionBranchExchange) ToMachineCode(labels map[string]uint32) ([]byte, error) {
	var binary uint32
	binary |= types.ConditionToBits[i.Condition] << 28 // Set condition bits
	binary |= types.MnemonicToBits[i.Mnemonic] << 4    // Branch exchange bits (SO MANY)
	binary |= i.BaseRegister & 0xF

	return utils.BitsToBytes(binary), nil
}
