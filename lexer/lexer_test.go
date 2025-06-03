package lexer

import (
	"testing"

	"github.com/robertjshirts/rogasmic/types"
)

func TestLexerMnemonics(t *testing.T) {
	cases := []struct {
		name           string
		input          string
		expectedTokens []types.Token
	}{}

	for k, v := range types.LiteralToMnemonicToken {
		cases = append(cases, struct {
			name           string
			input          string
			expectedTokens []types.Token
		}{
			name:  k + " operation",
			input: k,
			expectedTokens: []types.Token{
				{Type: v, Literal: k, Line: 1, Col: 1},
				{Type: types.TokenEOF, Literal: "", Line: -1, Col: -1},
			},
		})
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := NewLexer(c.input)
			tokens, err := l.Tokenize()
			if err != nil {
				t.Fatalf("unexpected error tokenizing: %v", err)
			}
			if len(tokens) != len(c.expectedTokens) {
				t.Fatalf("expected %d tokens, got %d", len(c.expectedTokens), len(tokens))
			}
			for i, token := range tokens {
				if token.Type != c.expectedTokens[i].Type || token.Literal != c.expectedTokens[i].Literal ||
					token.Line != c.expectedTokens[i].Line || token.Col != c.expectedTokens[i].Col {
					t.Errorf("token mismatch at index %d: expected %+v, got %+v", i, c.expectedTokens[i], token)
				}
			}
		})
	}
}

func TestLexercMnemonicsWithConditions(t *testing.T) {
	cases := []struct {
		name           string
		input          string
		expectedTokens []types.Token
	}{}

	for mk, mv := range types.LiteralToMnemonicToken {
		for ck := range types.LiteralToCondition {
			cases = append(cases, struct {
				name           string
				input          string
				expectedTokens []types.Token
			}{
				name:  mk + " operation with " + ck + " condition",
				input: mk + ck,
				expectedTokens: []types.Token{
					{Type: mv, Literal: mk + ck, Line: 1, Col: 1},
					{Type: types.TokenEOF, Literal: "", Line: -1, Col: -1},
				},
			})
		}
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := NewLexer(c.input)
			tokens, err := l.Tokenize()
			if err != nil {
				t.Fatalf("unexpected error tokenizing: %v", err)
			}
			if len(tokens) != len(c.expectedTokens) {
				t.Fatalf("expected %d tokens, got %d", len(c.expectedTokens), len(tokens))
			}
			for i, token := range tokens {
				if token.Type != c.expectedTokens[i].Type || token.Literal != c.expectedTokens[i].Literal ||
					token.Line != c.expectedTokens[i].Line || token.Col != c.expectedTokens[i].Col {
					t.Errorf("token mismatch at index %d: expected %+v, got %+v", i, c.expectedTokens[i], token)
				}
			}
		})
	}
}

func TestLexerBasics(t *testing.T) {
	cases := []struct {
		name           string
		input          string
		expectedTokens []types.Token
	}{
		{
			name:  "empty input",
			input: "",
			expectedTokens: []types.Token{
				{Type: types.TokenEOF, Literal: "", Line: -1, Col: -1},
			},
		},
		{
			name:  "single register",
			input: "r0",
			expectedTokens: []types.Token{
				{Type: types.TokenRegister, Literal: "r0", Line: 1, Col: 1},
				{Type: types.TokenEOF, Literal: "", Line: -1, Col: -1},
			},
		},
		{
			name:  "immediate value",
			input: "#42",
			expectedTokens: []types.Token{
				{Type: types.TokenImmediate, Literal: "42", Line: 1, Col: 1},
				{Type: types.TokenEOF, Literal: "", Line: -1, Col: -1},
			},
		},
		{
			name:  "label",
			input: "my_label:",
			expectedTokens: []types.Token{
				{Type: types.TokenLabel, Literal: "my_label", Line: 1, Col: 1},
				{Type: types.TokenEOF, Literal: "", Line: -1, Col: -1},
			},
		},
		{
			name:  "comma",
			input: ",",
			expectedTokens: []types.Token{
				{Type: types.TokenComma, Literal: ",", Line: 1, Col: 1},
				{Type: types.TokenEOF, Literal: "", Line: -1, Col: -1},
			},
		},
		{
			name:  "l bracket",
			input: "[",
			expectedTokens: []types.Token{
				{Type: types.TokenLBracket, Literal: "[", Line: 1, Col: 1},
				{Type: types.TokenEOF, Literal: "", Line: -1, Col: -1},
			},
		},
		{
			name:  "r bracket",
			input: "]",
			expectedTokens: []types.Token{
				{Type: types.TokenRBracket, Literal: "]", Line: 1, Col: 1},
				{Type: types.TokenEOF, Literal: "", Line: -1, Col: -1},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := NewLexer(c.input)
			tokens, err := l.Tokenize()
			if err != nil {
				t.Fatalf("unexpected error tokenizing: %v", err)
			}
			if len(tokens) != len(c.expectedTokens) {
				t.Fatalf("expected %d tokens, got %d", len(c.expectedTokens), len(tokens))
			}
			for i, token := range tokens {
				if token.Type != c.expectedTokens[i].Type || token.Literal != c.expectedTokens[i].Literal ||
					token.Line != c.expectedTokens[i].Line || token.Col != c.expectedTokens[i].Col {
					t.Errorf("token mismatch at index %d: expected %+v, got %+v", i, c.expectedTokens[i], token)
				}
			}
		})
	}
}

func TestLexerLines(t *testing.T) {
	cases := []struct {
		name           string
		input          string
		expectedTokens []types.Token
	}{
		{
			name:  "simple movt",
			input: "movt r4, #0x0000",
			expectedTokens: []types.Token{
				{Type: types.TokenMOVT, Literal: "movt", Line: 1, Col: 1},
				{Type: types.TokenRegister, Literal: "r4", Line: 1, Col: 6},
				{Type: types.TokenComma, Literal: ",", Line: 1, Col: 8},
				{Type: types.TokenImmediate, Literal: "0x0000", Line: 1, Col: 10},
				{Type: types.TokenEOF, Literal: "", Line: -1, Col: -1},
			},
		},
		{
			name:  "movt with label",
			input: "my_label:\nmovt r3, #0x0000",
			expectedTokens: []types.Token{
				{Type: types.TokenLabel, Literal: "my_label", Line: 1, Col: 1},
				{Type: types.TokenMOVT, Literal: "movt", Line: 2, Col: 1},
				{Type: types.TokenRegister, Literal: "r3", Line: 2, Col: 6},
				{Type: types.TokenComma, Literal: ",", Line: 2, Col: 8},
				{Type: types.TokenImmediate, Literal: "0x0000", Line: 2, Col: 10},
				{Type: types.TokenEOF, Literal: "", Line: -1, Col: -1},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := NewLexer(c.input)
			tokens, err := l.Tokenize()
			if err != nil {
				t.Fatalf("unexpected error tokenizing: %v", err)
			}
			if len(tokens) != len(c.expectedTokens) {
				t.Fatalf("expected %d tokens, got %d", len(c.expectedTokens), len(tokens))
			}
			for i, token := range tokens {
				if token.Type != c.expectedTokens[i].Type || token.Literal != c.expectedTokens[i].Literal ||
					token.Line != c.expectedTokens[i].Line || token.Col != c.expectedTokens[i].Col {
					t.Errorf("token mismatch at index %d: expected %+v, got %+v", i, c.expectedTokens[i], token)
				}
			}
		})
	}
}

func TestLexerComments(t *testing.T) {
	cases := []struct {
		name           string
		input          string
		expectedTokens []types.Token
	}{
		{
			name:  "only comment",
			input: "; this is a comment",
			expectedTokens: []types.Token{
				{Type: types.TokenEOF, Literal: "", Line: -1, Col: -1},
			},
		},
		{
			name:  "comment then code",
			input: "; comment line\nr1",
			expectedTokens: []types.Token{
				{Type: types.TokenRegister, Literal: "r1", Line: 2, Col: 1},
				{Type: types.TokenEOF, Literal: "", Line: -1, Col: -1},
			},
		},
		{
			name:  "code then comment",
			input: "r2 ; comment after code",
			expectedTokens: []types.Token{
				{Type: types.TokenRegister, Literal: "r2", Line: 1, Col: 1},
				{Type: types.TokenEOF, Literal: "", Line: -1, Col: -1},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := NewLexer(c.input)
			tokens, err := l.Tokenize()
			if err != nil {
				t.Fatalf("unexpected error tokenizing: %v", err)
			}
			if len(tokens) != len(c.expectedTokens) {
				t.Fatalf("expected %d tokens, got %d", len(c.expectedTokens), len(tokens))
			}
			for i, token := range tokens {
				if token.Type != c.expectedTokens[i].Type || token.Literal != c.expectedTokens[i].Literal ||
					token.Line != c.expectedTokens[i].Line || token.Col != c.expectedTokens[i].Col {
					t.Errorf("token mismatch at index %d: expected %+v, got %+v", i, c.expectedTokens[i], token)
				}
			}
		})
	}
}
