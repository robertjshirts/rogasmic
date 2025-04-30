package main

import (
	"fmt"
	"strconv"
	"strings"
	r o
)

type MOVInstruction struct {
	opcode         OpCode
	condition      CondCode
	DestRegister   uint32
	ImmediateValue uint32
}

func (i *MOVInstruction) GetOperation() OpCode {
	return i.opcode
}

func (i *MOVInstruction) GetCondition() CondCode {
	return i.condition
}

func NewMOVInstruction(opcode OpCode, condition CondCode) (*MOVInstruction, error) {
	if opcode != MOVW && opcode != MOVT {
		return nil, fmt.Errorf("invalid opcode: %s", opcode)
	}

	return &MOVInstruction{
		opcode:    opcode,
		condition: condition,
	}, nil
}

func (i *MOVInstruction) GetOperands(base int) []string {
	return []string{strconv.FormatInt(int64(i.DestRegister), base), strconv.FormatInt(int64(i.ImmediateValue), base)}
}

// Takes the assembly line, sans operator and condition code, and parses it into the struct
// Operands should not have any parenthesis, commas, spaces, or comments, or the MOV instruction
func (i *MOVInstruction) ParseOperands(operands []string) error {
	if len(operands) != 2 {
		return fmt.Errorf("invalid operands: %s", operands)
	}

	register, err := strconv.ParseInt(strings.ReplaceAll(operands[0], "R", ""), 10, 8)
	if err != nil || register < 0 || register > 15 {
		return fmt.Errorf("invalid register: %s", operands[0])
	}
	i.DestRegister = uint32(register)

	immediate, err := strconv.ParseInt(operands[1], 10, 16)
	if err != nil {
		return fmt.Errorf("invalid immediate value: %s", operands[1])
	}
	i.ImmediateValue = uint32(immediate)

	return nil
}

func (i *MOVInstruction) ToMachineCode() []byte {
	var binary uint32
	binary |= CondBits[i.condition] << 28
	binary |= 3 << 24 // 3 bits for MOV opcode
	binary |= OpBits[i.opcode] << 20
	binary |= (i.ImmediateValue >> 12) & 0x000F // Top 4 bits of immediate
	binary |= i.DestRegister << 12
	binary |= i.ImmediateValue & 0x0FFF // Bottom 12 bits of immediate
	machineCode := make([]byte, 4)
	machineCode[0] = byte((binary >> 24) & 0xFF)
	machineCode[1] = byte((binary >> 16) & 0xFF)
	machineCode[2] = byte((binary >> 8) & 0xFF)
	machineCode[3] = byte(binary & 0xFF)
	fmt.Printf("%02X %02X %02X %02X\n", machineCode[0], machineCode[1], machineCode[2], machineCode[3])
	return machineCode
}

// Expected output for MOVW R4, 0
// Binary 1110 0011 0000 0000 0100 0000 0000 0000
// Hex E3 00 40 00
// Little endian hex 00 40 00 E3
