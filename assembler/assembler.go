package assembler

import (
	"fmt"

	"github.com/robertjshirts/rogasmic/types"
)

// Assembler is responsible for converting parsed instructions into machine code.
type Assembler struct {
	instructions []types.Instruction
	labels       map[string]uint32 // Maps label names to instruction numbers
}

func NewAssembler(instructions []types.Instruction, labels map[string]uint32) *Assembler {
	return &Assembler{
		instructions: instructions,
		labels:       labels,
	}
}

func (a *Assembler) Assemble() ([]byte, error) {
	var machineCode []byte

	for _, instruction := range a.instructions {
		code, err := instruction.ToMachineCode(a.labels)
		if err != nil {
			return nil, fmt.Errorf("error converting instruction %T to machine code: %w", instruction, err)
		}
		machineCode = append(machineCode, code...)
	}

	return machineCode, nil
}
