package lexer

type NewTokenType int

const (
	TokenEOF NewTokenType = iota
	TokenError

	TokenComma
	TokenLBracket
	TokenRBracket

	TokenIdentifier
	TokenLabel

	TokenRegister
	TokenImmediate

	TokenSSuffix
	TokenLSuffix

	TokenMOVW
	TokenMOVT
	TokenLDR
	TokenSTR
	TokenADD
	TokenSUB
	TokenSUBS
	TokenAND
	TokenORR
	TokenBX
	TokenB
	TokenBL

	TokenEQ
	TokenPL
	TokenAL
)

var MnemonicsByLit = map[string]NewTokenType{
	"MOVW": TokenMOVW,
	"MOVT": TokenMOVT,
	"LDR":  TokenLDR,
	"STR":  TokenSTR,
	"ADD":  TokenADD,
	"SUB":  TokenSUB,
	"SUBS": TokenSUBS,
	"AND":  TokenAND,
	"ORR":  TokenORR,
	"BX":   TokenBX,
	"B":    TokenB,
	"BL":   TokenBL,
}

var ConditionCodesByLit = map[string]NewTokenType{
	"AL": TokenAL,
	"PL": TokenPL,
	"EQ": TokenEQ,
}

var NewTokenTypesByLit = map[string]NewTokenType{
	"MOVW": TokenMOVW,
	"MOVT": TokenMOVT,
	"LDR":  TokenLDR,
	"STR":  TokenSTR,
	"ADD":  TokenADD,
	"SUB":  TokenSUB,
	"SUBS": TokenSUBS,
	"AND":  TokenAND,
	"ORR":  TokenORR,
	"BX":   TokenBX,
	"B":    TokenB,
	"BL":   TokenBL,
	"AL":   TokenAL,
	"PL":   TokenPL,
	"EQ":   TokenEQ,
	"S":    TokenSSuffix,
	"L":    TokenLSuffix,
}

var NewTokenTypeInstructions = map[NewTokenType]struct{}{
	TokenMOVW: {},
	TokenMOVT: {},
	TokenLDR:  {},
	TokenSTR:  {},
	TokenADD:  {},
	TokenSUB:  {},
	TokenSUBS: {},
	TokenAND:  {},
	TokenORR:  {},
	TokenBX:   {},
	TokenB:    {},
	TokenBL:   {},
}

type NewToken struct {
	Type    NewTokenType
	Literal string
	Line    int
	Column  int
}
