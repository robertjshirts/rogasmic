// 32 bit arm assembler
package main

import (
	"bufio"
	"log"
	"os"

	"github.com/robertjshirts/rogasmic/lexer"
	"github.com/robertjshirts/rogasmic/parser"
)

func main() {
	//var instructions []types.Instruction
	if len(os.Args) > 3 || len(os.Args) < 2 {
		log.Fatal("Usage: rogasmic <input file> <output file>")
	}

	inputFile := os.Args[1]
	if inputFile == "" {
		log.Fatal("Please provide an input filename")
	}

	outputFile := "kernel7.img"
	if len(os.Args) == 3 {
		outputFile = os.Args[2]
		if outputFile == "" {
			log.Fatal("Please provide an output filename")
		}
	}

	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var allBytes []byte

	scanner := bufio.NewScanner(file)
	lineNo := 1
	for scanner.Scan() {
		line := scanner.Text()

		// Lex and parse
		toks := lexer.LexLine(line, lineNo)
		inst, err := parser.ParseInstruction(toks)
		if err != nil {
			log.Fatalf("Error parsing line #%d: %s, error: %v", lineNo, line, err)
		}

		allBytes = append(allBytes, inst.ToMachineCode()...)
		lineNo++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Parsed %d lines", lineNo)
	log.Printf("Writing %s", outputFile)
	err = os.WriteFile(outputFile, allBytes, 0644)
	if err != nil {
		log.Fatalf("Failed to write %s: %v", outputFile, err)
	}
}
