/*
 * --- Part Two ---
 *
 * An Elf just remembered one more important detail:
 * the two adjacent matching digits are not part of a larger group
 * of matching digits.
 *
 * Given this additional criterion, but still ignoring the range rule,
 * the following are now true:
 *
 *  - 112233 meets these criteria because the digits never decrease and all
 *    repeated digits are exactly two digits long.
 *  - 123444 no longer meets the criteria (the repeated 44 is part of a larger
 *    group of 444).
 *  - 111122 meets the criteria (even though 1 is repeated more than twice,
 *    it still contains a double 22).
 *
 * How many different passwords within the range given in your puzzle input meet all of the criteria?
 *
 */

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func getValidRangeFromStdin() (min int, max int) {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	validRange := strings.Split(string(input), "-")

	min, err = strconv.Atoi(validRange[0])
	if err != nil {
		log.Fatal(err)
	}

	max, err = strconv.Atoi(validRange[1])
	if err != nil {
		log.Fatal(err)
	}

	return
}

type Passcode struct {
	Code []int
}

func NewPasscode(inputCode int) Passcode {
	// Going int > string > []rune > []int seems a bit inefficient
	// TODO Extracts digits from input without going to string first
	codeString := strconv.Itoa(inputCode)
	digits := []rune(codeString)
	code := make([]int, len(digits))

	for i, digit := range digits {
		num, err := strconv.Atoi(string(digit))
		if err != nil {
			log.Fatal(err)
		}
		code[i] = num
	}

	return Passcode{code}
}

func (code *Passcode) HasExactlyTwoAdjacentDigits() bool {
	repetitions := make([]int, 0)
	var lastDigit int

	for i, digit := range code.Code {
		if i != 0 && lastDigit == digit {
			repetitions[len(repetitions)-1]++
		} else {
			repetitions = append(repetitions, 1)
		}
		lastDigit = digit
	}

	for _, rep := range repetitions {
		if rep == 2 {
			return true
		}
	}
	return false
}

func (code *Passcode) NeverDecreases() bool {
	var lastDigit int
	for i, digit := range code.Code {
		if i != 0 && lastDigit > digit {
			return false
		}
		lastDigit = digit
	}

	return true
}

func (code *Passcode) IsValid() bool {
	return code.HasExactlyTwoAdjacentDigits() && code.NeverDecreases()
}

func main() {
	min, max := getValidRangeFromStdin()

	valid := 0
	for input := min; input <= max; input++ {
		code := NewPasscode(input)
		if code.IsValid() {
			valid++
		}
	}

	fmt.Print(valid)
}
