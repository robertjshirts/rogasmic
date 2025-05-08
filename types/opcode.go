package types

// Opcodes
type OpCode uint32

const (
	OpMOVW OpCode = iota // Move bottom half of word
	OpMOVT               // Move top half of word
	OpLDR                // Load into register
	OpSTR                // Save from register
	OpADD                // Add
	OpSUB                // Subtract
	OpAND                // Bitwise AND
	OpORR                // Bitwise OR
	OpBX                 // Branch and exchange
	OpB                  // Branch
)

var OpCodesByLit = map[string]OpCode{
	"MOVW": OpMOVW,
	"MOVT": OpMOVT,
	"LDR":  OpLDR,
	"STR":  OpSTR,
	"ADD":  OpADD,
	"SUB":  OpSUB,
	"AND":  OpAND,
	"ORR":  OpORR,
	"BX":   OpBX,
	"B":    OpB,
}

var OpCodeBits = map[OpCode]uint32{
	OpMOVW: 0b0011_0000,
	OpMOVT: 0b0011_0100,
	OpLDR:  0b1,
	OpSTR:  0b0,
	OpADD:  0b0100,
	OpSUB:  0b0010,
	OpAND:  0b0000,
	OpORR:  0b1100,
	OpBX:   0b0001_0010_1111_1111_1111_0001,
	OpB:    0b101,
}

type InstructionType int

const (
	MOVType InstructionType = iota
	ArithmeticType
	MemoryType
	BranchExType
	BranchType
)

var InstructionTypes = map[OpCode]InstructionType{
	OpMOVT: MOVType,
	OpMOVW: MOVType,
	OpLDR:  MemoryType,
	OpSTR:  MemoryType,
	OpADD:  ArithmeticType,
	OpSUB:  ArithmeticType,
	OpAND:  ArithmeticType,
	OpORR:  ArithmeticType,
	OpBX:   BranchExType,
	OpB:    BranchType,
}
