package main

import (
	"bufio"
	"log"
	"math"
	"os"
)

func main() {
	input := readInput()

	oxygen := bitsToInt(progressivelyFilterInput(input, func(data [][]bool, position int) [][]bool {
		oxygenCriteria, _ := extractBitCriteriaAtPosition(data, position)
		return filterByBitCriteria(data, position, oxygenCriteria)
	}))
	co2 := bitsToInt(progressivelyFilterInput(input, func(data [][]bool, position int) [][]bool {
		_, co2Criteria := extractBitCriteriaAtPosition(data, position)
		return filterByBitCriteria(data, position, co2Criteria)
	}))

	log.Printf("oxygen: %d", oxygen)
	log.Printf("co2: %d", co2)
	log.Printf("product: %d", oxygen*co2)
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

func progressivelyFilterInput(input [][]bool, filter func(input [][]bool, position int) [][]bool) []bool {
	inputsToConsider := input
	i := 0
	for i < len(input[0]) {
		inputsToConsider = filter(inputsToConsider, i)
		if len(inputsToConsider) == 1 {
			return inputsToConsider[0]
		}
		i++
	}

	log.Fatal("Could not extract reading")
	return []bool{}
}

func filterByBitCriteria(input [][]bool, position int, criteria bool) [][]bool {
	filtered := make([][]bool, 0)
	for _, line := range input {
		if line[position] == criteria {
			filtered = append(filtered, line)
		}
	}
	return filtered
}

func extractBitCriteriaAtPosition(input [][]bool, position int) (bool, bool) {
	trueCount := 0
	falseCount := 0

	for _, line := range input {
		if line[position] {
			trueCount++
		} else {
			falseCount++
		}
	}

	oxygenCriteria := trueCount >= falseCount
	co2Criteria := trueCount < falseCount

	return oxygenCriteria, co2Criteria
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
