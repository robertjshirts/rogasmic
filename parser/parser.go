package parser

import (
	"fmt"

	"github.com/robertjshirts/rogasmic/types"
)

type Parser struct {
	pos          int
	tokens       []types.Token
	instructions []types.Instruction
	labels       types.LabelMap // Maps label names to instruction numbers
}

func NewParser(tokens []types.Token) *Parser {
	if len(tokens) == 0 || tokens[len(tokens)-1].Type != types.TokenEOF {
		// Ensure the last token is EOF
		tokens = append(tokens, types.Token{Type: types.TokenEOF, Literal: "", Line: -1, Col: -1})
	}
	return &Parser{
		pos:          0,
		tokens:       tokens,
		instructions: make([]types.Instruction, 0),
		labels:       make(types.LabelMap),
	}
}

func (p *Parser) current() types.Token {
	if p.pos >= len(p.tokens) {
		return types.Token{Type: types.TokenEOF, Literal: "", Line: -1, Col: -1} // EOF token
	}
	return p.tokens[p.pos]
}

func (p *Parser) consume() {
	if p.pos < len(p.tokens) {
		p.pos++
	}
}

func (p *Parser) peek() types.Token {
	if p.pos+1 >= len(p.tokens) {
		return types.Token{Type: types.TokenEOF, Literal: "", Line: -1, Col: -1} // EOF token
	}
	return p.tokens[p.pos+1]
}

func (p *Parser) Parse() ([]types.Instruction, types.LabelMap, error) {

	for p.current().Type != types.TokenEOF {
		instructionCategory, ok := types.MnemonicTokenToCategory[p.current().Type]
		if !ok {
			if p.current().Type != types.TokenLabel {
				// Any other token is unexpected
				return nil, nil, fmt.Errorf("unexpected token %s at line %d, col %d", p.current().Literal, p.current().Line, p.current().Col)
			}

			p.labels[p.current().Literal] = uint32(len(p.instructions))
			p.consume() // consume label token
			continue    // skip to next token
		}

		switch instructionCategory {
		case types.MnemonicCategoryMOV:
			instruction, err := p.parseMOV()
			if err != nil {
				return nil, nil, fmt.Errorf("error parsing MOV instruction at line %d, col %d: %w", p.current().Line, p.current().Col, err)
			}
			p.instructions = append(p.instructions, instruction)
		case types.MnemonicCategoryLoadStore:
			instruction, err := p.parseMemory()
			if err != nil {
				return nil, nil, fmt.Errorf("error parsing memory instruction at line %d, col %d: %w", p.current().Line, p.current().Col, err)
			}
			p.instructions = append(p.instructions, instruction)
		case types.MnemonicCategoryLoadStoreMultiple:
			instruction, err := p.parseMemoryMultiple()
			if err != nil {
				return nil, nil, fmt.Errorf("error parsing block memory instruction at line %d, col %d: %w", p.current().Line, p.current().Col, err)
			}
			p.instructions = append(p.instructions, instruction)
		case types.MnemonicCategoryArithmetic:
			instruction, err := p.parseArithmetic()
			if err != nil {
				return nil, nil, fmt.Errorf("error parsing arithmetic instruction at line %d, col %d: %w", p.current().Line, p.current().Col, err)
			}
			p.instructions = append(p.instructions, instruction)
		case types.MnemonicCategoryBranch:
			instruction, err := p.parseBranch()
			if err != nil {
				return nil, nil, fmt.Errorf("error parsing branch instruction at line %d, col %d: %w", p.current().Line, p.current().Col, err)
			}
			p.instructions = append(p.instructions, instruction)
		case types.MnemonicCategoryBranchExchange:
			instruction, err := p.parseBranchExchange()
			if err != nil {
				return nil, nil, fmt.Errorf("error parsing branch exchange instruction at line %d, col %d: %w", p.current().Line, p.current().Col, err)
			}
			p.instructions = append(p.instructions, instruction)
		default:
			return nil, nil, fmt.Errorf("unknown instruction category at line %d, col %d", p.current().Line, p.current().Col)
		}
	}

	return p.instructions, p.labels, nil
}
