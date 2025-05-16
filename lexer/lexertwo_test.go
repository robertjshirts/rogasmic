package lexer

import (
	"testing"
)

func TestLexerTwoBasicTokens(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []NewToken
	}{
		{
			"Lexes MOVW token",
			"MOVW",
			[]NewToken{
				{Type: TokenMOVW, Literal: "MOVW", Line: 1, Column: 1},
			},
		},
		{
			"Lexes MOVT token",
			"MOVT",
			[]NewToken{
				{Type: TokenMOVT, Literal: "MOVT", Line: 1, Column: 1},
			},
		},
		{
			"Lexes LDR token",
			"LDR",
			[]NewToken{
				{Type: TokenLDR, Literal: "LDR", Line: 1, Column: 1},
			},
		},
		{
			"Lexes STR token",
			"STR",
			[]NewToken{
				{Type: TokenSTR, Literal: "STR", Line: 1, Column: 1},
			},
		},
		{
			"Lexes ADD token",
			"ADD",
			[]NewToken{
				{Type: TokenADD, Literal: "ADD", Line: 1, Column: 1},
			},
		},
		{
			"Lexes SUB token",
			"SUB",
			[]NewToken{
				{Type: TokenSUB, Literal: "SUB", Line: 1, Column: 1},
			},
		},
		{
			"Lexes AND token",
			"AND",
			[]NewToken{
				{Type: TokenAND, Literal: "AND", Line: 1, Column: 1},
			},
		},
		{
			"Lexes ORR token",
			"ORR",
			[]NewToken{
				{Type: TokenORR, Literal: "ORR", Line: 1, Column: 1},
			},
		},
		{
			"Lexes BX token",
			"BX",
			[]NewToken{
				{Type: TokenBX, Literal: "BX", Line: 1, Column: 1},
			},
		},
		{
			"Lexes B token",
			"B",
			[]NewToken{
				{Type: TokenB, Literal: "B", Line: 1, Column: 1},
			},
		},
		{
			"Lexes Comma token",
			",",
			[]NewToken{
				{Type: TokenComma, Literal: ",", Line: 1, Column: 1},
			},
		},
		{
			"Lexes LParen token",
			"(",
			[]NewToken{
				{Type: TokenLParen, Literal: "(", Line: 1, Column: 1},
			},
		},
		{
			"Lexes RParen token",
			")",
			[]NewToken{
				{Type: TokenRParen, Literal: ")", Line: 1, Column: 1},
			},
		},
		{
			"Lexes Register token",
			"R7",
			[]NewToken{
				{Type: TokenRegister, Literal: "R7", Line: 1, Column: 1},
			},
		},
		{
			"Lexes Immediate token (decimal)",
			"#42",
			[]NewToken{
				{Type: TokenImmediate, Literal: "42", Line: 1, Column: 1},
			},
		},
		{
			"Lexes Immediate token (hex)",
			"#0x2A",
			[]NewToken{
				{Type: TokenImmediate, Literal: "42", Line: 1, Column: 1},
			},
		},
		{
			"Lexes S Suffix token",
			"S",
			[]NewToken{
				{Type: TokenSSuffix, Literal: "S", Line: 1, Column: 1},
			},
		},
		{
			"Lexes L Suffix token",
			"L",
			[]NewToken{
				{Type: TokenLSuffix, Literal: "L", Line: 1, Column: 1},
			},
		},
		{
			"Lexes ConditionCode token (AL)",
			"AL",
			[]NewToken{
				{Type: TokenConditionCode, Literal: "AL", Line: 1, Column: 1},
			},
		},
		{
			"Lexes ConditionCode token (EQ)",
			"EQ",
			[]NewToken{
				{Type: TokenConditionCode, Literal: "EQ", Line: 1, Column: 1},
			},
		},
		{
			"Lexes Identifier token",
			"foo",
			[]NewToken{
				{Type: TokenIdentifier, Literal: "foo", Line: 1, Column: 1},
			},
		},
		{
			"Lexes Label token",
			"label:",
			[]NewToken{
				{Type: TokenLabel, Literal: "label:", Line: 1, Column: 1},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := NewLexerTwo(c.input)
			tokens, errs := l.Tokenize()
			if len(errs) > 0 {
				t.Errorf("input=%s, unexpected errors: %v", c.input, errs)
			}
			if len(tokens) != len(c.expected) {
				t.Errorf("input=%s, wrong number of tokens: expected=%d, got=%d",
					c.input, len(c.expected), len(tokens))
			}
			for i, expected := range c.expected {
				if tokens[i] != expected {
					t.Errorf("input=%s, token[%d] wrong: expected=%#v, got=%#v",
						c.input, i, expected, tokens[i])
				}
			}
		})
	}
}

func TestLexerTwoBasicInstructions(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []NewToken
	}{
		{
			"Basic MOV Inst with immediate",
			"MOVW R4, 0",
			[]NewToken{
				{Type: TokenMOVW, Literal: "MOVW", Line: 1, Column: 1},
				{Type: TokenRegister, Literal: "R4", Line: 1, Column: 6},
				{Type: TokenComma, Literal: ",", Line: 1, Column: 8},
				{Type: TokenImmediate, Literal: "0", Line: 1, Column: 10},
			},
		},
		{
			"Basic MOV Inst with immediate hex",
			"MOVT R4, 0x3F20",
			[]NewToken{
				{Type: TokenMOVT, Literal: "MOVT", Line: 1, Column: 1},
				{Type: TokenRegister, Literal: "R4", Line: 1, Column: 6},
				{Type: TokenComma, Literal: ",", Line: 1, Column: 8},
				{Type: TokenImmediate, Literal: "0x3F20", Line: 1, Column: 10},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := NewLexerTwo(c.input)
			tokens, errs := l.Tokenize()
			if len(errs) > 0 {
				t.Errorf("input=%s, unexpected errors: %v", c.input, errs)
			}
			if len(tokens) != len(c.expected) {
				t.Errorf("input=%s, wrong number of tokens: expected=%d, got=%d",
					c.input, len(c.expected), len(tokens))
			}
			for i, expected := range c.expected {
				if tokens[i] != expected {
					t.Errorf("input=%s, token[%d] wrong: expected=%#v, got=%#v",
						c.input, i, expected, tokens[i])
				}
			}
		})
	}
}
