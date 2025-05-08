package types

// Tokens
type TokenType int

const (
	TokenEOF TokenType = iota
	TokenError
	TokenIdentifier
	TokenRegister
	TokenNumber
	TokenHexNumber

	TokenSBit
	TokenLBit

	TokenComma
	TokenLParen
	TokenRParen
	TokenSemicolon

	Comment
)

type Token struct {
	Type  TokenType
	Value string
	Line  int
	Col   int
}

type Instruction interface {
	ToMachineCode() []byte
}
