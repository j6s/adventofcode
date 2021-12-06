package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	lastDepth := -1
	depthIncreases := 0
	for scanner.Scan() {
		depth, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		if lastDepth != -1 && depth > lastDepth {
			depthIncreases++
		}

		lastDepth = depth
	}

	log.Printf("%d increases", depthIncreases)
}
