/*
 * --- Part Two ---
 *
 * It turns out that this circuit is very timing-sensitive; you actually need to minimize the signal delay.
 *
 * To do this, calculate the number of steps each wire takes to reach each intersection;
 * choose the intersection where the sum of both wires' steps is lowest. If a wire visits a position on the grid multiple times,
 * use the steps value from the first time it visits that position when calculating the total value of a specific intersection.
 *
 * The number of steps a wire takes is the total number of grid squares the wire has entered to get to that location,
 * including the intersection being considered. Again consider the example from above:
 *
 * ...........
 * .+-----+...
 * .|.....|...
 * .|..+--X-+.
 * .|..|..|.|.
 * .|.-X--+.|.
 * .|..|....|.
 * .|.......|.
 * .o-------+.
 * ...........
 *
 * In the above example, the intersection closest to the central port is reached after 8+5+5+2 = 20 steps by the first wire
 * and 7+6+4+3 = 20 steps by the second wire for a total of 20+20 = 40 steps.
 *
 * However, the top-right intersection is better: the first wire takes only 8+5+2 = 15 and the second wire takes only
 * 7+6+2 = 15, a total of 15+15 = 30 steps.
 *
 * Here are the best steps for the extra examples from above:
 *
 *  - R75,D30,R83,U83,L12,D49,R71,U7,L72
 *    U62,R66,U55,R34,D71,R55,D58,R83 = 610 steps
 *
 *  - R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51
 *    U98,R91,D20,R16,D67,R40,U7,R15,U6,R7 = 410 steps
 *
 * What is the fewest combined steps the wires must take to reach an intersection?
 */

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Position struct {
	X int
	Y int
}

func (positionA *Position) Distance(positionB Position) int {
	return int(math.Abs(float64(positionA.X-positionB.X)) + math.Abs(float64(positionA.Y-positionB.Y)))
}
func (position *Position) IsOrigin() bool {
	return position.X == 0 && position.Y == 0
}

type Path struct {
	StartPoint Position
	Direction  byte
	Distance   int
}

func (path *Path) IsHorizontal() bool {
	return path.Direction == 'R' || path.Direction == 'L'
}
func (path *Path) IsVertical() bool {
	return path.Direction == 'U' || path.Direction == 'D'
}
func (pathA *Path) MovementIsOnSameAxisAs(pathB Path) bool {
	return (pathA.IsHorizontal() && pathB.IsHorizontal()) || (pathA.IsVertical() && pathB.IsVertical())
}
func (path *Path) EndPoint() Position {
	switch path.Direction {
	case 'R':
		return Position{path.StartPoint.X + path.Distance, path.StartPoint.Y}
	case 'L':
		return Position{path.StartPoint.X - path.Distance, path.StartPoint.Y}
	case 'U':
		return Position{path.StartPoint.X, path.StartPoint.Y - path.Distance}
	case 'D':
		return Position{path.StartPoint.X, path.StartPoint.Y + path.Distance}
	}

	log.Fatalf("Unknown direction %v", path.Direction)
	// This return will never happen
	return Position{}
}
func (path *Path) Contains(point Position) bool {
	start := path.StartPoint
	end := path.EndPoint()

	xMin := int(math.Min(float64(start.X), float64(end.X)))
	xMax := int(math.Max(float64(start.X), float64(end.X)))
	yMin := int(math.Min(float64(start.Y), float64(end.Y)))
	yMax := int(math.Max(float64(start.Y), float64(end.Y)))

	return point.X >= xMin && point.X <= xMax && point.Y >= yMin && point.Y <= yMax
}
func (path *Path) HopsTo(point Position) int {
	if !path.Contains(point) {
		return int(math.Inf(1))
	}

	return path.StartPoint.Distance(point)
}

func (pathA *Path) Crosses(pathB Path) (bool, Position) {
	// if they go in the same direction we assume they don't cross
	// (note: If they are at exactly the same position this assumption will not hold true)
	if pathA.MovementIsOnSameAxisAs(pathB) {
		return false, Position{}
	}

	var horizontalPath, verticalPath Path
	if pathA.IsHorizontal() {
		horizontalPath = *pathA
		verticalPath = pathB
	} else {
		horizontalPath = pathB
		verticalPath = *pathA
	}

	potentialCrossingPoint := Position{verticalPath.StartPoint.X, horizontalPath.StartPoint.Y}
	if pathA.Contains(potentialCrossingPoint) && pathB.Contains(potentialCrossingPoint) {
		return true, potentialCrossingPoint
	}

	return false, Position{}
}

type Wire struct {
	Path []Path
}

func (wireA *Wire) CrossingPoints(wireB Wire) []Position {
	crossingPoints := make([]Position, 0)

	for _, pathA := range wireA.Path {
		for _, pathB := range wireB.Path {
			crosses, point := pathA.Crosses(pathB)
			if crosses && !point.IsOrigin() {
				crossingPoints = append(crossingPoints, point)
			}
		}
	}

	return crossingPoints
}

func (wire *Wire) HopsTo(point Position) int {

	hops := 0
	for _, path := range wire.Path {
		if path.Contains(point) {
			hops += path.HopsTo(point)
			return hops
		}
		hops += path.Distance
	}

	// If we have arrived at this point the point does not seem to be on the wire
	// at all. This should not happen, but if it does the hops are returned as infinity
	return int(math.Inf(1))
}

func NewWire(commaSeparatedInstructions string) Wire {
	currentPosition := Position{0, 0}
	split := strings.Split(commaSeparatedInstructions, ",")
	paths := make([]Path, len(split))

	for i, instruction := range split {
		runeInstruction := []rune(instruction)
		direction := byte(runeInstruction[0])
		distance, err := strconv.Atoi(string(runeInstruction[1:]))
		if err != nil {
			log.Fatal(err)
		}
		paths[i] = Path{currentPosition, direction, distance}
		currentPosition = paths[i].EndPoint()
	}

	return Wire{paths}
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(input), "\n")
	wireA := NewWire(lines[0])
	wireB := NewWire(lines[1])

	crossingPoints := wireA.CrossingPoints(wireB)
	closestCrossingPoint := crossingPoints[0]
	closestDistance := wireA.HopsTo(closestCrossingPoint) + wireB.HopsTo(closestCrossingPoint)
	for _, point := range crossingPoints {
		distance := wireA.HopsTo(point) + wireB.HopsTo(point)
		if distance < closestDistance {
			closestDistance = distance
			closestCrossingPoint = point
		}
	}

	fmt.Print(closestDistance)
}
