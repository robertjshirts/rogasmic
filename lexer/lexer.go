package lexer

import (
	"fmt"
	"unicode"

	"github.com/robertjshirts/rogasmic/types"
	"github.com/robertjshirts/rogasmic/utils"
)

type lexer struct {
	input  string
	pos    int
	line   int
	col    int
	tokens []types.Token
}

func NewLexer(input string) *lexer {
	return &lexer{
		input:  input,
		pos:    0,
		line:   1,
		col:    1,
		tokens: []types.Token{},
	}
}

// current returns the current character. returns 0 at EOF.
func (l *lexer) current() byte {
	if l.pos >= len(l.input) {
		return 0 // EOF
	}
	return l.input[l.pos]
}

// consume advances to the next char. updates pos, col, and line.
func (l *lexer) consume() {
	if l.current() == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
	l.pos++
}

// peek returns the next character without consuming it. returns 0 at EOF.
func (l *lexer) peek() byte {
	if l.pos+1 >= len(l.input) {
		return 0 // EOF
	}
	return l.input[l.pos+1]
}

// skipWhitespace consumes spaces, tabs, and newlines
func (l *lexer) skipWhitespace() {
	for unicode.IsSpace(rune(l.current())) {
		l.consume()
	}
}

// consumeLit consumes a literal until it hits a non-literal character, and returns the literal
func (l *lexer) consumeLit() string {
	start := l.pos
	for utils.IsLiteralChar(l.current()) {
		l.consume()
	}
	return l.input[start:l.pos]
}

func (l *lexer) consumeComment() {
	// Consume until the end of the line or EOF
	for l.current() != '\n' && l.current() != 0 {
		l.consume()
	}
	if l.current() == '\n' {
		l.consume() // Consume the newline character
	}
}

// appendToken appends a new token to the lexer tokens slice. sets the column to the start of the literal.
func (l *lexer) appendToken(tokenType types.TokenType, literal string, startRow, startCol int) {
	l.tokens = append(l.tokens, types.Token{
		Type:    tokenType,
		Literal: literal,
		Line:    startRow,
		Col:     startCol,
	})
}

// Tokenize processes some input string and returns a slice of tokens.
func (l *lexer) Tokenize() ([]types.Token, error) {
	for l.current() != 0 {
		l.skipWhitespace()
		startRow := l.line // Store the starting row for the token
		startCol := l.col  // Store the starting column for the token
		switch l.current() {
		case ';':
			// Skip comments
			l.consumeComment()
		case ',':
			l.appendToken(types.TokenComma, string(l.current()), startRow, startCol)
			l.consume()
		case '[':
			l.appendToken(types.TokenLBracket, string(l.current()), startRow, startCol)
			l.consume()
		case ']':
			l.appendToken(types.TokenRBracket, string(l.current()), startRow, startCol)
			l.consume()
		case '#':
			l.consume() // Skip the #
			lit := l.consumeLit()
			if !utils.IsImmediate(lit) {
				return nil, fmt.Errorf("invalid immediate value: %s at line %d, col %d", lit, l.line, l.col)
			}
			l.appendToken(types.TokenImmediate, lit, startRow, startCol)
		default: // Handle registers, mnemonics, and labels/identifiers
			lit := l.consumeLit()
			if lit == "" {
				return nil, fmt.Errorf("unexpected input (not a valid lit or identifier) at line %d, col %d", l.line, l.col)
			}
			if utils.IsRegister(lit) {
				lit := utils.NormalizeRegister(lit) // For lr, sp, and pc, switch the actual register nums
				l.appendToken(types.TokenRegister, lit, startRow, startCol)
			} else if utils.IsOperation(lit) {
				l.appendToken(utils.GetMnemonicTokenType(lit), lit, startRow, startCol)
			} else { // After checking the literal against all reserved keywords, it is assumed to be an identifier
				if l.current() == ':' { // Check for label
					l.consume() // Consume the ':'
					l.appendToken(types.TokenLabel, lit, startRow, startCol)
				} else { // Otherwise, it's just an identifier
					l.appendToken(types.TokenIdentifier, lit, startRow, startCol)
				}
			}
		}
	}
	l.appendToken(types.TokenEOF, "", -1, -1) // Append EOF token at the end
	return l.tokens, nil
}
