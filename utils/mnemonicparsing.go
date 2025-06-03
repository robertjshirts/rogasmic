package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/robertjshirts/rogasmic/types"
)

func ParseRegister(registerLiteral string) (uint32, error) {
	if len(registerLiteral) < 2 || len(registerLiteral) > 3 || (registerLiteral[0] != 'R' && registerLiteral[0] != 'r') {
		return 0, fmt.Errorf("invalid register format: %s", registerLiteral)
	}

	reg, err := strconv.ParseUint(registerLiteral[1:], 0, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid register number: %s", registerLiteral[1:])
	}
	if reg > 15 {
		return 0, fmt.Errorf("register number out of range: %s", registerLiteral[1:])
	}
	return uint32(reg), nil
}

func ParseImmediate(immediateLiteral string) (uint32, error) {
	value, err := strconv.ParseUint(immediateLiteral, 0, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid immediate value: %s", immediateLiteral[1:])
	}
	return uint32(value), nil
}

func ParseMOVSuffixes(mnemonicLiteral string) (types.ConditionType, error) {
	// Allows MOVT and MOVT w/ condition
	if len(mnemonicLiteral) < 3 || len(mnemonicLiteral) > 6 {
		return types.ConditionAL, fmt.Errorf("invalid MOV mnemonic length: %s", mnemonicLiteral)
	}
	mnemonicLiteral = strings.ToLower(mnemonicLiteral[4:]) // Remove MOVT/MOVW
	if mnemonicLiteral == "" {
		return types.ConditionAL, nil
	}
	condition, ok := types.LiteralToCondition[mnemonicLiteral]
	if !ok {
		return types.ConditionAL, fmt.Errorf("invalid MOV condition: %s", mnemonicLiteral)
	}

	return condition, nil
}

// For now we return false on all suffixes
// Returns condition, pBit, uBit, error
func ParseMemorySuffixes(mnemonicLiteral string) (types.ConditionType, uint32, uint32, error) {
	if len(mnemonicLiteral) < 3 || len(mnemonicLiteral) > 8 {
		return types.ConditionAL, 0, 0, fmt.Errorf("invalid memory mnemonic length: %s", mnemonicLiteral)
	}

	var pBit, uBit uint32
	if mnemonicLiteral[0] == 's' || mnemonicLiteral[0] == 'S' {
		uBit = 1 // add offset for store, subtract for load
	} else {
		pBit = 1 // preIndexed for LDR
	}

	mnemonicLiteral = strings.ToLower(mnemonicLiteral[3:]) // Remove LDR/STR/LDM/STM
	if mnemonicLiteral == "" {
		return types.ConditionAL, pBit, uBit, nil
	}

	// Parse condition
	condition, ok := types.LiteralToCondition[mnemonicLiteral]
	if !ok {
		return types.ConditionAL, pBit, uBit, fmt.Errorf("invalid memory condition: %s", mnemonicLiteral)
	}

	return condition, pBit, uBit, nil
}

func ParseArithmeticSuffixes(mnemonicLiteral string) (types.ConditionType, uint32, error) {
	if len(mnemonicLiteral) < 3 || len(mnemonicLiteral) > 7 {
		return types.ConditionAL, 0, fmt.Errorf("invalid arithmetic mnemonic length: %s", mnemonicLiteral)
	}

	mnemonicLiteral = strings.ToLower(mnemonicLiteral[3:]) // Remove ADD/SUB/ORR/AND
	if mnemonicLiteral == "" {
		return types.ConditionAL, 0, nil
	}

	var sBit uint32
	if len(mnemonicLiteral) > 0 {
		if mnemonicLiteral[0] == 's' {
			sBit = 1
			mnemonicLiteral = mnemonicLiteral[1:] // Remove 's' if present
		}
	}

	if len(mnemonicLiteral) == 0 {
		return types.ConditionAL, sBit, nil // No condition specified, default to AL
	}

	// Parse condition
	condition, ok := types.LiteralToCondition[mnemonicLiteral]
	if !ok {
		return types.ConditionAL, 0, fmt.Errorf("invalid arithmetic condition: %s", mnemonicLiteral)
	}

	return condition, sBit, nil
}

func ParseBranchSuffixes(mnemonicLiteral string) (types.ConditionType, uint32, error) {
	if len(mnemonicLiteral) < 1 || len(mnemonicLiteral) > 4 {
		return types.ConditionAL, 0, fmt.Errorf("invalid branch mnemonic length: %s", mnemonicLiteral)
	}
	mnemonicLiteral = strings.ToLower(mnemonicLiteral[1:]) // Remove the B

	var lBit uint32
	if len(mnemonicLiteral) > 0 {
		if mnemonicLiteral[0] == 'l' {
			lBit = 1
			mnemonicLiteral = mnemonicLiteral[1:]
		}
	}

	if len(mnemonicLiteral) == 0 {
		return types.ConditionAL, lBit, nil // No condition specified, default to AL
	}

	// Parse condition
	condition, ok := types.LiteralToCondition[mnemonicLiteral]
	if !ok {
		return types.ConditionAL, 0, fmt.Errorf("invalid branch condition: %s", mnemonicLiteral)
	}

	return condition, lBit, nil
}

func ParseBranchExchangeSuffixes(mnemonicLiteral string) (types.ConditionType, error) {
	if len(mnemonicLiteral) < 2 || len(mnemonicLiteral) > 4 {
		return types.ConditionAL, fmt.Errorf("invalid branch exchange mnemonic length: %s", mnemonicLiteral)
	}
	mnemonicLiteral = strings.ToLower(mnemonicLiteral[2:]) // Remove BX

	if len(mnemonicLiteral) == 0 {
		return types.ConditionAL, nil // No condition specified, default to AL
	}

	// Parse condition
	condition, ok := types.LiteralToCondition[mnemonicLiteral]
	if !ok {
		return types.ConditionAL, fmt.Errorf("invalid branch exchange condition: %s", mnemonicLiteral)
	}

	return condition, nil
}

// Little endian
func BitsToBytes(bits uint32) []byte {
	bytes := make([]byte, 4)
	bytes[0] = byte(bits & 0xFF)
	bytes[1] = byte((bits >> 8) & 0xFF)
	bytes[2] = byte((bits >> 16) & 0xFF)
	bytes[3] = byte((bits >> 24) & 0xFF)
	return bytes
}

// Direct, not little endian
func BytesToBits(bytes []byte) uint32 {
	if len(bytes) != 4 {
		panic("bytes slice must have exactly 4 elements")
	}

	bits := uint32(bytes[3]) |
		(uint32(bytes[2]) << 8) |
		(uint32(bytes[1]) << 16) |
		(uint32(bytes[0]) << 24)

	return bits
}
