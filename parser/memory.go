package parser

import (
	"fmt"

	"github.com/robertjshirts/rogasmic/types"
)

type MemoryInstruction struct {
	OpCode   types.OpCode
	CondCode types.CondCode
	Rd       uint32
	Rn       uint32
	IBit     uint32
	PBit     uint32
	UBit     uint32
	BBit     uint32
	WBit     uint32
}

func NewMemoryInstruction(opCode types.OpCode, condCode types.CondCode, tokens []types.Token) (*MemoryInstruction, error) {
	// Need at least 2 tokens, Rd and Rn
	if len(tokens) < 2 {
		return nil, fmt.Errorf("invalid memory instruction: expected at least 2 tokens (Rd, Rn), got %d", len(tokens))
	}

	if instType := types.InstructionTypes[opCode]; instType != types.MemoryType {
		return nil, fmt.Errorf("invalid memory instruction: opcode %d is not a memory type", opCode)
	}

	instruction := &MemoryInstruction{
		OpCode:   opCode,
		CondCode: condCode,
	}

	setRd := false
	setRn := false

	for _, token := range tokens {
		switch token.Type {
		case types.TokenRegister:
			reg, err := parseRegister(token.Value)
			if err != nil {
				return nil, err
			}
			if !setRd {
				instruction.Rd = reg
				setRd = true
			} else {
				instruction.Rn = reg
				setRn = true
			}
		default:
			return nil, fmt.Errorf("unexpected token in memory instruction: %s", token.Value)
		}
	}

	if !setRd {
		return nil, fmt.Errorf("missing destination register (Rd) in memory instruction")
	}
	if !setRn {
		return nil, fmt.Errorf("missing source register (Rn) in memory instruction")
	}

	instruction.IBit = 0
	instruction.PBit = 0
	instruction.UBit = 0
	instruction.BBit = 0
	instruction.WBit = 0

	return instruction, nil
}

func (m *MemoryInstruction) ToMachineCode() []byte {
	var binary uint32
	binary |= types.CondCodeBits[m.CondCode] << 28 // Condition code
	binary |= 1 << 26                              // Data loading instruction
	binary |= m.IBit << 25                         // I bit, is offset immediate or from a register
	binary |= m.PBit << 24                         // P bit, pre or post index
	binary |= m.UBit << 23                         // U bit, add or subtract offset
	binary |= m.BBit << 22                         // B bit, byte or word
	binary |= m.WBit << 21                         // W bit, write back
	binary |= types.OpCodeBits[m.OpCode] << 20     // Opcode bit (1 bit)
	binary |= m.Rn << 16                           // Source register
	binary |= m.Rd << 12                           // Destination register
	binary |= 0 << 11                              // Offset

	return bitsToBytes(binary)
}
