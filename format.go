package main

import (
	"strings"
)

// Separate the operation, condition code, and S bit with white space
// Remove all parenthesis, commas, and comments (; delimited)
func FormatLine(line string) string {
	strings.Trim(line, " ")
	var formatted_line strings.Builder
	// var previous string
	var current strings.Builder
	for _, char := range line {
		switch char {
		case ';':
			return strings.Trim(formatted_line.String(), " ")
		case ',':
		case '(':
		case ')':
			break
		case ' ':
			// TODO later
			continue
		default:
			current.WriteRune(char)
			// Put space between operation and condition code
			if IsValidOpCode(current.String()) || IsValidCondCode(current.String()) {

			}

		}
		// Stop reading line when we reach a comment
		if char == ';' {
			break
		}

		// Skip parenthesis
		if char == '(' || char == ')' {
			continue
		}

		if char == ' ' {
			formatted_line.WriteString(current.String())
			formatted_line.WriteString(" ")
			current.Reset()
		} else {
			current.WriteRune(char)
		}
	}

	return strings.Trim(formatted_line.String(), " ")
}
