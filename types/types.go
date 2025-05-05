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

var CondCodeBits = map[string]uint32{
	"AL": 0b1110,
	"PL": 0b1010,
}

var OpCodeMeta2 = [...]struct {
	Name string
	Bits uint32
}{
	OpMOVW: {"MOVW", 0b0010},
	OpMOVT: {"MOVT", 0b0010},
	OpLDR:  {"LDR", 0b0100},
	OpSTR:  {"STR", 0b0101},
	OpADD:  {"ADD", 0b0011},
	OpSUB:  {"SUB", 0b0110},
	OpORR:  {"ORR", 0b0111},
	OpB:    {"B", 0b1000},
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

/**
0000 = AND
0010 = SUB
0011 = RSB
0100 = ADD
0101 = ADC
0110 = SBC
0111 = RSC
1000 = TST
1001 = TEQ
1010 = CMP
1011 = CMN
1100 = ORR
1101 = MOV
1110 = BIC
1111 = MVN

**/

var OpCodeBits = map[string]uint32{
	"MOVW": 0b0000,
	"MOVT": 0b0100,
	"AND":  0b0000,
	"EOR":  0b0001,
	"SUB":  0b0010,
	"RSB":  0b0011,
	"ADD":  0b0100,
	"ADC":  0b0101,
	"SBC":  0b0110,
	"RSC":  0b0111,
	"TST":  0b1000,
	"TEQ":  0b1001,
	"CMP":  0b1010,
	"CMN":  0b1011,
	"ORR":  0b1100,
	"BIC":  0b1110,
	"MOV":  0b1101,
	"MVN":  0b1111,
	"LDR":  0b1,
	"STR":  0b0,
	"B":    0b101,
}

var (
	OpCodeByName map[string]OpCode
	CondByName   map[string]CondCode
)
var MovSet = map[string]struct{}{
	"MOVW": {},
	"MOVT": {},
}
var ArithmeticSet = map[string]struct{}{
	"ADD": {},
	"SUB": {},
	"ORR": {},
}
var MemorySet = map[string]struct{}{
	"LDR": {},
	"STR": {},
}
var BranchSet = map[string]struct{}{
	"B": {},
}

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
