package parser

import (
	"fmt"

	"github.com/robertjshirts/rogasmic/types"
)

type MOVInstruction struct {
	OpCode    types.OpCode
	CondCode  types.CondCode
	Rd        uint32
	Immediate uint32
}

func NewMOVInstruction(opCode types.OpCode, condCode types.CondCode, tokens []types.Token) (*MOVInstruction, error) {
	// Need at least 2 tokens: Rd and immediate
	if len(tokens) < 2 {
		return nil, fmt.Errorf("invalid MOV instruction: expected at least 2 tokens (Rd, immediate), got %d", len(tokens))
	}

	if opCode != types.OpMOVW && opCode != types.OpMOVT {
		return nil, fmt.Errorf("invalid MOV instruction: opcode %d is not MOVW or MOVT", opCode)
	}

	instruction := &MOVInstruction{
		OpCode:   opCode,
		CondCode: condCode,
	}

	for _, token := range tokens {
		if isSugar(token) {
			continue
		}
		switch token.Type {
		case types.TokenRegister:
			rd, err := parseRegister(token.Value)
			if err != nil {
				return nil, fmt.Errorf("error parsing MOV instruction: %v", err)
			}
			instruction.Rd = rd
		case types.TokenNumber, types.TokenHexNumber:
			immediate, err := parseImmediate(token.Value, 16)
			if err != nil {
				return nil, fmt.Errorf("error parsing MOV instruction: %v", err)
			}
			instruction.Immediate = immediate
		default:
			return nil, fmt.Errorf("unexpected token in MOV instruction: %s", token.Value)
		}
	}

	return instruction, nil
}

func (m *MOVInstruction) ToMachineCode() []byte {
	var binary uint32
	binary |= types.CondCodeBits[m.CondCode] << 28 // Condition code
	binary |= types.OpCodeBits[m.OpCode] << 20     // 8 whole bits
	binary |= (m.Immediate >> 12) & 0xF << 16      // top 4 bits of immediate
	binary |= m.Rd << 12                           // Destination register
	binary |= m.Immediate & 0x0FFF                 // bottom 12 bits of immediate

	return bitsToBytes(binary)
}
