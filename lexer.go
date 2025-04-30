package main

import "unicode"

type Lexer struct {
	input   string
	pos     int
	readPos int
	ch      int
	line    int
}

func (l *Lexer) LexLine(line string, lineNo int) []Token {
	var toks []Token

	i, end := 0, len(line)

	for i < end {
		char := line[i]
		switch {
		case unicode.IsSpace(rune(char)):
			i++

		case char == ',':
			toks = append(toks, Token{Type: TokenComma, Value: string(char), Line: lineNo, Col: i})
			i++

		case char == '(':
			toks = append(toks, Token{Type: TokenLParen, Value: string(char), Line: lineNo, Col: i})

		case char == ')':
			toks = append(toks, Token{Type: TokenRParen, Value: string(char), Line: lineNo, Col: i})

		case unicode.IsLetter(rune(char)):
			start := i
			for i < end && (unicode.IsLetter(rune(line[i])) || unicode.IsDigit(rune(line[i]))) {
				i++
			}

		case char == ';':
			toks = append(toks, Token{Type: Comment, Value: line[i:]})
			i = end
		}
	}
}

func peek(line string, i int) rune {
	if i >= len(line) {
		return 0
	}
	return rune(line[i])
}
