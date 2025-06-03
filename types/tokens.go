package types

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Col     int
}

type TokenType int

const (
	TokenEOF   TokenType = iota
	TokenError           // Only used when we need to return a token type but there is an error

	TokenComma
	TokenLBracket
	TokenRBracket
	TokenBang // For write back bit on registers (specifically for LDM/STM)
	TokenLBrace
	TokenRBrace
	TokenDash

	TokenIdentifier
	TokenLabel

	TokenRegister
	TokenImmediate

	TokenMOVW
	TokenMOVT
	TokenLDR
	TokenSTR
	TokenLDM
	TokenSTM
	TokenADD
	TokenSUB
	TokenAND
	TokenORR
	TokenBX
	TokenB
	TokenBL

	TokenS

	TokenEQ
	TokenPL
	TokenAL
)

var TokenToLiteral = map[TokenType]string{
	TokenEOF:        "EOF",
	TokenError:      "ERROR",
	TokenComma:      "COMMA",
	TokenLBracket:   "LBRACKET",
	TokenRBracket:   "RBRACKET",
	TokenIdentifier: "IDENTIFIER",
	TokenLabel:      "LABEL",
	TokenRegister:   "REGISTER",
	TokenImmediate:  "IMMEDIATE",
	TokenMOVW:       "MOVW",
	TokenMOVT:       "MOVT",
	TokenLDR:        "LDR",
	TokenSTR:        "STR",
	TokenLDM:        "LDM",
	TokenSTM:        "STM",
	TokenADD:        "ADD",
	TokenSUB:        "SUB",
	TokenAND:        "AND",
	TokenORR:        "ORR",
	TokenBX:         "BX",
	TokenB:          "B",
	TokenBL:         "BL",
}
