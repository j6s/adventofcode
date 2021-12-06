package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

func main() {
	depths := readDepths()
	slidingWindows := computeSlidingWindows(depths, 3)
	sums := sumSlidingWindows(slidingWindows)

	increases := 0
	lastSum := -1
	for _, sum := range sums {
		if lastSum != -1 && sum > lastSum {
			increases++
		}
		lastSum = sum
	}

	log.Printf("%d increases", increases)
}

func readDepths() []int {
	scanner := bufio.NewScanner(os.Stdin)
	depths := make([]int, 0)

	for scanner.Scan() {
		depth, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		depths = append(depths, depth)
	}

	return depths
}

func computeSlidingWindows(depths []int, size int) [][]int {
	slidingWindows := make([][]int, 0)

	i := size - 1
	for i < len(depths) {
		slidingWindows = append(
			slidingWindows,
			[]int{depths[i-2], depths[i-1], depths[i]},
		)
		i++
	}
	return slidingWindows
}

func sumSlidingWindows(slidingWindows [][]int) []int {
	sums := make([]int, 0)

	for _, window := range slidingWindows {
		sum := 0
		for _, depth := range window {
			sum += depth
		}
		sums = append(sums, sum)
	}

	return sums
}
