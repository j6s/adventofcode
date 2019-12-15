/*
 * --- Part Two ---
 *
 * The air conditioner comes online! Its cold air feels good for a while, but then the TEST alarms start to go off.
 * Since the air conditioner can't vent its heat anywhere but back into the spacecraft, it's actually making the air
 * inside the ship warmer.
 *
 * Instead, you'll need to use the TEST to extend the thermal radiators. Fortunately, the diagnostic program
 * (your puzzle input) is already equipped for this. Unfortunately, your Intcode computer is not.
 *
 * Your computer is only missing a few opcodes:
 *
 *  - Opcode 5 is jump-if-true: if the first parameter is non-zero, it sets the instruction pointer to the value from the
 *    second parameter. Otherwise, it does nothing.
 *  - Opcode 6 is jump-if-false: if the first parameter is zero, it sets the instruction pointer to the value from the second
 *    parameter. Otherwise, it does nothing.
 *  - Opcode 7 is less than: if the first parameter is less than the second parameter, it stores 1 in the position given by the
 *    third parameter. Otherwise, it stores 0.
 *  - Opcode 8 is equals: if the first parameter is equal to the second parameter, it stores 1 in the position given by the third
 *    parameter. Otherwise, it stores 0.
 *
 * Like all instructions, these instructions need to support parameter modes as described above.
 *
 * Normally, after an instruction is finished, the instruction pointer increases by the number of values in that instruction.
 * However, if the instruction modifies the instruction pointer, that value is used and the instruction pointer is not automatically increased.
 *
 * For example, here are several programs that take one input, compare it to the value 8, and then produce one output:
 *
 *  - 3,9,8,9,10,9,4,9,99,-1,8 - Using position mode, consider whether the input is equal to 8; output 1 (if it is) or 0 (if it is not).
 *  - 3,9,7,9,10,9,4,9,99,-1,8 - Using position mode, consider whether the input is less than 8; output 1 (if it is) or 0 (if it is not).
 *  - 3,3,1108,-1,8,3,4,3,99 - Using immediate mode, consider whether the input is equal to 8; output 1 (if it is) or 0 (if it is not).
 *  - 3,3,1107,-1,8,3,4,3,99 - Using immediate mode, consider whether the input is less than 8; output 1 (if it is) or 0 (if it is not).
 *
 * Here are some jump tests that take an input, then output 0 if the input was zero or 1 if the input was non-zero:
 *
 *  - 3,12,6,12,15,1,13,14,13,4,13,99,-1,0,1,9 (using position mode)
 *  - 3,3,1105,-1,9,1101,0,0,12,4,12,99,1 (using immediate mode)
 *
 * Here's a larger example:
 *
 * 3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,
 * 1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,
 * 999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99
 *
 * The above example program uses an input instruction to ask for a single number. The program will then output 999 if the input value
 * is below 8, output 1000 if the input value is equal to 8, or output 1001 if the input value is greater than 8.
 *
 * This time, when the TEST diagnostic program runs its input instruction to get the ID of the system to test, provide it 5, the ID for
 * the ship's thermal radiator controller. This diagnostic test suite only outputs one number, the diagnostic code.
 *
 * What is the diagnostic code for system ID 5?
 *
 */

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type IntCode struct {
	code               []int
	instructionPointer int

	inputBuffer int
	Debug       bool
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

	return IntCode{code: code}
}

func (intCode *IntCode) PrintDebug(additionalInfo string) {
	if !intCode.Debug {
		return
	}

	start := intCode.instructionPointer
	end := int(math.Min(float64(len(intCode.code)-1), float64(intCode.instructionPointer+10)))
	log.Printf(
		"%s instructionPointer=%d inputBuffer=%d current=%v",
		additionalInfo,
		intCode.instructionPointer,
		intCode.inputBuffer,
		intCode.code[start:end],
	)
}

func (intCode *IntCode) Get(address int) int {
	return intCode.code[address]
}

func (intCode *IntCode) Set(address int, value int) {
	intCode.code[address] = value
}

func (intCode *IntCode) getCurrentInstruction() (instruction int, paramModes int) {
	instructionCode := intCode.Get(intCode.instructionPointer)
	instruction = instructionCode % 100
	paramModes = int(math.Floor(float64(instructionCode / 100)))
	return
}

func (intCode *IntCode) parametersForCurrentInstruction(length int, paramModes int) (parameters []int, rawParameters []int) {
	start := intCode.instructionPointer + 1
	end := start + length

	rawParameters = intCode.code[start:end]
	parameters = make([]int, len(rawParameters))

	for i, param := range rawParameters {
		paramMode := paramModes / int(math.Pow10(i)) % 10
		switch paramMode {
		case 0:
			// position mode
			parameters[i] = intCode.Get(param)
			break
		case 1:
			// immediate mode: nothing changes
			parameters[i] = param
			break
		default:
			log.Fatalf("Unknown parameter mode %d for parameter %d of intcode at position %d", paramMode, i, intCode.instructionPointer)
		}
	}

	return
}

func (intCode *IntCode) incrementInstructionPointerBasedOnNumberOfParameters(params []int) {
	intCode.instructionPointer += len(params) + 1
}

func (intCode *IntCode) RunStep() (isDone bool, err error) {
	intCode.PrintDebug("")
	instruction, paramModes := intCode.getCurrentInstruction()

	// Note about params and raw: params are the parameters with the parameter mode applied to it, raw without.
	// In most cases (e.g. calculation, comparison) you will want to use the params with parameter mode.
	// Only use raw if you need to set an offset to update the intcode on the fly.
	var params, raw []int
	isDone = false

	switch instruction {
	case 1:
		// Opcode 1 adds together numbers read from two positions and stores the result in a third position.
		// The three integers immediately after the opcode tell you these three positions - the first two indicate
		// the positions from which you should read the input values, and the third indicates the position at which
		// the output should be stored.
		params, raw = intCode.parametersForCurrentInstruction(3, paramModes)
		intCode.Set(raw[2], params[0]+params[1])
		intCode.incrementInstructionPointerBasedOnNumberOfParameters(params)
		break
	case 2:
		// Opcode 2 works exactly like opcode 1, except it multiplies the two inputs instead of adding them.
		// Again, the three integers after the opcode indicate where the inputs and outputs are, not their values.
		params, raw = intCode.parametersForCurrentInstruction(3, paramModes)
		intCode.Set(raw[2], params[0]*params[1])
		intCode.incrementInstructionPointerBasedOnNumberOfParameters(params)
		break
	case 3:
		// Opcode 3 takes a single integer as input and saves it to the position given by its only parameter.
		// For example, the instruction 3,50 would take an input value and store it at address 50.
		params, raw = intCode.parametersForCurrentInstruction(1, paramModes)
		intCode.Set(raw[0], intCode.inputBuffer)
		intCode.incrementInstructionPointerBasedOnNumberOfParameters(params)
		break
	case 4:
		// Opcode 4 outputs the value of its only parameter. For example, the instruction 4,50 would output the value at address 50.
		params, _ = intCode.parametersForCurrentInstruction(1, paramModes)
		intCode.inputBuffer = params[0]
		intCode.incrementInstructionPointerBasedOnNumberOfParameters(params)
		break
	case 5:
		// Opcode 5 is jump-if-true: if the first parameter is non-zero, it sets the instruction pointer to the
		// value from the second parameter. Otherwise, it does nothing.
		params, _ = intCode.parametersForCurrentInstruction(2, paramModes)
		if params[0] != 0 {
			intCode.instructionPointer = params[1]
		} else {
			intCode.incrementInstructionPointerBasedOnNumberOfParameters(params)
		}
		break
	case 6:
		// Opcode 6 is jump-if-false: if the first parameter is zero, it sets the instruction pointer to the value
		// from the second parameter. Otherwise, it does nothing.
		params, _ = intCode.parametersForCurrentInstruction(2, paramModes)
		if params[0] == 0 {
			intCode.instructionPointer = params[1]
		} else {
			intCode.incrementInstructionPointerBasedOnNumberOfParameters(params)
		}
		break
	case 7:
		// Opcode 7 is less than: if the first parameter is less than the second parameter, it stores 1 in the position
		// given by the third parameter. Otherwise, it stores 0.
		params, raw = intCode.parametersForCurrentInstruction(3, paramModes)
		if params[0] < params[1] {
			intCode.Set(raw[2], 1)
		} else {
			intCode.Set(raw[2], 0)
		}
		intCode.incrementInstructionPointerBasedOnNumberOfParameters(params)
		break
	case 8:
		// Opcode 8 is equals: if the first parameter is equal to the second parameter, it stores 1 in the position
		// given by the third parameter. Otherwise, it stores 0.
		params, raw = intCode.parametersForCurrentInstruction(3, paramModes)
		if params[0] == params[1] {
			intCode.Set(raw[2], 1)
		} else {
			intCode.Set(raw[2], 0)
		}
		intCode.incrementInstructionPointerBasedOnNumberOfParameters(params)
		break
	case 99:
		// Opcode 99 terminates the program
		isDone = true
	default:
		err = errors.New(fmt.Sprintf("Invalid intcode instruction %v encountered", instruction))
	}

	return
}

func (intCode *IntCode) Run(input int) (int, error) {
	intCode.inputBuffer = input
	intCode.instructionPointer = 0

	i := 1
	for true {
		isDone, err := intCode.RunStep()

		if err != nil {
			return intCode.inputBuffer, err
		}
		if isDone {
			return intCode.inputBuffer, nil
		}

		i++
	}

	return intCode.inputBuffer, errors.New("No exit instruction (99) encountered")
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	code := NewIntCode(string(input))
	// code.Debug = true
	result, err := code.Run(5)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
