package types

type MnemonicType uint32

const (
	MnemonicMOVW MnemonicType = iota
	MnemonicMOVT
	MnemonicLDR
	MnemonicSTR
	MnemonicLDM
	MnemonicSTM
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
	MnemonicCategoryLoadStoreMultiple
	MnemonicCategoryArithmetic
	MnemonicCategoryBranch
	MnemonicCategoryBranchExchange
)
