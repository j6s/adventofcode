/*
* --- Part Two ---
*
* "Good, the new computer seems to be working correctly! Keep it nearby during this mission - you'll probably use it again.
* Real Intcode computers support many more features than your new one, but we'll let you know what they are as you need them."
*
* "However, your current priority should be to complete your gravity assist around the Moon.
* For this mission to succeed, we should settle on some terminology for the parts you've already built."
*
* Intcode programs are given as a list of integers; these values are used as the initial state for the computer's memory.
* When you run an Intcode program, make sure to start by initializing memory to the program's values.
* A position in memory is called an address (for example, the first value in memory is at "address 0").
*
* Opcodes (like 1, 2, or 99) mark the beginning of an instruction.
* The values used immediately after an opcode, if any, are called the instruction's parameters.
* For example, in the instruction 1,2,3,4, 1 is the opcode; 2, 3, and 4 are the parameters.
* The instruction 99 contains only an opcode and has no parameters.
*
* The address of the current instruction is called the instruction pointer; it starts at 0.
* After an instruction finishes, the instruction pointer increases by the number of values in the instruction;
* until you add more instructions to the computer, this is always 4 (1 opcode + 3 parameters) for the add and multiply instructions.
* (The halt instruction would increase the instruction pointer by 1, but it halts the program instead.)
*
* "With terminology out of the way, we're ready to proceed. To complete the gravity assist,
* you need to determine what pair of inputs produces the output 19690720."
*
* The inputs should still be provided to the program by replacing the values at addresses 1 and 2, just like before.
* In this program, the value placed in address 1 is called the noun, and the value placed in address 2 is called the verb.
* Each of the two input values will be between 0 and 99, inclusive.
*
* Once the program has halted, its output is available at address 0, also just like before.
* Each time you try a pair of inputs, make sure you first reset the computer's memory to the values in the program (your puzzle input)
* - in other words, don't reuse memory from a previous attempt.
*
* Find the input noun and verb that cause the program to produce the output 19690720.
* What is 100 * noun + verb? (For example, if noun=12 and verb=2, the answer would be 1202.)
*
 */
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type IntCode struct {
	code               []int
	instructionPointer int
}

func NewIntCode(commaSeparated string) IntCode {
	split := strings.Split(commaSeparated, ",")
	code := make([]int, len(split))

	for i, str := range split {
		num, err := strconv.Atoi(str)
		if err != nil {
			log.Fatal(err)
		}
		code[i] = num
	}

	return IntCode{code, 0}
}

func (intCode *IntCode) Get(address int) int {
	return intCode.code[address]
}

func (intCode *IntCode) Set(address int, value int) {
	intCode.code[address] = value
}

func (intCode *IntCode) Clone() IntCode {
	code := make([]int, len(intCode.code))
	copy(code, intCode.code)
	return IntCode{code, intCode.instructionPointer}
}

func (intCode *IntCode) GetSlice(start int, length int) []int {
	return intCode.code[start:(start + length)]
}

func (intCode *IntCode) String() string {
	str := make([]string, len(intCode.code))
	for i, code := range intCode.code {
		str[i] = strconv.Itoa(code)
	}
	return strings.Join(str, ",")
}

func (intCode *IntCode) RunStep() (isDone bool, err error) {
	instruction := intCode.Get(intCode.instructionPointer)
	if instruction == 99 {
		isDone = true
	} else {
		isDone = false
		param := intCode.GetSlice(intCode.instructionPointer+1, 3)
		switch instruction {
		case 1:
			intCode.Set(param[2], intCode.Get(param[0])+intCode.Get(param[1]))
			break
		case 2:
			intCode.Set(param[2], intCode.Get(param[0])*intCode.Get(param[1]))
			break
		default:
			err = errors.New(fmt.Sprintf("Invalid intcode instruction %v encountered", instruction))
		}

	}

	intCode.instructionPointer += 4
	return
}

func (intCode *IntCode) Run() error {
	intCode.instructionPointer = 0

	for true {
		isDone, err := intCode.RunStep()
		if err != nil {
			return err
		}
		if isDone {
			return nil
		}
	}

	return errors.New("No exit instruction (99) encountered")
}

func main() {
	lines, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	desired := 19690720
	intCode := NewIntCode(string(lines))

	for noun := 0; noun <= 99; noun++ {
		for verb := 0; verb <= 99; verb++ {
			tempCode := intCode.Clone()
			tempCode.Set(1, noun)
			tempCode.Set(2, verb)
			err := tempCode.Run()
			if err != nil {
				continue
			}
			result := tempCode.Get(0)
			if result == desired {
				fmt.Print(100*noun + verb)
				os.Exit(0)
			}
		}
	}

	fmt.Printf("Could not find noun & verb producing result %d\n", desired)
	os.Exit(1)
}
