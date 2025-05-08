package parser

import (
	"fmt"

	"github.com/robertjshirts/rogasmic/types"
)

type ArithmeticInstruction struct {
	OpCode    types.OpCode
	CondCode  types.CondCode
	Rd        uint32
	Rn        uint32
	Immediate uint32
	SBit      uint32
	IBit      uint32
}

func NewArithmeticInstruction(opCode types.OpCode, condCode types.CondCode, tokens []types.Token) (*ArithmeticInstruction, error) {
	// Need at least 3 tokens, Rd, Rn, and immediate
	if len(tokens) < 3 { // Adjusted to 3 for Rd, Rn, and immediate
		return nil, fmt.Errorf("invalid arithmetic instruction: expected at least 3 tokens (Rd, Rn, immediate), got %d", len(tokens))
	}

	opCodeLit := ""
	for lit, code := range types.OpCodesByLit {
		if code == opCode {
			opCodeLit = lit
			break
		}
	}

	if instType := types.InstructionTypes[opCode]; instType != types.ArithmeticType {
		return nil, fmt.Errorf("invalid arithmetic instruction: opcode %s is not an arithmetic type", opCodeLit)
	}

	setRd := false
	setRn := false

	instruction := &ArithmeticInstruction{
		OpCode:   opCode,
		CondCode: condCode,
	}

	for _, token := range tokens {
		switch token.Type {
		case types.TokenRegister:
			reg, err := parseRegister(token.Value)
			if err != nil {
				return nil, fmt.Errorf("error parsing arithmetic instruction: %v", err)
			}
			if !setRd {
				instruction.Rd = reg
				setRd = true
			} else {
				instruction.Rn = reg
				setRn = true
			}
		case types.TokenNumber, types.TokenHexNumber:
			immediate, err := parseImmediate(token.Value, 8)
			if err != nil {
				return nil, fmt.Errorf("error parsing arithmetic instruction: %v", err)
			}
			instruction.Immediate = immediate
		case types.TokenSBit:
			instruction.SBit = 1
		default:
			return nil, fmt.Errorf("unexpected token in arithmetic instruction: %s", token.Value)
		}
	}

	if !setRd {
		return nil, fmt.Errorf("missing destination register (Rd) in arithmetic instruction")
	}
	if !setRn {
		return nil, fmt.Errorf("missing source register (Rn) in arithmetic instruction")
	}

	instruction.IBit = 1 // Not handling shift operations yet, so set IBit to 1

	return instruction, nil
}

func (a *ArithmeticInstruction) ToMachineCode() []byte {
	var binary uint32
	binary |= types.CondCodeBits[a.CondCode] << 28 // Condition code
	binary |= 0 << 26                              // Always 0 for arithmetic instructions
	binary |= a.IBit << 25                         // I bit, is offset immediate or from a register
	binary |= types.OpCodeBits[a.OpCode] << 21     // 4 bits
	binary |= a.SBit << 20                         // S bit, update CSPR register if set
	binary |= a.Rn << 16                           // Source register
	binary |= a.Rd << 12                           // Destination register
	binary |= a.Immediate & 0x0FF                  // Top 4 bits are used for something called 'rotate'

	return bitsToBytes(binary)
}
