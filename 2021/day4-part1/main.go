package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type BingoBoard struct {
	numbers [][]int
	flagged [][]bool
	rows    int
	cols    int
}

func NewBingoBoard(numbers [][]int) BingoBoard {
	rows := len(numbers)
	flagged := make([][]bool, rows)

	cols := 0
	if rows > 0 {
		cols = len(numbers[0])
	}

	for i, _ := range flagged {
		flagged[i] = make([]bool, cols)
	}

	return BingoBoard{
		numbers: numbers,
		flagged: flagged,
		rows:    rows,
		cols:    cols,
	}
}

func (this *BingoBoard) Validate() error {
	for i, row := range this.numbers {
		if len(row) != this.cols {
			return fmt.Errorf("Row %d does not have %d columns", i, this.cols)
		}
	}
	return nil
}

func (this *BingoBoard) Flag(number int) bool {
	hasWon := false

	for rowIndex, row := range this.numbers {
		for colIndex, boardNumber := range row {
			if boardNumber == number {
				this.flagged[rowIndex][colIndex] = true
				hasWon = hasWon || this.rowHasWon(rowIndex) || this.colHasWon(colIndex)
			}
		}
	}

	return hasWon
}

func (this *BingoBoard) rowHasWon(rowIndex int) bool {
	for _, flagged := range this.flagged[rowIndex] {
		if !flagged {
			return false
		}
	}
	return true
}

func (this *BingoBoard) colHasWon(colIndex int) bool {
	for _, row := range this.flagged {
		if !row[colIndex] {
			return false
		}
	}
	return true
}

func (this *BingoBoard) UnflaggedNumbers() []int {
	unflagged := make([]int, 0)
	for rowIndex, row := range this.flagged {
		for colIndex, flagged := range row {
			if !flagged {
				unflagged = append(unflagged, this.numbers[rowIndex][colIndex])
			}
		}
	}
	return unflagged
}

type BingoGame struct {
	inputs []int
	boards []BingoBoard
}

func NewBingoGame(inputs []int, boards []BingoBoard) BingoGame {
	return BingoGame{
		inputs: inputs,
		boards: boards,
	}
}

func (this *BingoGame) Validate() error {

	rows := -1
	cols := -1
	for i, board := range this.boards {
		if rows == -1 {
			rows = board.rows
		}
		if cols == -1 {
			cols = board.cols
		}
		err := board.Validate()
		if err != nil {
			return fmt.Errorf("Board %d is invalid: %s", i, err)
		}
		if rows != board.rows {
			return fmt.Errorf("Board %d does not have the same number of rows as other boards", i)
		}
		if cols != board.cols {
			return fmt.Errorf("Board %d does not have the same number of cols as other boards", i)
		}
	}

	return nil
}

func (this *BingoGame) playRound(input int) (bool, BingoBoard) {
	for _, board := range this.boards {
		if board.Flag(input) {
			return true, board
		}
	}

	return false, BingoBoard{}
}

func (this *BingoGame) Play() (int, BingoBoard) {
	for _, input := range this.inputs {
		hasWinner, winner := this.playRound(input)
		if hasWinner {
			return input, winner
		}
	}

	log.Fatal("No winner exists")
	return -1, BingoBoard{}
}

func main() {
	game := readBingoGameFromStdin()
	err := game.Validate()
	if err != nil {
		log.Fatal(err)
	}

	winningInput, winningBoard := game.Play()

	result := sum(winningBoard.UnflaggedNumbers()) * winningInput
	log.Printf("First winning board score: %d", result)
}

func readBingoGameFromStdin() BingoGame {
	scanner := bufio.NewScanner(os.Stdin)

	// First line contains inputs
	scanner.Scan()
	inputs := lineToNumbers(scanner.Text(), regexp.MustCompile(","))

	// Second line is spacer
	scanner.Scan()

	currentBoardNumbers := make([][]int, 0)
	separator := regexp.MustCompile(`\s+`)
	boards := make([]BingoBoard, 0)
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " ")
		if line == "" {
			boards = append(boards, NewBingoBoard(currentBoardNumbers))
			currentBoardNumbers = make([][]int, 0)
			continue
		}

		currentBoardNumbers = append(currentBoardNumbers, lineToNumbers(line, separator))
	}

	boards = append(boards, NewBingoBoard(currentBoardNumbers))
	return NewBingoGame(inputs, boards)
}

func lineToNumbers(line string, separator *regexp.Regexp) []int {
	parts := separator.Split(line, -1)
	numbers := make([]int, len(parts))

	for i, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			log.Fatal(err)
		}
		numbers[i] = num
	}

	return numbers
}

func sum(numbers []int) int {
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	return sum
}
