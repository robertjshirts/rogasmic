package utils

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/robertjshirts/rogasmic/types"
)

// IsLiteralChar checks if a byte is a valid literal char. Returns true on letter, digit, and underscore
func IsLiteralChar(c byte) bool {
	return unicode.IsLetter(rune(c)) || unicode.IsDigit(rune(c)) || c == '_'
}

// IsImmediate checks if a string is a valid decimal or hexadecimal value. Hex values should start with "0x" or "0X".
func IsImmediate(lit string) bool {
	validChars := unicode.Digit

	isHex := strings.HasPrefix(lit, "0X") || strings.HasPrefix(lit, "0x")
	if isHex {
		lit = lit[2:] // trim 0x
		validChars = unicode.Hex_Digit
	}

	for _, c := range lit {
		if !unicode.Is(validChars, c) {
			return false
		}
	}

	return true
}

// IsRegister checks if a string is a valid register name. Registers should start with 'r' or 'R' followed by digits.
func IsRegister(lit string) bool {
	if strings.ToLower(lit) == "sp" || strings.ToLower(lit) == "lr" || strings.ToLower(lit) == "pc" {
		return true
	}

	if len(lit) < 2 || (lit[0] != 'r' && lit[0] != 'R') {
		return false
	}

	register, err := strconv.Atoi(lit[1:])
	if err != nil || register < 0 || register > 15 {
		return false
	}

	return true
}

func NormalizeRegister(lit string) string {
	if strings.ToLower(lit) == "sp" {
		return "r13"
	}
	if strings.ToLower(lit) == "lr" {
		return "r14"
	}
	if strings.ToLower(lit) == "pc" {
		return "r15"
	}
	return lit
}

/*
IsOperation checks if a string begins with a valid operation.
*/
func IsOperation(lit string) bool {
	typ := GetMnemonicTokenType(lit)
	return typ != types.TokenError
}

/*
GetMnemonicTokenType returns the TokenType for an operation literal. Matches the longest mnemonic first.
*/
func GetMnemonicTokenType(lit string) types.TokenType {
	if len(lit) == 0 {
		return types.TokenError
	}

	lit = strings.ToLower(lit)

	// Iterate backwards through the literal to find longest operation
	for i := len(lit); i > 0; i-- {
		if tokenType, ok := types.LiteralToMnemonicToken[lit[0:i]]; ok {
			return tokenType
		}
	}

	// We didn't find an operation
	return types.TokenError
}
