package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	SPAWN_DAYS             = 7
	FIRST_GENERATION_DELAY = 2
	SIMULATE_DAYS          = 256
)

type FishGeneration struct {
	daysUntilGrowth int
	size            int
}

func (this *FishGeneration) Process1Day() bool {
	if this.daysUntilGrowth == 0 {
		return true
	}

	this.daysUntilGrowth--

	return false
}

type FishCollection struct {
	fish []*FishGeneration
}

func (this *FishCollection) String() string {
	elements := make([]string, len(this.fish))
	for i, fish := range this.fish {
		elements[i] = fmt.Sprintf("%d", fish.daysUntilGrowth)
	}
	return strings.Join(elements, ",")
}

func (this *FishCollection) Process1Day() {
	newFish := 0

	for _, generation := range this.fish {
		if generation.Process1Day() {
			generation.daysUntilGrowth = SPAWN_DAYS - 1
			newFish += generation.size
		}
	}

	if newFish > 0 {
		this.fish = append(this.fish, &FishGeneration{SPAWN_DAYS - 1 + FIRST_GENERATION_DELAY, newFish})
	}
}

func (this *FishCollection) Count() int {
	count := 0
	for _, generation := range this.fish {
		count += generation.size
	}
	return count
}

func main() {
	fish := readFishFromStdin()

	for i := 0; i < SIMULATE_DAYS; i++ {
		fish.Process1Day()
		// log.Printf("After %d days: %s", i, fish.String())
	}

	log.Printf("%d fish after %d days", fish.Count(), SIMULATE_DAYS)
}

func readFishFromStdin() FishCollection {
	fish := make([]*FishGeneration, 0)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		numStrings := strings.Split(scanner.Text(), ",")
		for _, numString := range numStrings {
			num, err := strconv.Atoi(numString)
			if err != nil {
				log.Fatal(err)
			}
			fish = append(fish, &FishGeneration{num, 1})
		}
	}

	return FishCollection{fish}
}
