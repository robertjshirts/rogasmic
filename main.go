// 32 bit arm assembler
package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	var instructions []Instruction
	file, err := os.Open("asm.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	col := 0
	for scanner.Scan() {
		line := scanner.Text()
		instruction, err := parseLine(line, col)
		if err != nil {
			log.Printf("Error parsing line: %s, error: %v", line, err)
			os.Exit(1)
		}
		instructions = append(instructions, instruction)
		col++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func parseLine(line string, col int) (Instruction, error) {
	toks := LexLine(line)
	return nil, nil
}
