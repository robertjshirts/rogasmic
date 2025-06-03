package types

type ConditionType uint32

const (
	ConditionEQ ConditionType = iota
	ConditionNE
	ConditionCS
	ConditionCC
	ConditionMI
	ConditionPL
	ConditionVS
	ConditionVC
	ConditionHI
	ConditionLS
	ConditionGE
	ConditionLT
	ConditionGT
	ConditionLE
	ConditionAL
)