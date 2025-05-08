package parser

import (
	"fmt"
	"strconv"

	"github.com/robertjshirts/rogasmic/types"
)

func ParseInstruction(tokens []types.Token) (types.Instruction, error) {
	if len(tokens) < 2 {
		return nil, fmt.Errorf("not enough tokens for instruction: expected at least 2, got %d", len(tokens))
	}
	// Parse the first token as opcode
	firstToken := tokens[0]
	tokens = tokens[1:]
	opCode, ok := types.OpCodesByLit[firstToken.Value]
	if !ok {
		return nil, fmt.Errorf("expected opcode, got %q", firstToken.Value)
	}

	// Try to parse second token as condition code, otherwise default to AL
	secondToken := tokens[0]
	condCode := types.CondCodesByLit["AL"]
	if secondToken.Type == types.TokenIdentifier {
		if tempCond, ok := types.CondCodesByLit[secondToken.Value]; ok {
			condCode = tempCond
			tokens = tokens[1:]
		}
	}

	// Using maps and switches for maximum performance
	switch types.InstructionTypes[opCode] {
	case types.MOVType:
		instruction, err := NewMOVInstruction(opCode, condCode, tokens)
		if err != nil {
			return nil, err
		}
		return instruction, nil
	case types.ArithmeticType:
		instruction, err := NewArithmeticInstruction(opCode, condCode, tokens)
		if err != nil {
			return nil, err
		}
		return instruction, nil
	case types.MemoryType:
		instruction, err := NewMemoryInstruction(opCode, condCode, tokens)
		if err != nil {
			return nil, err
		}
		return instruction, nil
	case types.BranchType:
		instruction, err := NewBranchInstruction(opCode, condCode, tokens)
		if err != nil {
			return nil, err
		}
		return instruction, nil
	case types.BranchExType:
		instruction, err := NewBranchExchangeInstruction(opCode, condCode, tokens)
		if err != nil {
			return nil, err
		}
		return instruction, nil
	}

	return nil, fmt.Errorf("unknown instruction type: %s", firstToken.Value)
}

// isSugar returns true if the token is a syntactic sugar token
func isSugar(token types.Token) bool {
	switch token.Type {
	case types.TokenComma:
	case types.TokenLParen:
	case types.TokenRParen:
	case types.TokenSemicolon:
		return true
	default:
		return false
	}
	return false
}

func parseRegister(tokenLit string) (uint32, error) {
	if len(tokenLit) < 2 {
		return 0, fmt.Errorf("invalid register literal: %q is too short", tokenLit)
	}
	if tokenLit[0] != 'R' && tokenLit[0] != 'r' {
		return 0, fmt.Errorf("invalid register literal: %q does not start with 'R' or 'r'", tokenLit)
	}
	register, err := strconv.ParseUint(tokenLit[1:], 10, 8)
	if err != nil {
		return 0, fmt.Errorf("invalid register literal: %q - %v", tokenLit, err)
	}
	if register > 14 { // 15 is the program counter
		return 0, fmt.Errorf("invalid register number: %d, must be between 0 and 14", register)
	}

	return uint32(register), nil
}

func parseImmediate(tokenLit string, maxBits int) (uint32, error) {
	immediateValue64, err := strconv.ParseUint(tokenLit, 0, maxBits)
	if err != nil {
		return 0, fmt.Errorf("invalid immediate value: %v", err)
	}
	return uint32(immediateValue64), nil
}

// Little endian
func bitsToBytes(bits uint32) []byte {
	bytes := make([]byte, 4)
	bytes[0] = byte(bits & 0xFF)
	bytes[1] = byte((bits >> 8) & 0xFF)
	bytes[2] = byte((bits >> 16) & 0xFF)
	bytes[3] = byte((bits >> 24) & 0xFF)
	return bytes
}
