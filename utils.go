package main

func IsValidOpCode(opcode string) bool {
	switch opcode {
	case string(MOVW), string(MOVT), string(LDR), string(STR), string(ADD), string(SUB), string(ORR), string(B):
		return true
	default:
		return false
	}
}

func IsValidCondCode(cond string) bool {
	switch cond {
	case string(AL), string(PL):
		return true
	default:
		return false
	}
}
