// 32 bit arm assembler
package main

import (
	"bufio"
	"log"
	"os"

	"github.com/robertjshirts/rogasmic/lexer"
	"github.com/robertjshirts/rogasmic/parser"
	"github.com/robertjshirts/rogasmic/types"
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
		inst, err := parseLine(line, lineNo)
		if err != nil {
			log.Fatalf("Error parsing line: %s, error: %v", line, err)
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

func parseLine(line string, lineNo int) (types.Instruction, error) {
	toks := lexer.LexLine(line, lineNo)
	bytes, err := parser.ParseInstruction(toks)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
