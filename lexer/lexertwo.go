package lexer

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type LexerTwo struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
	column       int
}

func NewLexerTwo(input string) *LexerTwo {
	l := &LexerTwo{
		input:  input,
		line:   1,
		column: 0,
	}
	l.consumeChar()
	return l
}

// consumeChar advances the lexer by one character, updating line and column.
func (l *LexerTwo) consumeChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII code 0, signifies EOF
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++

	// Update line and column
	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

func (l *LexerTwo) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *LexerTwo) lexToken() (NewToken, error) {
	for unicode.IsSpace(rune(l.ch)) {
		l.consumeChar()
	}

	tok := NewToken{
		Line:   l.line,
		Column: l.column,
	}

	// Single character tokens
	switch l.ch {
	case 0:
		tok.Type = TokenEOF
		tok.Literal = ""
		// Don't consume EOF
		return tok, nil
	case ',':
		tok.Type = TokenComma
		tok.Literal = string(l.ch)
		l.consumeChar()
		return tok, nil
	case '(':
		tok.Type = TokenLParen
		tok.Literal = string(l.ch)
		l.consumeChar()
		return tok, nil
	case ')':
		tok.Type = TokenRParen
		tok.Literal = string(l.ch)
		l.consumeChar()
		return tok, nil
	}

	if l.ch == '#' {
		l.consumeChar()
		lit := l.consumeLiteral()
		if lit == "" {
			tok.Type = TokenError
			tok.Literal = "immediate value is empty"
			return tok, fmt.Errorf("immediate value is empty at %d:%d", l.line, l.column)
		}
		imm, err := strconv.ParseUint(lit, 0, 32)
		if err != nil {
			tok.Type = TokenError
			tok.Literal = lit
			return tok, fmt.Errorf("error parsing immediate value: %v", err)
		}
		tok.Type = TokenImmediate
		tok.Literal = strconv.FormatUint(imm, 10) // Format to base 10
		return tok, nil
	} else if unicode.IsLetter(rune(l.ch)) {
		lit := l.consumeLiteral()
		if isLabel(lit) {
			tok.Type = TokenLabel
			tok.Literal = lit
			return tok, nil
		} else if isRegister(lit) {
			tok.Type = TokenRegister
			tok.Literal = lit
			return tok, nil
		} else if typ, ok := NewTokenTypesByLit[strings.ToUpper(lit)]; ok {
			tok.Type = typ
			tok.Literal = lit
			return tok, nil
		} else {
			tok.Type = TokenIdentifier
			tok.Literal = lit
			return tok, nil
		}
	}

	return tok, nil
}

func (l *LexerTwo) Tokenize() ([]NewToken, []error) {
	var tokens []NewToken
	var errs []error
	for {
		tok, err := l.lexToken()
		if err != nil {
			errs = append(errs, err)
		}
		if tok.Type == TokenEOF {
			break
		}
		tokens = append(tokens, tok)
	}
	return tokens, errs
}

func (l *LexerTwo) consumeLiteral() string {
	start := l.position
	for unicode.IsLetter(rune(l.ch)) || unicode.IsDigit(rune(l.ch)) || l.ch == '_' || l.ch == ':' {
		l.consumeChar()
	}
	return l.input[start:l.position]
}

func isLabel(lit string) bool {
	return lit[len(lit)-1] == ':' && len(lit) > 1
}

func isRegister(lit string) bool {
	return len(lit) > 1 && lit[0] == 'R' && isDigit(lit[1])
}
