package types

type MnemonicType uint32

const (
	MnemonicMOVW MnemonicType = iota
	MnemonicMOVT
	MnemonicLDR
	MnemonicSTR
	MnemonicADD
	MnemonicSUB
	MnemonicAND
	MnemonicORR
	MnemonicBX
	MnemonicB
	MnemonicBL
)

type MnemonicCategory uint32

const (
	MnemonicCategoryMOV MnemonicCategory = iota
	MnemonicCategoryLoadStore
	MnemonicCategoryArithmetic
	MnemonicCategoryBranch
	MnemonicCategoryBranchExchange
)
