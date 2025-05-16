package lexer

type NewTokenType int

const (
	TokenEOF NewTokenType = iota
	TokenError

	TokenComma
	TokenLParen
	TokenRParen

	TokenIdentifier
	TokenLabel

	TokenRegister
	TokenImmediate

	TokenSSuffix
	TokenLSuffix
	TokenConditionCode

	TokenMOVW
	TokenMOVT
	TokenLDR
	TokenSTR
	TokenADD
	TokenSUB
	TokenAND
	TokenORR
	TokenBX
	TokenB
)

var NewTokenTypesByLit = map[string]NewTokenType{
	"MOVW": TokenMOVW,
	"MOVT": TokenMOVT,
	"LDR":  TokenLDR,
	"STR":  TokenSTR,
	"ADD":  TokenADD,
	"SUB":  TokenSUB,
	"AND":  TokenAND,
	"ORR":  TokenORR,
	"BX":   TokenBX,
	"B":    TokenB,
	"AL":   TokenConditionCode,
	"PL":   TokenConditionCode,
	"EQ":   TokenConditionCode,
	"NE":   TokenConditionCode,
	"LT":   TokenConditionCode,
	"LE":   TokenConditionCode,
	"GT":   TokenConditionCode,
	"GE":   TokenConditionCode,
	"S":    TokenSSuffix,
	"L":    TokenLSuffix,
}

type NewToken struct {
	Type    NewTokenType
	Literal string
	Line    int
	Column  int
}
