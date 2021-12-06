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
	SIMULATE_DAYS          = 80
)

type Fish struct {
	daysUntilGrowth int
}

func (this *Fish) Process1Day() bool {
	if this.daysUntilGrowth == 0 {
		return true
	}

	this.daysUntilGrowth--

	return false
}

type FishCollection struct {
	fish []*Fish
}

func (this *FishCollection) String() string {
	elements := make([]string, len(this.fish))
	for i, fish := range this.fish {
		elements[i] = fmt.Sprintf("%d", fish.daysUntilGrowth)
	}
	return strings.Join(elements, ",")
}

func (this *FishCollection) Process1Day() {
	newFish := make([]*Fish, 0)

	for _, fish := range this.fish {
		if fish.Process1Day() {
			fish.daysUntilGrowth = SPAWN_DAYS - 1
			newFish = append(newFish, &Fish{SPAWN_DAYS - 1 + FIRST_GENERATION_DELAY})
		}
	}

	this.fish = append(this.fish, newFish...)
}

func (this *FishCollection) Count() int {
	return len(this.fish)
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
	fish := make([]*Fish, 0)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		numStrings := strings.Split(scanner.Text(), ",")
		for _, numString := range numStrings {
			num, err := strconv.Atoi(numString)
			if err != nil {
				log.Fatal(err)
			}
			fish = append(fish, &Fish{num})
		}
	}

	return FishCollection{fish}
}
