package parser

import (
	"bytes"
	"testing"

	"github.com/robertjshirts/rogasmic/lexer"
)

func TestParserMOV(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		expected      [][]byte
		expectedError bool
	}{
		{
			name:          "MOVW with register and immediate",
			input:         "MOVW R5, #0xFFFF",
			expected:      [][]byte{{0xFF, 0x5F, 0x0F, 0xE3}},
			expectedError: false,
		},
		{
			name:          "MOVT with register and immediate",
			input:         "MOVT R4, #0x3F20",
			expected:      [][]byte{{0x20, 0x4F, 0x43, 0xE3}},
			expectedError: false,
		},
		{
			name:          "MOVW with condition",
			input:         "MOVWPL R5, #0xFFFF",
			expected:      [][]byte{{0xFF, 0x5F, 0x0F, 0x53}},
			expectedError: false,
		},
		{
			name:          "MOVT with invalid register",
			input:         "MOVT R16, #128",
			expected:      [][]byte{},
			expectedError: true,
		},
		{
			name:          "MOVW with missing immediate",
			input:         "MOVW R3,",
			expected:      [][]byte{},
			expectedError: true,
		},
		{
			name:          "MOVT with invalid condition",
			input:         "MOVTXYZ R4, #64",
			expected:      [][]byte{},
			expectedError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := lexer.NewLexer(c.input)
			toks, err := l.Tokenize()
			if err != nil {
				t.Fatalf("unexpected error tokenizing: %v", err)
			}
			p := NewParser(toks)
			instructions, _, err := p.Parse()
			if c.expectedError {
				if err == nil {
					t.Fatalf("expected error but got none for input: %s", c.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error parsing: %v", err)
			}

			if len(instructions) != len(c.expected) {
				t.Fatalf("expected %d instructions, got %d", len(c.expected), len(instructions))
			}

			for i, inst := range instructions {
				machineCode, err := inst.ToMachineCode(map[string]uint32{})
				if err != nil {
					t.Fatalf("unexpected error converting instruction to machine code: %v", err)
				}
				if !bytes.Equal(machineCode, c.expected[i]) {
					t.Errorf("instruction mismatch at index %d: expected %v, got %v", i, c.expected[i], machineCode)
				}
			}
		})
	}
}

func TestParserMemory(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		expected      [][]byte
		expectedError bool
	}{
		{
			name:          "LDR with base register",
			input:         "LDR R3, [R2]",
			expected:      [][]byte{{0x00, 0x30, 0x12, 0xE4}},
			expectedError: false,
		},
		{
			name:          "STR with base register",
			input:         "STR R3, [R2]",
			expected:      [][]byte{{0x00, 0x30, 0x02, 0xE4}},
			expectedError: false,
		},
		{
			name:          "STR with invalid register",
			input:         "STR R16, [R1]",
			expected:      [][]byte{},
			expectedError: true,
		},
		{
			name:          "LDR with missing base register",
			input:         "LDR R6, []",
			expected:      [][]byte{},
			expectedError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := lexer.NewLexer(c.input)
			toks, err := l.Tokenize()
			if err != nil {
				t.Fatalf("unexpected error tokenizing: %v", err)
			}
			p := NewParser(toks)
			instructions, _, err := p.Parse()
			if c.expectedError {
				if err == nil {
					t.Fatalf("expected error but got none for input: %s", c.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error parsing: %v", err)
			}

			if len(instructions) != len(c.expected) {
				t.Fatalf("expected %d instructions, got %d", len(c.expected), len(instructions))
			}

			for i, inst := range instructions {
				machineCode, err := inst.ToMachineCode(map[string]uint32{})
				if err != nil {
					t.Fatalf("unexpected error converting instruction to machine code: %v", err)
				}
				if !bytes.Equal(machineCode, c.expected[i]) {
					t.Errorf("instruction mismatch at index %d: expected %v, got %v", i, c.expected[i], machineCode)
				}
			}
		})
	}
}

func TestParserArithmetic(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		expected      [][]byte
		expectedError bool
	}{
		{
			name:          "ADD with registers and immediate",
			input:         "ADD R3, R4, #0x1C",
			expected:      [][]byte{{0x1C, 0x30, 0x84, 0xE2}},
			expectedError: false,
		},
		{
			name:          "SUB with S bit",
			input:         "SUBS R5, R5, #0x01",
			expected:      [][]byte{{0x01, 0x50, 0x55, 0xE2}},
			expectedError: false,
		},
		{
			name:          "SUB with S bit and condition",
			input:         "SUBSAL R5, R5, #0x01",
			expected:      [][]byte{{0x01, 0x50, 0x55, 0xE2}},
			expectedError: false,
		},
		{
			name:          "ORR with registers and immediate",
			input:         "ORR R3, R3, #0x8",
			expected:      [][]byte{{0x08, 0x30, 0x83, 0xE3}},
			expectedError: false,
		},
		{
			name:          "ADD with invalid register",
			input:         "ADD R16, R4, #0x1C",
			expected:      [][]byte{},
			expectedError: true,
		},
		{
			name:          "SUB with missing immediate",
			input:         "SUB R5, R5,",
			expected:      [][]byte{},
			expectedError: true,
		},
		{
			name:          "ORR with invalid suffix bit",
			input:         "ORRT R3, R3, #0x8",
			expected:      [][]byte{},
			expectedError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := lexer.NewLexer(c.input)
			toks, err := l.Tokenize()
			if err != nil {
				t.Fatalf("unexpected error tokenizing: %v", err)
			}
			p := NewParser(toks)
			instructions, _, err := p.Parse()
			if c.expectedError {
				if err == nil {
					t.Fatalf("expected error but got none for input: %s", c.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error parsing: %v", err)
			}

			if len(instructions) != len(c.expected) {
				t.Fatalf("expected %d instructions, got %d", len(c.expected), len(instructions))
			}

			for i, inst := range instructions {
				machineCode, err := inst.ToMachineCode(map[string]uint32{})
				if err != nil {
					t.Fatalf("unexpected error converting instruction to machine code: %v", err)
				}
				if !bytes.Equal(machineCode, c.expected[i]) {
					t.Errorf("instruction mismatch at index %d: expected %v, got %v", i, c.expected[i], machineCode)
				}
			}
		})
	}
}

func TestParserBranch(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		expected      [][]byte
		expectedError bool
	}{
		{
			name:          "B with immediate",
			input:         "B #0xFFFFEE",
			expected:      [][]byte{{0xEE, 0xFF, 0xFF, 0xEA}},
			expectedError: false,
		},
		{
			name:          "BL with immediate",
			input:         "BL #0xFFFFEE",
			expected:      [][]byte{{0xEE, 0xFF, 0xFF, 0xEB}},
			expectedError: false,
		},
		{
			name:          "B with condition and immediate",
			input:         "BAL #0xFFFFEE",
			expected:      [][]byte{{0xEE, 0xFF, 0xFF, 0xEA}},
			expectedError: false,
		},
		{
			name:          "B with register (invalid)",
			input:         "B R0",
			expected:      [][]byte{},
			expectedError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := lexer.NewLexer(c.input)
			toks, err := l.Tokenize()
			if err != nil {
				t.Fatalf("unexpected error tokenizing: %v", err)
			}
			p := NewParser(toks)
			instructions, _, err := p.Parse()
			if c.expectedError {
				if err == nil {
					t.Fatalf("expected error but got none for input: %s", c.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error parsing: %v", err)
			}

			if len(instructions) != len(c.expected) {
				t.Fatalf("expected %d instructions, got %d", len(c.expected), len(instructions))
			}

			for i, inst := range instructions {
				machineCode, err := inst.ToMachineCode(map[string]uint32{})
				if err != nil {
					t.Fatalf("unexpected error converting instruction to machine code: %v", err)
				}
				if !bytes.Equal(machineCode, c.expected[i]) {
					t.Errorf("instruction mismatch at index %d: expected %v, got %v", i, c.expected[i], machineCode)
				}
			}
		})
	}
}

func TestParserBranchExchange(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		expected      [][]byte
		expectedError bool
	}{
		{
			name:          "BX with register",
			input:         "BX R14",
			expected:      [][]byte{{0x1E, 0xFF, 0x2F, 0xE1}},
			expectedError: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := lexer.NewLexer(c.input)
			toks, err := l.Tokenize()
			if err != nil {
				t.Fatalf("unexpected error tokenizing: %v", err)
			}
			p := NewParser(toks)
			instructions, _, err := p.Parse()
			if c.expectedError {
				if err == nil {
					t.Fatalf("expected error but got none for input: %s", c.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error parsing: %v", err)
			}

			if len(instructions) != len(c.expected) {
				t.Fatalf("expected %d instructions, got %d", len(c.expected), len(instructions))
			}

			for i, inst := range instructions {
				machineCode, err := inst.ToMachineCode(map[string]uint32{})
				if err != nil {
					t.Fatalf("unexpected error converting instruction to machine code: %v", err)
				}
				if !bytes.Equal(machineCode, c.expected[i]) {
					t.Errorf("instruction mismatch at index %d: expected %v, got %v", i, c.expected[i], machineCode)
				}
			}
		})
	}
}

func TestParserLabels(t *testing.T) {
	cases := []struct {
		name string
		input string
		expected [][]byte
		expectedError bool
	}{
		{
			name: "Looping label",
			input: "loop:\nB loop",
			expected: [][]byte{{0xFE, 0xFF, 0xFF, 0xEA}},
			expectedError: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := lexer.NewLexer(c.input)
			toks, err := l.Tokenize()
			if err != nil {
				t.Fatalf("unexpected error tokenizing: %v", err)
			}
			p := NewParser(toks)
			instructions, labelMap, err := p.Parse()
			if c.expectedError {
				if err == nil {
					t.Fatalf("expected error but got none for input: %s", c.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error parsing: %v", err)
			}

			if len(instructions) != len(c.expected) {
				t.Fatalf("expected %d instructions, got %d", len(c.expected), len(instructions))
			}

			for i, inst := range instructions {
				machineCode, err := inst.ToMachineCode(labelMap)
				if err != nil {
					t.Fatalf("unexpected error converting instruction to machine code: %v", err)
				}
				if !bytes.Equal(machineCode, c.expected[i]) {
					t.Errorf("instruction mismatch at index %d: expected %v, got %v", i, c.expected[i], machineCode)
				}
			}
		})
	}
}