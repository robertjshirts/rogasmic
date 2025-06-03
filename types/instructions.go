package types

type Instruction interface {
	ToMachineCode(labels map[string]uint32) ([]byte, error)
}

type LabelMap map[string]uint32
