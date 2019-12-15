/*
 *  --- Day 3: Crossed Wires ---
 *
 * The gravity assist was successful, and you're well on your way to the Venus refuelling station.
 * During the rush back on Earth, the fuel management system wasn't completely installed, so that's
 * next on the priority list.
 *
 * Opening the front panel reveals a jumble of wires. Specifically, two wires are connected to a central
 * port and extend outward on a grid. You trace the path each wire takes as it leaves the central port,
 * one wire per line of text (your puzzle input).
 *
 * The wires twist and turn, but the two wires occasionally cross paths. To fix the circuit, you need to
 * find the intersection point closest to the central port. Because the wires are on a grid, use the
 * Manhattan distance for this measurement. While the wires do technically cross right at the central port
 * where they both start, this point does not count, nor does a wire count as crossing with itself.
 *
 * For example, if the first wire's path is R8,U5,L5,D3, then starting from the central port (o),
 * it goes right 8, up 5, left 5, and finally down 3:
 *
 * ...........
 * ...........
 * ...........
 * ....+----+.
 * ....|....|.
 * ....|....|.
 * ....|....|.
 * .........|.
 * .o-------+.
 * ...........
 *
 * Then, if the second wire's path is U7,R6,D4,L4, it goes up 7, right 6, down 4, and left 4:
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
 * These wires cross at two locations (marked X), but the lower-left one is closer to the central port: its distance is 3 + 3 = 6.
 *
 * Here are a few more examples:
 *
 *     R75,D30,R83,U83,L12,D49,R71,U7,L72
 *     U62,R66,U55,R34,D71,R55,D58,R83 = distance 159
 *     R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51
 *     U98,R91,D20,R16,D67,R40,U7,R15,U6,R7 = distance 135
 *
 * What is the Manhattan distance from the central port to the closest intersection?
 *
 */

/*
 * Not about my approach:
 * My first version of this calculated every single point on the wires grid and compared every point with
 * every other point. This was very inefficient which then lead me to compare paths as whole instead of indivitual points.
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
	closestDistance := closestCrossingPoint.Distance(Position{0, 0})
	for _, point := range crossingPoints {
		distance := point.Distance(Position{0, 0})
		if distance < closestDistance {
			closestDistance = distance
			closestCrossingPoint = point
		}
	}

	fmt.Print(closestDistance)
}
