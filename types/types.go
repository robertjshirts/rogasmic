package types

// Tokens
type TokenType int

const (
	TokenEOF TokenType = iota
	TokenError
	TokenIdentifier
	TokenRegister
	TokenSBit
	TokenNumber
	TokenHexNumber

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
	OpMOVW OpCode = iota // Move bottom half of word
	OpMOVT               // Move top half of word
	OpLDR                // Load into register
	OpSTR                // Save from register
	OpADD                // Add
	OpSUB                // Subtract
	OpORR                // Bitwise OR
	OpB                  // Branch
)

var OpCodeMeta = [...]struct {
	Name string
	Bits uint32
}{
	OpMOVW: {"MOVW", 0b0010},
	OpMOVT: {"MOVT", 0b0010},
}

var (
	OpCodeByName map[string]OpCode
	CondByName   map[string]CondCode
)

type Instruction interface {
	Parse(tokens []Token) error
	ToMachineCode() []byte
}

var TokenMeta = [...]struct {
	Name string
}{
	TokenEOF:        {"EOF"},
	TokenError:      {"Error"},
	TokenIdentifier: {"Identifier"},
	TokenRegister:   {"Register"},
	TokenSBit:       {"SBit"},
	TokenNumber:     {"Number"},
	TokenHexNumber:  {"HexNumber"},
	TokenComma:      {"Comma"},
	TokenLParen:     {"LParen"},
	TokenRParen:     {"RParen"},
	TokenSemicolon:  {"Semicolon"},
	Comment:         {"Comment"},
}

func Init() {
	OpCodeByName = make(map[string]OpCode, len(OpCodeMeta))
	for code, meta := range OpCodeMeta {
		OpCodeByName[meta.Name] = OpCode(code)
	}

	CondByName = make(map[string]CondCode, len(CondMeta))
	for cond, meta := range CondMeta {
		CondByName[meta.Name] = CondCode(cond)
	}
}
