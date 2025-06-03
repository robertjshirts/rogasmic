package types

var LiteralToMnemonicToken = map[string]TokenType{
	"movw": TokenMOVW,
	"movt": TokenMOVT,
	"ldr":  TokenLDR,
	"str":  TokenSTR,
	"add":  TokenADD,
	"sub":  TokenSUB,
	"and":  TokenAND,
	"orr":  TokenORR,
	"b":    TokenB,
	"bl":   TokenBL,
	"bx":   TokenBX,
}

var TokenToMnemonic = map[TokenType]MnemonicType{
	TokenMOVW: MnemonicMOVW,
	TokenMOVT: MnemonicMOVT,
	TokenLDR:  MnemonicLDR,
	TokenSTR:  MnemonicSTR,
	TokenADD:  MnemonicADD,
	TokenSUB:  MnemonicSUB,
	TokenAND:  MnemonicAND,
	TokenORR:  MnemonicORR,
	TokenBX:   MnemonicBX,
	TokenB:    MnemonicB,
	TokenBL:   MnemonicBL,
}

var LiteralToCondition = map[string]ConditionType{
	"eq": ConditionEQ,
	"pl": ConditionPL,
	"al": ConditionAL,
	// Only supporting a few conditions for now
}

var ConditionToBits = map[ConditionType]uint32{
	ConditionEQ: 0b0000,
	ConditionAL: 0b1110,
	ConditionPL: 0b0101,
	// "eq": 0b0000,
	// "ne": 0b0001,
	// "cs": 0b0010,
	// "cc": 0b0011,
	// "mi": 0b0100,
	// "pl": 0b0101,
	// "vs": 0b0110,
	// "vc": 0b0111,
	// "hi": 0b1000,
	// "ls": 0b1001,
	// "ge": 0b1010,
	// "lt": 0b1011,
	// "gt": 0b1100,
	// "le": 0b1101,
	// "al": 0b1110,
}

var MnemonicToBits = map[MnemonicType]uint32{
	MnemonicMOVW: 0b0011_0000,
	MnemonicMOVT: 0b0011_0100,
	MnemonicLDR:  0b1,
	MnemonicSTR:  0b0,
	MnemonicADD:  0b0100,
	MnemonicSUB:  0b0010,
	MnemonicAND:  0b0000,
	MnemonicORR:  0b1100,
	MnemonicBX:   0b0001_0010_1111_1111_1111_0001,
	MnemonicB:    0b101,
	MnemonicBL:   0b101,
}

var MnemonicToCategory = map[MnemonicType]MnemonicCategory{
	MnemonicMOVW: MnemonicCategoryMOV,
	MnemonicMOVT: MnemonicCategoryMOV,
	MnemonicLDR:  MnemonicCategoryLoadStore,
	MnemonicSTR:  MnemonicCategoryLoadStore,
	MnemonicADD:  MnemonicCategoryArithmetic,
	MnemonicSUB:  MnemonicCategoryArithmetic,
	MnemonicAND:  MnemonicCategoryArithmetic,
	MnemonicORR:  MnemonicCategoryArithmetic,
	MnemonicBX:   MnemonicCategoryBranchExchange,
	MnemonicB:    MnemonicCategoryBranch,
	MnemonicBL:   MnemonicCategoryBranch,
}

var MnemonicTokenToCategory = map[TokenType]MnemonicCategory{
	TokenMOVW: MnemonicCategoryMOV,
	TokenMOVT: MnemonicCategoryMOV,
	TokenLDR:  MnemonicCategoryLoadStore,
	TokenSTR:  MnemonicCategoryLoadStore,
	TokenADD:  MnemonicCategoryArithmetic,
	TokenSUB:  MnemonicCategoryArithmetic,
	TokenAND:  MnemonicCategoryArithmetic,
	TokenORR:  MnemonicCategoryArithmetic,
	TokenBX:   MnemonicCategoryBranchExchange,
	TokenB:    MnemonicCategoryBranch,
	TokenBL:   MnemonicCategoryBranch,
}