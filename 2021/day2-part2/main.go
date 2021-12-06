package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	INSTRUCTION_FORWARD = 'f'
	INSTRUCTION_DOWN    = 'd'
	INSTRUCTION_UP      = 'u'
)

type Instruction struct {
	Type   byte
	Amount int
}

func main() {
	instructions := readInstructions()

	aim := 0
	depth := 0
	horizontalPosition := 0

	for _, instruction := range instructions {
		switch instruction.Type {
		case INSTRUCTION_FORWARD:
			horizontalPosition += instruction.Amount
			depth += aim * instruction.Amount
			break
		case INSTRUCTION_DOWN:
			aim += instruction.Amount
			break
		case INSTRUCTION_UP:
			aim -= instruction.Amount
			break
		}
	}

	log.Printf("depth: %d", depth)
	log.Printf("horizontalPosition: %d", horizontalPosition)
	log.Printf("product: %d", depth*horizontalPosition)
}

func readInstructionType(instructionType string) byte {
	switch instructionType {
	case "forward":
		return INSTRUCTION_FORWARD
	case "up":
		return INSTRUCTION_UP
	case "down":
		return INSTRUCTION_DOWN
	}

	log.Fatalf("Unkonwn instruction %s", instructionType)
	return '-'
}

func readInstructions() []Instruction {
	scanner := bufio.NewScanner(os.Stdin)

	instructions := make([]Instruction, 0)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")

		instructionType := readInstructionType(parts[0])
		amount, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Fatal(err)
		}

		instructions = append(instructions, Instruction{instructionType, amount})
	}

	return instructions
}
