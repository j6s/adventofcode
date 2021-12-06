package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Point struct {
	x int
	y int
}

func (this *Point) String() string {
	return fmt.Sprintf("%d,%d", this.x, this.y)
}

type AffectedPoint struct {
	Point
	crossings int
}

type Line struct {
	begin Point
	end   Point
}

func NewLine(x1 int, y1 int, x2 int, y2 int) Line {
	return Line{
		begin: Point{x: x1, y: y1},
		end:   Point{x: x2, y: y2},
	}
}

func (this *Line) String() string {
	return fmt.Sprintf("%s -> %s", this.begin.String(), this.end.String())
}

func (this *Line) IsHorizontal() bool {
	return this.begin.y == this.end.y
}
func (this *Line) IsVertical() bool {
	return this.begin.x == this.end.x
}
func (this *Line) GetCrossingPoints() []Point {

	// TODO Handle non-horizontal and vertical lines
	if this.IsHorizontal() {
		xs := []int{this.begin.x, this.end.x}
		sort.Ints(xs)

		points := make([]Point, 0)
		for x := xs[0]; x <= xs[1]; x++ {
			points = append(points, Point{x: x, y: this.begin.y})
		}
		return points
	}

	if this.IsVertical() {
		ys := []int{this.begin.y, this.end.y}
		sort.Ints(ys)

		points := make([]Point, 0)
		for y := ys[0]; y <= ys[1]; y++ {
			points = append(points, Point{x: this.begin.x, y: y})
		}
		return points
	}

	return []Point{}
}

type Grid struct {
	lines []Line
}

func NewGrid(lines []Line) Grid {
	return Grid{lines: lines}
}

func (this *Grid) GetAffectedPoints() []AffectedPoint {
	crossings := make(map[int]map[int]int)

	for _, line := range this.lines {
		for _, point := range line.GetCrossingPoints() {
			if _, ok := crossings[point.x]; !ok {
				crossings[point.x] = make(map[int]int)
			}
			if _, ok := crossings[point.x][point.y]; !ok {
				crossings[point.x][point.y] = 0
			}
			crossings[point.x][point.y]++
		}
	}

	points := make([]AffectedPoint, 0)
	for x, ys := range crossings {
		for y, num := range ys {
			points = append(points, AffectedPoint{Point{x, y}, num})
		}
	}

	return points
}

func (this *Grid) String() string {
	size := Point{0, 0}
	for _, line := range this.lines {
		if line.begin.x > size.x {
			size.x = line.begin.x
		}
		if line.end.x > size.x {
			size.x = line.end.x
		}
		if line.begin.y > size.y {
			size.y = line.begin.y
		}
		if line.end.y > size.y {
			size.y = line.end.y
		}
	}

	lines := make([]string, size.y+1)
	for y := 0; y < size.y+1; y++ {
		lines[y] = strings.Repeat(".", size.x+1)
	}

	for _, point := range this.GetAffectedPoints() {
		lines[point.y] = replaceAtIndex(lines[point.y], fmt.Sprintf("%d", point.crossings)[0], point.x)
	}

	return strings.Join(lines, "\n")
}

func main() {
	grid := readGridFromStdin()

	numberOfDangerousPoints := 0
	for _, point := range grid.GetAffectedPoints() {
		if point.crossings >= 2 {
			numberOfDangerousPoints++
		}
	}

	// log.Printf("\n%s", grid.String())
	log.Printf("Number of dangerous crossings: %d", numberOfDangerousPoints)
}

func readGridFromStdin() Grid {
	lines := make([]Line, 0)
	scanner := bufio.NewScanner(os.Stdin)
	separator := regexp.MustCompile(`(,| -> )`)
	for scanner.Scan() {
		points := lineToNumbers(strings.Trim(scanner.Text(), " "), separator)
		line := NewLine(points[0], points[1], points[2], points[3])

		if line.IsHorizontal() || line.IsVertical() {
			lines = append(lines, line)
		}
	}

	return NewGrid(lines)
}

func lineToNumbers(line string, separator *regexp.Regexp) []int {
	parts := separator.Split(line, -1)
	numbers := make([]int, len(parts))

	for i, part := range parts {
		num, err := strconv.Atoi(strings.Trim(part, " "))
		if err != nil {
			log.Fatal(err)
		}
		numbers[i] = num
	}

	return numbers
}

func replaceAtIndex(in string, r byte, i int) string {
	out := []byte(in)
	out[i] = r
	return string(out)
}
