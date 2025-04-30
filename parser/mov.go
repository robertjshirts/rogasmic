package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/robertjshirts/rogasmic/types"
)

type MOVInstruction struct {
	opcode         types.OpCode
	condition      types.CondCode
	DestRegister   uint32
	ImmediateValue uint32
}

func isValidMOVTokens(tokens []types.Token) bool {
	if len(tokens) < 3 {
		return false
	}

	if tokens[0].Type != types.TokenIdentifier || (tokens[0].Value != "MOVW" && tokens[0].Value != "MOVT") {
		return false
	}

	if tokens[1].Type != types.TokenRegister {
		return false
	}

	if tokens[2].Type != types.TokenNumber && tokens[2].Type != types.TokenHexNumber {
		return false
	}

	return true
}

func NewMOVInstruction(tokens []types.Token) (*MOVInstruction, error) {
	if !isValidMOVTokens(tokens) {
		return nil, fmt.Errorf("invalid MOV instruction tokens: %v", tokens)
	}

	opcode := types.MOVW
	if tokens[0].Value == "MOVT" {
		opcode = types.MOVT
	}



	instruction := &MOVInstruction{
		opcode:         types.MOVW,
		condition:      types.CondAL,
		DestRegister:   (tokens[1].Value[1:]),
		ImmediateValue: 0,
	}
	if tokens[0].Value == "MOVT" {
		instruction.opcode = types.MOVT
	} 

	



	return instruction, nil
}

func (i *MOVInstruction) Parse(operands []string) error {
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
	// binary |= CondBits[i.condition] << 28
	binary |= 3 << 24 // 3 bits for MOV opcode
	// binary |= OpBits[i.opcode] << 20
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
