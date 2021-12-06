package main

import (
	"bufio"
	"log"
	"math"
	"os"
)

func main() {
	input := readInput()
	gammaBits, epsilonBits := extractMostAndLeastCommonBits(input)
	gammaRate := bitsToInt(gammaBits)
	epsilonRate := bitsToInt(epsilonBits)

	log.Printf("gammaRate: %d", gammaRate)
	log.Printf("epsilonRate: %d", epsilonRate)
	log.Printf("product: %d", gammaRate*epsilonRate)
}

func readInput() [][]bool {
	scanner := bufio.NewScanner(os.Stdin)

	input := make([][]bool, 0)
	for scanner.Scan() {
		text := scanner.Text()

		currentLine := make([]bool, len(text))
		for i, bit := range text {
			if bit == '1' {
				currentLine[i] = true
			} else if bit == '0' {
				currentLine[i] = false
			} else {
				log.Fatalf("%s is not a valid bit value (must be 1 or 0)")
			}
		}

		input = append(input, currentLine)
	}

	return input
}

func extractMostAndLeastCommonBits(input [][]bool) ([]bool, []bool) {
	trueCount := make([]int, len(input[0]))
	falseCount := make([]int, len(input[0]))

	for _, line := range input {
		for i, bit := range line {
			if bit {
				trueCount[i]++
			} else {
				falseCount[i]++
			}
		}
	}

	mostCommon := make([]bool, len(input[0]))
	leastCommon := make([]bool, len(input[0]))
	for i, _ := range trueCount {
		mostCommon[i] = trueCount[i] > falseCount[i]
		leastCommon[i] = !mostCommon[i]
	}

	return mostCommon, leastCommon
}

func bitsToInt(bits []bool) int {
	result := 0

	for i, bit := range bits {
		if bit {
			result += int(math.Pow(2, float64(len(bits)-i-1)))
		}
	}

	return result
}
