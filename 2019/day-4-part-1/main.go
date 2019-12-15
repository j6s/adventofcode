/*
 * --- Day 4: Secure Container ---
 *
 * You arrive at the Venus fuel depot only to discover it's protected by a password.
 * The Elves had written the password on a sticky note, but someone threw it out.
 *
 * However, they do remember a few key facts about the password:
 *
 *  - It is a six-digit number.
 *  - The value is within the range given in your puzzle input.
 *  - Two adjacent digits are the same (like 22 in 122345).
 *  - Going from left to right, the digits never decrease; they only ever
 *    increase or stay the same (like 111123 or 135679).
 *
 * Other than the range rule, the following are true:
 *
 *  - 111111 meets these criteria (double 11, never decreases).
 *  - 223450 does not meet these criteria (decreasing pair of digits 50).
 *  - 123789 does not meet these criteria (no double).
 *
 * How many different passwords within the range given in your puzzle input meet these criteria?
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

func (code *Passcode) HasTwoAdjacentDigits() bool {
	var lastDigit int
	for i, digit := range code.Code {
		if i != 0 && lastDigit == digit {
			return true
		}
		lastDigit = digit
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
	return code.HasTwoAdjacentDigits() && code.NeverDecreases()
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
