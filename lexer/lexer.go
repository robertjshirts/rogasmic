package lexer

import (
	"unicode"

	"github.com/robertjshirts/rogasmic/types"
)

type Lexer struct {
	input   string
	pos     int
	readPos int
	ch      byte
	line    int
}

func NewLexer(input string, line int) *Lexer {
	l := &Lexer{
		input:   input,
		line:    line,
		readPos: 0,
	}
	l.consumeChar()
	return l
}

func (l *Lexer) consumeChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
}

func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

func (l *Lexer) readIdentifier() string {
	start := l.pos
	for isLetter(l.ch) || isDigit(l.ch) {
		l.consumeChar()
	}
	return l.input[start:l.pos]
}

func (l *Lexer) readNumber() (string, types.TokenType) {
	start := l.pos
	// hex?
	if l.ch == '0' && (l.peekChar() == 'x' || l.peekChar() == 'X') {
		l.consumeChar() // 0
		l.consumeChar() // x
		for isHexDigit(l.ch) {
			l.consumeChar()
		}
		return l.input[start:l.pos], types.TokenHexNumber
	}

	// decimal
	for isDigit(l.ch) {
		l.consumeChar()
	}

	return l.input[start:l.pos], types.TokenNumber
	// No binary. fuck binary. me and my homies hate binary
}

func (l *Lexer) nextToken() types.Token {
	var tok types.Token

	for unicode.IsSpace(rune(l.ch)) || l.ch == ',' || l.ch == '(' || l.ch == ')' {
		l.consumeChar()
	}

	tok.Line = l.line
	tok.Col = l.pos

	switch l.ch {
	case 0:
		tok.Type = types.TokenEOF
		tok.Value = ""
	case ';':
		panic("Comment not implemented")
	case 'S':
		if unicode.IsSpace(rune(l.peekChar())) {
			tok.Type = types.TokenSBit
			tok.Value = "S"
			l.consumeChar()
			return tok
		}
		fallthrough
	default:
		if unicode.IsLetter(rune(l.ch)) {
			lit := l.readIdentifier()
			tok.Value = lit
			tok.Type = types.TokenIdentifier
			if isRegisterLiteral(lit) {
				tok.Type = types.TokenRegister
			}
			return tok
		}

		if unicode.IsDigit(rune(l.ch)) {
			start := l.pos
			lit, typ := l.readNumber()
			return types.Token{Type: typ, Value: lit, Line: l.line, Col: start}
		}

		tok = types.Token{Type: types.TokenError, Value: string(l.ch), Line: l.line, Col: l.pos}
	}

	l.consumeChar()
	return tok
}

func LexLine(line string, lineNo int) []types.Token {
	l := NewLexer(line, lineNo)
	var tokens []types.Token
	for tok := l.nextToken(); tok.Type != types.TokenEOF; tok = l.nextToken() {
		if tok.Type == types.TokenError {
			panic("Lexer - Unexpected token: " + tok.Value)
		}
		tokens = append(tokens, tok)
	}
	return tokens
}

func isRegisterLiteral(lit string) bool {
	if len(lit) < 2 {
		return false
	}

	if lit[0] != 'R' && lit[0] != 'r' {
		return false
	}
	for i := 1; i < len(lit); i++ {
		if !unicode.IsDigit(rune(lit[i])) {
			return false
		}
	}
	return true
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isHexDigit(ch byte) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}
