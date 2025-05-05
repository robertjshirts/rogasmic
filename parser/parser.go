package parser

import (
	"fmt"
	"strconv"

	"github.com/robertjshirts/rogasmic/types"
)

type parser struct {
	tokens  []types.Token
	token   types.Token
	pos     int
	readPos int
}

func NewParser(tokens []types.Token) *parser {
	return &parser{
		tokens: tokens,
		token:  types.Token{},
		pos:    0,
	}
}

// next advances to the next token and returns it.
// If we run past the end, we synthesize an EOF token.
func (p *parser) next() types.Token {
	p.pos++
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	return types.Token{Type: types.TokenEOF}
}

// peek looks one token ahead without consuming.
func (p *parser) peek() types.Token {
	if p.pos+1 < len(p.tokens) {
		return p.tokens[p.pos+1]
	}
	return types.Token{Type: types.TokenEOF}
}

func ParseInstruction(tokens []types.Token) ([]byte, error) {
	if len(tokens) == 0 {
		return nil, nil
	}
	op := tokens[0]
	if op.Type != types.TokenIdentifier {
		return nil, fmt.Errorf("expected opcode, got %q", op.Value)
	}

	if _, ok := types.MovSet[op.Value]; ok {
		inst, err := parseMOVOperands(tokens)
		if err != nil {
			return nil, err
		}
		return inst, nil
	} else if _, ok := types.ArithmeticSet[op.Value]; ok {
		inst, err := parseArithmeticOperands(tokens)
		if err != nil {
			return nil, err
		}
		return inst, nil
	} else if _, ok := types.MemorySet[op.Value]; ok {
		inst, err := parseMemoryOperands(tokens)
		if err != nil {
			return nil, err
		}
		return inst, nil
	} else if _, ok := types.BranchSet[op.Value]; ok {
		inst, err := parseBranchOperands(tokens)
		if err != nil {
			return nil, err
		}
		return inst, nil
	}

	return nil, nil
}

func parseBranchOperands(tokens []types.Token) ([]byte, error) {
	if len(tokens) != 2 && len(tokens) != 3 {
		return nil, fmt.Errorf("branch instruction requires 2 operands, got %d", len(tokens))
	}
	condBits := types.CondCodeBits["AL"]
	var opBits uint32
	var offset uint32
	for _, token := range tokens {
		if token.Type == types.TokenIdentifier {
			// Try to parse as a condition code or opcode
			if bits, ok := types.CondCodeBits[token.Value]; ok {
				condBits = uint32(bits)
			} else if bits, ok := types.OpCodeBits[token.Value]; ok {
				opBits = uint32(bits)
			}
			continue
		} else if token.Type == types.TokenNumber || token.Type == types.TokenHexNumber {
			offset64, err := strconv.ParseUint(token.Value, 0, 24)
			if err != nil || offset64 > 0xFFFFFF {
				return nil, fmt.Errorf("invalid offset: %v", err)
			}
			offset = uint32(offset64)
			continue
		}
	}

	var binary uint32
	binary |= condBits << 28      // Condition code
	binary |= opBits << 25        // Opcode
	binary |= 0 << 24             // L bit, link bit
	binary |= offset & 0x00FFFFFF // offset

	machineCode := make([]byte, 4)
	// little endian
	machineCode[0] = byte(binary & 0xFF)
	machineCode[1] = byte((binary >> 8) & 0xFF)
	machineCode[2] = byte((binary >> 16) & 0xFF)
	machineCode[3] = byte((binary >> 24) & 0xFF)

	return machineCode, nil
}

func parseMemoryOperands(tokens []types.Token) ([]byte, error) {
	if len(tokens) != 3 && len(tokens) != 4 {
		return nil, fmt.Errorf("memory instruction requires 4-5 operands, got %d", len(tokens))
	}
	condBits := types.CondCodeBits["AL"]
	var opBit uint32
	var sourceRegister uint32
	var destinationRegister uint32
	// var offsetBits uint32 // Not calculating offset yet
	for _, token := range tokens {
		if token.Type == types.TokenIdentifier {
			// Try to parse as a condition code or opcode
			if bits, ok := types.CondCodeBits[token.Value]; ok {
				condBits = uint32(bits)
			} else if bits, ok := types.OpCodeBits[token.Value]; ok {
				opBit = uint32(bits)
			}
			continue
		}
		if token.Type == types.TokenRegister {
			register64, err := strconv.ParseUint(token.Value[1:], 10, 8)
			if err != nil || register64 > 15 {
				return nil, fmt.Errorf("invalid register: %v", err)
			}
			register := uint32(register64)
			if destinationRegister == 0 {
				destinationRegister = register
			} else {
				sourceRegister = register
			}
			continue
		}
		if token.Type == types.TokenNumber || token.Type == types.TokenHexNumber {
			panic("LDR/STR Offsets not implemented yet")
		}
	}

	var binary uint32
	binary |= condBits << 28                     // Condition code
	binary |= 1 << 26                            // Data loading instruction
	binary |= 0 << 25                            // I bit, is offset immediate or from a register
	binary |= 0 << 24                            // P bit, pre or post index
	binary |= 0 << 23                            // U bit, add or subtract offset
	binary |= 0 << 22                            // B bit, byte or word
	binary |= 0 << 21                            // W bit, write back
	binary |= opBit << 20                        // Opcode
	binary |= (sourceRegister & 0xFF) << 16      // Source register
	binary |= (destinationRegister & 0xFF) << 12 // Destination register
	binary |= 0 << 11                            // offset

	machineCode := make([]byte, 4)
	// little endian
	machineCode[0] = byte(binary & 0xFF)
	machineCode[1] = byte((binary >> 8) & 0xFF)
	machineCode[2] = byte((binary >> 16) & 0xFF)
	machineCode[3] = byte((binary >> 24) & 0xFF)
	return machineCode, nil
}

func parseArithmeticOperands(tokens []types.Token) ([]byte, error) {
	if len(tokens) != 4 && len(tokens) != 5 {
		return nil, fmt.Errorf("arithmetic instruction requires 4-5 operands, got %d", len(tokens))
	}
	condBits := types.CondCodeBits["AL"]
	var opBits uint32
	var operandRegister uint32     // register
	var destinationRegister uint32 // destination register
	var immediateValue uint32      // Technically could be a register, but we aren't handling that for now
	var sBit uint32                // S bit, set if we are setting flags
	for _, token := range tokens {
		if token.Type == types.TokenIdentifier {
			// Try to parse as a condition code or opcode
			if bits, ok := types.CondCodeBits[token.Value]; ok {
				condBits = uint32(bits)
			} else if bits, ok := types.OpCodeBits[token.Value]; ok {
				opBits = uint32(bits)
			}
			continue
		}
		if token.Type == types.TokenRegister {
			register64, err := strconv.ParseUint(token.Value[1:], 10, 8)
			if err != nil || register64 > 15 {
				return nil, fmt.Errorf("invalid register: %v", err)
			}
			register := uint32(register64)
			if destinationRegister == 0 {
				destinationRegister = register
			} else {
				operandRegister = register
			}
			continue
		}
		if token.Type == types.TokenNumber || token.Type == types.TokenHexNumber {
			immediateValue64, err := strconv.ParseUint(token.Value, 0, 8)
			if err != nil || immediateValue64 > 255 {
				return nil, fmt.Errorf("invalid immediate value: %v", err)
			}
			immediateValue = uint32(immediateValue64)
			continue
		}
		if token.Type == types.TokenSBit {
			sBit = 1
			continue
		}
	}
	var binary uint32
	binary |= condBits << 28            // Condition code
	binary |= 0 << 26                   // Data processing instruction
	binary |= 1 << 25                   // I bit, always immediate bc we aren't handling registers yet
	binary |= opBits << 21              // Opcode
	binary |= sBit << 20                // S bit
	binary |= operandRegister << 16     // operand register
	binary |= destinationRegister << 12 // destination register
	binary |= immediateValue & 0x0FF    // immediate value, ignore shift

	machineCode := make([]byte, 4)
	// little endian
	machineCode[0] = byte(binary & 0xFF)
	machineCode[1] = byte((binary >> 8) & 0xFF)
	machineCode[2] = byte((binary >> 16) & 0xFF)
	machineCode[3] = byte((binary >> 24) & 0xFF)

	return machineCode, nil
}

func parseMOVOperands(tokens []types.Token) ([]byte, error) {
	if len(tokens) != 3 && len(tokens) != 4 {
		return nil, fmt.Errorf("MOVW requires 3-4 operands, got %d", len(tokens))
	}

	var opBits uint32
	condBits := types.CondCodeBits["AL"]
	var immediateValue uint32
	var register uint32
	for _, token := range tokens {
		if token.Type == types.TokenIdentifier {
			// Try to parse as a condition code or opcode
			if bits, ok := types.CondCodeBits[token.Value]; ok {
				condBits = uint32(bits)
			} else if bits, ok := types.OpCodeBits[token.Value]; ok {
				opBits = uint32(bits)
			}
			continue
		}
		if token.Type == types.TokenRegister {
			register64, err := strconv.ParseUint(token.Value[1:], 10, 8)
			if err != nil || register64 > 15 {
				return nil, fmt.Errorf("invalid register: %v", err)
			}
			register = uint32(register64)
			continue
		}
		if token.Type == types.TokenNumber || token.Type == types.TokenHexNumber {
			immediateValue64, err := strconv.ParseUint(token.Value, 0, 16)
			if err != nil {
				return nil, fmt.Errorf("invalid immediate value: %v", err)
			}
			immediateValue = uint32(immediateValue64)
		}
	}

	var binary uint32
	binary |= condBits << 28
	binary |= 3 << 24                              // MOV opcode
	binary |= opBits << 20                         // MOVW or MOVT
	binary |= (immediateValue >> 12) & 0x00F << 16 // top 4 bits of immediate value
	binary |= register << 12                       // Destination register
	binary |= immediateValue & 0x0FFF              // bottom 12 bits of immediate value

	machineCode := make([]byte, 4)
	// little endian
	machineCode[0] = byte(binary & 0xFF)
	machineCode[1] = byte((binary >> 8) & 0xFF)
	machineCode[2] = byte((binary >> 16) & 0xFF)
	machineCode[3] = byte((binary >> 24) & 0xFF)

	return machineCode, nil
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
