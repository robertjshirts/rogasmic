package parser

import (
	"fmt"

	"github.com/robertjshirts/rogasmic/types"
)

type BranchInstruction struct {
	OpCode   types.OpCode
	CondCode types.CondCode
	LBit     uint32
	Offset   uint32
}

func NewBranchInstruction(opCode types.OpCode, condCode types.CondCode, tokens []types.Token) (*BranchInstruction, error) {
	// Need at least one token, the offset
	if len(tokens) < 1 {
		return nil, fmt.Errorf("invalid branch instruction: expected at least 1 token (offset), got %d", len(tokens))
	}

	if instType := types.InstructionTypes[opCode]; instType != types.BranchType {
		return nil, fmt.Errorf("invalid branch instruction: opcode %d is not a branch type", opCode)
	}

	instruction := &BranchInstruction{
		OpCode:   opCode,
		CondCode: condCode,
	}

	for _, token := range tokens {
		switch token.Type {
		case types.TokenLBit:
			instruction.LBit = 1
		case types.TokenNumber, types.TokenHexNumber:
			offset, err := parseImmediate(token.Value, 24)
			if err != nil {
				return nil, err
			}
			instruction.Offset = offset
		default:
			return nil, fmt.Errorf("unexpected token in branch instruction: %s", token.Value)
		}
	}

	return instruction, nil
}

func (b *BranchInstruction) ToMachineCode() []byte {
	var binary uint32
	binary |= types.CondCodeBits[b.CondCode] << 28 // Condition code
	binary |= types.OpCodeBits[b.OpCode] << 25     // Branch opcode
	binary |= b.LBit << 24                         // Link bit
	binary |= b.Offset & 0xFFFFFF                  // Offset

	return bitsToBytes(binary)
}
