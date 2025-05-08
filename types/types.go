package types

// Tokens
type TokenType int

const (
	// Control tokens
	TokenEOF TokenType = iota
	TokenError

	// Literals and identifiers
	TokenIdentifier
	TokenRegister
	TokenNumber
	TokenHexNumber

	// ARM-specific bits
	TokenSBit
	TokenLBit

	// Punctuation
	TokenComma
	TokenLParen
	TokenRParen
	TokenSemicolon
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
