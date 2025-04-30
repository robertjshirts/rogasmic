package parser

import (
	"fmt"

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

func ParseInstruction(tokens []types.Token) (types.Instruction, error) {
	if len(tokens) == 0 {
		return nil, nil
	}
	p := NewParser(tokens)
	tok := p.next()
	if tok.Type != types.TokenIdentifier {
		return nil, fmt.Errorf("expected opcode, got %q", tok.Value)
	}

	opcode, ok := types.OpCodeByName[tok.Value]
	if !ok {
		return nil, fmt.Errorf("unknown opcode %q", tok.Value)
	}

	tok = p.next()
	condCode, ok := types.CondByName[tok.Value]
	if !ok {
		return nil, fmt.Errorf("unknown condcode %q", tok.Value)
	}

	switch opcode {
	case types.OpMOVW, types.OpMOVT:
		return p.parseMOV()
	}

	return nil, nil
}

func (p *parser) parseMOV(opcode types.OpCode, condcode types.CondCode) (types.Instruction, error) {
	if len(p.tokens) == 0 {
		return nil, nil
	}

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
