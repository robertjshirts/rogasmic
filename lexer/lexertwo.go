package lexer

import (
	"fmt"
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

func (l *LexerTwo) Tokenize() ([]NewToken, []error) {
	var tokens []NewToken
	var errs []error
	for {
		// skip whitespace
		for unicode.IsSpace(rune(l.ch)) {
			l.consumeChar()
		}

		tok := NewToken{
			Line:   l.line,
			Column: l.column,
		}

		switch l.ch {
		case 0:
			return tokens, errs
		case ',':
			tok.Type = TokenComma
			tok.Literal = string(l.ch)
			l.consumeChar()
		case '[':
			tok.Type = TokenLBracket
			tok.Literal = string(l.ch)
			l.consumeChar()
		case ']':
			tok.Type = TokenRBracket
			tok.Literal = string(l.ch)
			l.consumeChar()
		case '#':
			l.consumeChar()
			lit := l.consumeLiteral()
			if !isImmediate(lit) {
				tok.Type = TokenError
				tok.Literal = lit
				errs = append(errs, fmt.Errorf("invalid immediate value: %s at %d:%d", lit, l.line, l.column))
			} else {
				tok.Type = TokenImmediate
				tok.Literal = lit
			}
		default:
			lit := l.consumeLiteral()

			splitFound := false
			// Try to split the literal into mnemonic and condition code
			for i := range lit {
				// Check if valid mnemonic and condition code
				mnemonic := strings.ToUpper(lit[:i])
				condition := strings.ToUpper(lit[i:])
				mnemonicTyp, mnemonicOk := MnemonicsByLit[strings.ToUpper(mnemonic)]
				conditionTyp, conditionOk := ConditionCodesByLit[strings.ToUpper(condition)]

				if mnemonicOk && conditionOk {
					mnemonicToken := NewToken{
						Type:    mnemonicTyp,
						Literal: mnemonic,
						Line:    tok.Line,
						Column:  tok.Column,
					}
					conditionToken := NewToken{
						Type:    conditionTyp,
						Literal: condition,
						Line:    tok.Line,
						Column:  tok.Column + i,
					}

					tokens = append(tokens, mnemonicToken, conditionToken)

					splitFound = true
					break
				}
			}

			if splitFound {
				continue
			}

			if isLabel(lit) {
				tok.Type = TokenLabel
				tok.Literal = lit
			} else if isRegister(lit) {
				tok.Type = TokenRegister
				tok.Literal = lit
			} else if typ, ok := NewTokenTypesByLit[strings.ToUpper(lit)]; ok {
				tok.Type = typ
				tok.Literal = lit
			} else {
				tok.Type = TokenIdentifier
				tok.Literal = lit
			}
		}

		tokens = append(tokens, tok)
	}
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

func isImmediate(lit string) bool {
	lit = strings.ToUpper(lit)
	// Trim off the '0x' prefix if present so we can check for hex digits
	isHex := len(lit) > 2 && lit[0] == '0' && lit[1] == 'X'
	if isHex {
		lit = lit[2:]
	}
	for _, c := range lit {
		if !unicode.Is(unicode.Hex_Digit, c) {
			return false
		}
	}
	return true
}
