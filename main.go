// 32 bit arm assembler
package main

import (
	"bufio"
	"log"
	"os"

	"github.com/robertjshirts/rogasmic/lexer"
	"github.com/robertjshirts/rogasmic/types"
)

func main() {
	types.Init()

	//var instructions []types.Instruction
	file, err := os.Open("asm.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNo := 0
	for scanner.Scan() {
		line := scanner.Text()
		_, err := parseLine(line, lineNo)
		if err != nil {
			log.Printf("Error parsing line: %s, error: %v", line, err)
			os.Exit(1)
		}
		lineNo++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func parseLine(line string, lineNo int) (types.Instruction, error) {
	toks := lexer.LexLine(line, lineNo)
	instruction := instructions.Parse(toks)

	return nil, nil
}
