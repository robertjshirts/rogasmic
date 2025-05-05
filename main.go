// 32 bit arm assembler
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/robertjshirts/rogasmic/lexer"
	"github.com/robertjshirts/rogasmic/parser"
)

func main() {
	//var instructions []types.Instruction
	file, err := os.Open("asm.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var allBytes []byte

	scanner := bufio.NewScanner(file)
	lineNo := 0
	for scanner.Scan() {
		line := scanner.Text()
		bytes, err := parseLine(line, lineNo)
		if err != nil {
			log.Printf("Error parsing line: %s, error: %v", line, err)
			os.Exit(1)
		}
		allBytes = append(allBytes, bytes...)
		lineNo++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("kernel7.img", allBytes, 0644)
	if err != nil {
		log.Fatalf("Failed to write kernel7.img: %v", err)
	}
}

func parseLine(line string, lineNo int) ([]byte, error) {
	toks := lexer.LexLine(line, lineNo)
	bytes, err := parser.ParseInstruction(toks)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%d: % X\n", lineNo+1, bytes)

	return bytes, nil
}
