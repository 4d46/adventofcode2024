package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	fmt.Println("Hello, day 03!")

	// Part 1
	// Load the input data
	inputData := LoadInputData("input1.txt")
	// inputData := LoadInputData("test1.txt")
	// inputData := LoadInputData("test2.txt")

	{
		// The strings in the input data are a set of instructions
		// Each instruction is in the format "mul(1,20)", but there are lots of errors in the data too
		// We need to use a regular expression to extract the instruction strings
		regex := `mul\(\d+,\d+\)`
		m := regexp.MustCompile(regex)

		allInstructions := make([]string, 0)

		// Loop through the input data and find all instances of the instruction strings
		for _, line := range inputData {
			// Find the instruction strings by running the regular expression on the line
			instructions := m.FindAllString(line, -1)
			// Append the instruction strings to the allInstructions slice
			allInstructions = append(allInstructions, instructions...)
			// spew.Dump(instructions)
		}
		spew.Dump(allInstructions)

		// Now execute the instructions
		// The instructions are in the format "mul(1,20)"
		// We need to extract the 2 numbers from the instruction string
		// and multiply them together
		// The result of each multiplication should be added to a total
		// The total is the answer to the puzzle
		total := 0
		regex2 := `mul\((\d+),(\d+)\)`
		m2 := regexp.MustCompile(regex2)
		for _, instruction := range allInstructions {
			// Capture the 2 numbers from the instruction string using the regular expression
			matches := m2.FindStringSubmatch(instruction)
			// Convert the numbers to integers
			num1, err := strconv.Atoi(matches[1])
			if err != nil {
				log.Fatalf("failed to convert string to integer: %s", err)
			}
			num2, err := strconv.Atoi(matches[2])
			if err != nil {
				log.Fatalf("failed to convert string to integer: %s", err)
			}
			// Multiply the 2 numbers together and add the result to the total
			total += num1 * num2
		}

		fmt.Printf("Part 1: total = %d\n", total)
	}

	fmt.Println("-----------------------------")

	// Part 2
	// This time we need to identify and process additional conditional statements in the input data
	{
		regex := `mul\(\d+,\d+\)|do\(\)|don't\(\)`
		m := regexp.MustCompile(regex)

		allInstructions := make([]string, 0)

		// Loop through the input data and find all instances of the instruction strings
		for _, line := range inputData {
			// Find the instruction strings by running the regular expression on the line
			instructions := m.FindAllString(line, -1)
			// Append the instruction strings to the allInstructions slice
			allInstructions = append(allInstructions, instructions...)
			// spew.Dump(instructions)
		}
		spew.Dump(allInstructions)

		// Now execute the instructions
		enabled := true
		total := 0

		regex2 := `mul\((\d+),(\d+)\)`
		m2 := regexp.MustCompile(regex2)

		for _, instruction := range allInstructions {
			if instruction == "do()" {
				enabled = true
			} else if instruction == "don't()" {
				enabled = false
			} else if enabled {
				// Capture the 2 numbers from the instruction string using the regular expression
				matches := m2.FindStringSubmatch(instruction)
				// Convert the numbers to integers
				num1, err := strconv.Atoi(matches[1])
				if err != nil {
					log.Fatalf("failed to convert string to integer: %s", err)
				}
				num2, err := strconv.Atoi(matches[2])
				if err != nil {
					log.Fatalf("failed to convert string to integer: %s", err)
				}
				// Multiply the 2 numbers together and add the result to the total
				total += num1 * num2
			}
		}

		fmt.Printf("Part 2: total = %d\n", total)

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

// Difference returns the difference between 2 integers
func Difference(x, y int) int {
	if x == y {
		return 0
	} else if x > y {
		return x - y
	} else {
		return y - x
	}
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
