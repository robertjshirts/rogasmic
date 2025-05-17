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
			"Lexes SUBS token",
			"SUBS",
			[]NewToken{
				{Type: TokenSUBS, Literal: "SUBS", Line: 1, Column: 1},
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
			"Lexes LBracket token",
			"[",
			[]NewToken{
				{Type: TokenLBracket, Literal: "[", Line: 1, Column: 1},
			},
		},
		{
			"Lexes RPBracket token",
			"]",
			[]NewToken{
				{Type: TokenRBracket, Literal: "]", Line: 1, Column: 1},
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
				{Type: TokenImmediate, Literal: "0x2A", Line: 1, Column: 1},
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
			"Lexes AL token",
			"AL",
			[]NewToken{
				{Type: TokenAL, Literal: "AL", Line: 1, Column: 1},
			},
		},
		{
			"Lexes EQ token",
			"EQ",
			[]NewToken{
				{Type: TokenEQ, Literal: "EQ", Line: 1, Column: 1},
			},
		},
		{
			"Lexes PL token",
			"PL",
			[]NewToken{
				{Type: TokenPL, Literal: "PL", Line: 1, Column: 1},
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

func TestLexerMnemonicWithSuffix(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []NewToken
	}{
		{
			"Branch with PL condition code",
			"BPL label",
			[]NewToken{
				{Type: TokenB, Literal: "B", Line: 1, Column: 1},
				{Type: TokenPL, Literal: "PL", Line: 1, Column: 2},
				{Type: TokenIdentifier, Literal: "label", Line: 1, Column: 5},
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
			"MOVW R4, #0",
			[]NewToken{
				{Type: TokenMOVW, Literal: "MOVW", Line: 1, Column: 1},
				{Type: TokenRegister, Literal: "R4", Line: 1, Column: 6},
				{Type: TokenComma, Literal: ",", Line: 1, Column: 8},
				{Type: TokenImmediate, Literal: "0", Line: 1, Column: 10},
			},
		},
		{
			"Basic MOV Inst with immediate hex",
			"MOVT R4, #0x3F20",
			[]NewToken{
				{Type: TokenMOVT, Literal: "MOVT", Line: 1, Column: 1},
				{Type: TokenRegister, Literal: "R4", Line: 1, Column: 6},
				{Type: TokenComma, Literal: ",", Line: 1, Column: 8},
				{Type: TokenImmediate, Literal: "0x3F20", Line: 1, Column: 10},
			},
		},
		{
			"ADD with two registers and immediate",
			"ADD R2, R4, #8",
			[]NewToken{
				{Type: TokenADD, Literal: "ADD", Line: 1, Column: 1},
				{Type: TokenRegister, Literal: "R2", Line: 1, Column: 5},
				{Type: TokenComma, Literal: ",", Line: 1, Column: 7},
				{Type: TokenRegister, Literal: "R4", Line: 1, Column: 9},
				{Type: TokenComma, Literal: ",", Line: 1, Column: 11},
				{Type: TokenImmediate, Literal: "8", Line: 1, Column: 13},
			},
		},
		{
			"LDR with register and memory reference",
			"LDR R3, [R2]",
			[]NewToken{
				{Type: TokenLDR, Literal: "LDR", Line: 1, Column: 1},
				{Type: TokenRegister, Literal: "R3", Line: 1, Column: 5},
				{Type: TokenComma, Literal: ",", Line: 1, Column: 7},
				{Type: TokenLBracket, Literal: "[", Line: 1, Column: 9},
				{Type: TokenRegister, Literal: "R2", Line: 1, Column: 10},
				{Type: TokenRBracket, Literal: "]", Line: 1, Column: 12},
			},
		},
		{
			"ORR with two registers and immediate",
			"ORR R3, R3, #8",
			[]NewToken{
				{Type: TokenORR, Literal: "ORR", Line: 1, Column: 1},
				{Type: TokenRegister, Literal: "R3", Line: 1, Column: 5},
				{Type: TokenComma, Literal: ",", Line: 1, Column: 7},
				{Type: TokenRegister, Literal: "R3", Line: 1, Column: 9},
				{Type: TokenComma, Literal: ",", Line: 1, Column: 11},
				{Type: TokenImmediate, Literal: "8", Line: 1, Column: 13},
			},
		},
		{
			"STR with register and memory reference",
			"STR R3, [R2]",
			[]NewToken{
				{Type: TokenSTR, Literal: "STR", Line: 1, Column: 1},
				{Type: TokenRegister, Literal: "R3", Line: 1, Column: 5},
				{Type: TokenComma, Literal: ",", Line: 1, Column: 7},
				{Type: TokenLBracket, Literal: "[", Line: 1, Column: 9},
				{Type: TokenRegister, Literal: "R2", Line: 1, Column: 10},
				{Type: TokenRBracket, Literal: "]", Line: 1, Column: 12},
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

func TestLexerTwoMultiLineInstructions(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []NewToken
	}{
		{
			"Label, SUB, and B",
			"label:\nSUBS R5, R5, #0x01\nBPL label",
			[]NewToken{
				{Type: TokenLabel, Literal: "label:", Line: 1, Column: 1},
				{Type: TokenSUBS, Literal: "SUBS", Line: 2, Column: 1},
				{Type: TokenRegister, Literal: "R5", Line: 2, Column: 6},
				{Type: TokenComma, Literal: ",", Line: 2, Column: 8},
				{Type: TokenRegister, Literal: "R5", Line: 2, Column: 10},
				{Type: TokenComma, Literal: ",", Line: 2, Column: 12},
				{Type: TokenImmediate, Literal: "0x01", Line: 2, Column: 14},
				{Type: TokenB, Literal: "B", Line: 3, Column: 1},
				{Type: TokenPL, Literal: "PL", Line: 3, Column: 2},
				{Type: TokenIdentifier, Literal: "label", Line: 3, Column: 5},
			},
		},
		{
			"Multi-line with labels, MOV, ADD, STR, and branch",
			`start:
ADD R3, R4, #0x1C
MOVW R2, #0x0000
MOVT R2, #0x0020
STR R2, [R3]
BL delay`,
			[]NewToken{
				{Type: TokenLabel, Literal: "start:", Line: 1, Column: 1},
				{Type: TokenADD, Literal: "ADD", Line: 2, Column: 1},
				{Type: TokenRegister, Literal: "R3", Line: 2, Column: 5},
				{Type: TokenComma, Literal: ",", Line: 2, Column: 7},
				{Type: TokenRegister, Literal: "R4", Line: 2, Column: 9},
				{Type: TokenComma, Literal: ",", Line: 2, Column: 11},
				{Type: TokenImmediate, Literal: "0x1C", Line: 2, Column: 13},
				{Type: TokenMOVW, Literal: "MOVW", Line: 3, Column: 1},
				{Type: TokenRegister, Literal: "R2", Line: 3, Column: 6},
				{Type: TokenComma, Literal: ",", Line: 3, Column: 8},
				{Type: TokenImmediate, Literal: "0x0000", Line: 3, Column: 10},
				{Type: TokenMOVT, Literal: "MOVT", Line: 4, Column: 1},
				{Type: TokenRegister, Literal: "R2", Line: 4, Column: 6},
				{Type: TokenComma, Literal: ",", Line: 4, Column: 8},
				{Type: TokenImmediate, Literal: "0x0020", Line: 4, Column: 10},
				{Type: TokenSTR, Literal: "STR", Line: 5, Column: 1},
				{Type: TokenRegister, Literal: "R2", Line: 5, Column: 5},
				{Type: TokenComma, Literal: ",", Line: 5, Column: 7},
				{Type: TokenLBracket, Literal: "[", Line: 5, Column: 9},
				{Type: TokenRegister, Literal: "R3", Line: 5, Column: 10},
				{Type: TokenRBracket, Literal: "]", Line: 5, Column: 12},
				{Type: TokenBL, Literal: "BL", Line: 6, Column: 1},
				{Type: TokenIdentifier, Literal: "delay", Line: 6, Column: 4},
			},
		},
		{
			"Loop with SUBS and BPL",
			`delay:
MOVW R5, #0xFFFF
MOVT R5, #0x000F
loop:
SUBS R5, R5, #0x01
BPL loop`,
			[]NewToken{
				{Type: TokenLabel, Literal: "delay:", Line: 1, Column: 1},
				{Type: TokenMOVW, Literal: "MOVW", Line: 2, Column: 1},
				{Type: TokenRegister, Literal: "R5", Line: 2, Column: 6},
				{Type: TokenComma, Literal: ",", Line: 2, Column: 8},
				{Type: TokenImmediate, Literal: "0xFFFF", Line: 2, Column: 10},
				{Type: TokenMOVT, Literal: "MOVT", Line: 3, Column: 1},
				{Type: TokenRegister, Literal: "R5", Line: 3, Column: 6},
				{Type: TokenComma, Literal: ",", Line: 3, Column: 8},
				{Type: TokenImmediate, Literal: "0x000F", Line: 3, Column: 10},
				{Type: TokenLabel, Literal: "loop:", Line: 4, Column: 1},
				{Type: TokenSUBS, Literal: "SUBS", Line: 5, Column: 1},
				{Type: TokenRegister, Literal: "R5", Line: 5, Column: 6},
				{Type: TokenComma, Literal: ",", Line: 5, Column: 8},
				{Type: TokenRegister, Literal: "R5", Line: 5, Column: 10},
				{Type: TokenComma, Literal: ",", Line: 5, Column: 12},
				{Type: TokenImmediate, Literal: "0x01", Line: 5, Column: 14},
				{Type: TokenB, Literal: "B", Line: 6, Column: 1},
				{Type: TokenPL, Literal: "PL", Line: 6, Column: 2},
				{Type: TokenIdentifier, Literal: "loop", Line: 6, Column: 5},
			},
		},
		{
			"BX instruction",
			"BX R14",
			[]NewToken{
				{Type: TokenBX, Literal: "BX", Line: 1, Column: 1},
				{Type: TokenRegister, Literal: "R14", Line: 1, Column: 4},
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
