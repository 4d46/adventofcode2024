package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Cost struct {
	x int
	y int
}

type PrizeInstruction struct {
	buttonA   Cost
	buttonB   Cost
	prizeLoc  Cost
	pressesA  int
	pressesB  int
	totalCost int
	valid     bool
}

const (
	epsilon                   = 1e-9
	buttonACostPerPress       = 3
	buttonBCostPerPress       = 1
	part2UnitConversionFactor = 10000000000000
)

func main() {
	fmt.Println("Hello, day 13!")

	// Part 1
	// Load the input data
	inputData := LoadInputData("input1.txt")
	// inputData := LoadInputData("test1.txt")

	{
		// The strings in the input data are a set of instructions
		// Each instruction is in the format:
		// Button A: X+94, Y+34
		// Button B: X+22, Y+67
		// Prize: X=8400, Y=5400
		// With a blank line separting each instruction

		// We need to use a regular expression to extract the instruction strings
		regexButtonA := `Button A: X\+(\d+), Y\+(\d+)`
		regexButtonB := `Button B: X\+(\d+), Y\+(\d+)`
		regexPrize := `Prize: X=(\d+), Y=(\d+)`
		m1 := regexp.MustCompile(regexButtonA)
		m2 := regexp.MustCompile(regexButtonB)
		m3 := regexp.MustCompile(regexPrize)

		allInstructions := make([]PrizeInstruction, 0)

		currentInstruction := PrizeInstruction{}
		// Loop through the input data and find all instances of the instruction strings
		for _, line := range inputData {
			if strings.HasPrefix(line, "Button A") {
				// Find the instruction strings by running the regular expression on the line
				instructions := m1.FindAllStringSubmatch(line, -1)
				// Convert the strings to integers
				butA, _ := strconv.Atoi(instructions[0][1])
				butB, _ := strconv.Atoi(instructions[0][2])
				// Append the instruction strings to the allInstructions slice
				currentInstruction.buttonA.x = butA
				currentInstruction.buttonA.y = butB
			} else if strings.HasPrefix(line, "Button B") {
				// Find the instruction strings by running the regular expression on the line
				instructions := m2.FindAllStringSubmatch(line, -1)
				// Convert the strings to integers
				butA, _ := strconv.Atoi(instructions[0][1])
				butB, _ := strconv.Atoi(instructions[0][2])
				// Append the instruction strings to the allInstructions slice
				currentInstruction.buttonB.x = butA
				currentInstruction.buttonB.y = butB
			} else if strings.HasPrefix(line, "Prize") {
				// Find the instruction strings by running the regular expression on the line
				instructions := m3.FindAllStringSubmatch(line, -1)
				// Convert the strings to integers
				prizeA, _ := strconv.Atoi(instructions[0][1])
				prizeB, _ := strconv.Atoi(instructions[0][2])
				// Append the instruction strings to the allInstructions slice
				currentInstruction.prizeLoc.x = prizeA
				currentInstruction.prizeLoc.y = prizeB
			} else if line == "" {
				allInstructions = append(allInstructions, currentInstruction)
				currentInstruction = PrizeInstruction{}
			} else {
				log.Fatalf("unexpected line: %s", line)
			}
		}
		// spew.Dump(allInstructions)

		// Part 1
		// Solve the 2 simultaneous equations for x and y
		for i, in := range allInstructions {
			// Solve for x
			pressesB := float64((in.buttonA.x*in.prizeLoc.y)-(in.buttonA.y*in.prizeLoc.x)) / float64((in.buttonA.x*in.buttonB.y)-(in.buttonB.x*in.buttonA.y))
			pressesA := (float64(in.prizeLoc.y) - (float64(in.buttonB.y) * pressesB)) / float64(in.buttonA.y)
			totalCost := (pressesA * buttonACostPerPress) + (pressesB * buttonBCostPerPress)
			valid := true
			if _, frac := math.Modf(math.Abs(pressesA)); frac < epsilon || frac > 1.0-epsilon {
				allInstructions[i].pressesA = int(pressesA)
			} else {
				valid = false
			}
			if _, frac := math.Modf(math.Abs(pressesB)); frac < epsilon || frac > 1.0-epsilon {
				allInstructions[i].pressesB = int(pressesB)
			} else {
				valid = false
			}
			if _, frac := math.Modf(math.Abs(totalCost)); frac < epsilon || frac > 1.0-epsilon {
				allInstructions[i].totalCost = int(totalCost)
			} else {
				valid = false
			}
			allInstructions[i].valid = valid
		}

		// Find cost of winning all possible prizes
		totalCost := 0
		for _, in := range allInstructions {
			if in.valid {
				totalCost += in.totalCost
			}
		}

		fmt.Printf("Part 1: Total cost of winning all possible prizes: %d\n", totalCost)
	}

	fmt.Println("-----------------------------")

	{
		// The strings in the input data are a set of instructions
		// Each instruction is in the format:
		// Button A: X+94, Y+34
		// Button B: X+22, Y+67
		// Prize: X=8400, Y=5400
		// With a blank line separting each instruction

		// We need to use a regular expression to extract the instruction strings
		regexButtonA := `Button A: X\+(\d+), Y\+(\d+)`
		regexButtonB := `Button B: X\+(\d+), Y\+(\d+)`
		regexPrize := `Prize: X=(\d+), Y=(\d+)`
		m1 := regexp.MustCompile(regexButtonA)
		m2 := regexp.MustCompile(regexButtonB)
		m3 := regexp.MustCompile(regexPrize)

		allInstructions := make([]PrizeInstruction, 0)

		currentInstruction := PrizeInstruction{}
		// Loop through the input data and find all instances of the instruction strings
		for _, line := range inputData {
			if strings.HasPrefix(line, "Button A") {
				// Find the instruction strings by running the regular expression on the line
				instructions := m1.FindAllStringSubmatch(line, -1)
				// Convert the strings to integers
				butA, _ := strconv.Atoi(instructions[0][1])
				butB, _ := strconv.Atoi(instructions[0][2])
				// Append the instruction strings to the allInstructions slice
				currentInstruction.buttonA.x = butA
				currentInstruction.buttonA.y = butB
			} else if strings.HasPrefix(line, "Button B") {
				// Find the instruction strings by running the regular expression on the line
				instructions := m2.FindAllStringSubmatch(line, -1)
				// Convert the strings to integers
				butA, _ := strconv.Atoi(instructions[0][1])
				butB, _ := strconv.Atoi(instructions[0][2])
				// Append the instruction strings to the allInstructions slice
				currentInstruction.buttonB.x = butA
				currentInstruction.buttonB.y = butB
			} else if strings.HasPrefix(line, "Prize") {
				// Find the instruction strings by running the regular expression on the line
				instructions := m3.FindAllStringSubmatch(line, -1)
				// Convert the strings to integers
				prizeA, _ := strconv.Atoi(instructions[0][1])
				prizeB, _ := strconv.Atoi(instructions[0][2])
				// Append the instruction strings to the allInstructions slice
				currentInstruction.prizeLoc.x = prizeA + part2UnitConversionFactor
				currentInstruction.prizeLoc.y = prizeB + part2UnitConversionFactor
			} else if line == "" {
				allInstructions = append(allInstructions, currentInstruction)
				currentInstruction = PrizeInstruction{}
			} else {
				log.Fatalf("unexpected line: %s", line)
			}
		}
		// spew.Dump(allInstructions)

		// Part 1
		// Solve the 2 simultaneous equations for x and y
		for i, in := range allInstructions {
			// Solve for x
			pressesB := float64((in.buttonA.x*in.prizeLoc.y)-(in.buttonA.y*in.prizeLoc.x)) / float64((in.buttonA.x*in.buttonB.y)-(in.buttonB.x*in.buttonA.y))
			pressesA := (float64(in.prizeLoc.y) - (float64(in.buttonB.y) * pressesB)) / float64(in.buttonA.y)
			totalCost := (pressesA * buttonACostPerPress) + (pressesB * buttonBCostPerPress)
			valid := true
			if _, frac := math.Modf(math.Abs(pressesA)); frac < epsilon || frac > 1.0-epsilon {
				allInstructions[i].pressesA = int(pressesA)
			} else {
				valid = false
			}
			if _, frac := math.Modf(math.Abs(pressesB)); frac < epsilon || frac > 1.0-epsilon {
				allInstructions[i].pressesB = int(pressesB)
			} else {
				valid = false
			}
			if _, frac := math.Modf(math.Abs(totalCost)); frac < epsilon || frac > 1.0-epsilon {
				allInstructions[i].totalCost = int(totalCost)
			} else {
				valid = false
			}
			allInstructions[i].valid = valid
		}

		// fmt.Printf("Part 2: %v\n", allInstructions)

		// Find cost of winning all possible prizes
		totalCost := 0
		for _, in := range allInstructions {
			if in.valid {
				totalCost += in.totalCost
			}
		}

		fmt.Printf("Part 2: Total cost of winning all possible prizes: %d\n", totalCost)
	}
}

// LoadInputData reads the input file and returns a slice of strings
func LoadInputData(filename string) []string {
	// Read the file as strings a line at a time
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	// Create a new reader
	reader := bufio.NewReader(file)

	var inputData []string

	for {
		// Read until we encounter a newline character
		line, err := reader.ReadString('\n')
		if err != nil {
			// If we encounter EOF, break out of the loop
			if err.Error() == "EOF" {
				break
			}
			log.Fatalf("error reading file: %s", err)
		}
		// Remove the newline character from the end of the line
		line = strings.TrimSuffix(line, "\n")
		// Append the line to the inputData slice
		inputData = append(inputData, line)
	}
	return inputData
}
