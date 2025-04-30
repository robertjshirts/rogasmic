package main

// Tokens
type TokenType int

const (
	TokenEOF TokenType = iota
	TokenError
	TokenIdentifier
	TokenSBit
	TokenNumber

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

type CondCode uint8

const (
	CondAL CondCode = iota
	CondPL
)

var CondMeta = [...]struct {
	Name string
	Bits uint32
}{
	CondAL: {"AL", 0b1110},
	CondPL: {"PL", 0b0101},
}

// Opcodes
type OpCode uint32

const (
	MOVW OpCode = iota // Move bottom half of word
	MOVT               // Move top half of word
	LDR                // Load into register
	STR                // Save from register
	ADD                // Add
	SUB                // Subtract
	ORR                // Bitwise OR
	B                  // Branch
)

var OpCodeMeta = [...]struct {
	Name string
	Bits uint32
}{
	MOVW: {"MOVW", 0b0010},
	MOVT: {"MOVT", 0b0010},
}

type Instruction interface {
	GetOperation() OpCode
	GetCondition() CondCode
	ParseOperands(string) error
	ToMachineCode() []byte
}
