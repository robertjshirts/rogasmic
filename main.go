package main

import (
	"fmt"
	"os"

	"github.com/robertjshirts/rogasmic/assembler"
	"github.com/robertjshirts/rogasmic/lexer"
	"github.com/robertjshirts/rogasmic/parser"
	"github.com/robertjshirts/rogasmic/types"
)

func main() {
	inputFile := "labels.asm"
	if len(os.Args) > 1 {
		inputFile = os.Args[1]
	}
	file, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	l := lexer.NewLexer(string(file))
	tokens, err := l.Tokenize()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Tokenized %d tokens:\n", len(tokens))
	for _, token := range tokens {
		fmt.Printf("Token Type=%s, Literal=%s, Line=%d, Col=%d\n", types.TokenToLiteral[token.Type], token.Literal, token.Line, token.Col)
	}

	p := parser.NewParser(tokens)
	instructions, labelMap, err := p.Parse()
	if err != nil {
		fmt.Printf("Error parsing instructions: %v\n", err)
		return
	}

	fmt.Printf("Parsed %d instructions:\n", len(instructions))

	a := assembler.NewAssembler(instructions, labelMap)
	machineCode, err := a.Assemble()
	if err != nil {
		fmt.Printf("Error assembling instructions: %v\n", err)
		return
	}
	fmt.Printf("Assembled machine code: %x\n", machineCode)

	fmt.Printf("Writing kernel7.img...\n")
	outputFile := "kernel7.img"
	if len(os.Args) > 2 {
		outputFile = os.Args[2]
	}
	err = os.WriteFile(outputFile, machineCode, 0644)
	if err != nil {
		fmt.Printf("Error writing output file %s: %v\n", outputFile, err)
		return
	}
	fmt.Printf("Successfully wrote %d bytes to %s\n", len(machineCode), outputFile)
	fmt.Printf("Done!\n")
}
