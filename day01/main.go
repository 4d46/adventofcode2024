package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Hello, day 01!")

	// Part 1
	// Load the input data
	inputData := LoadInputData("input1.txt")
	// inputData := LoadInputData("test1.txt")

	// Loop through input data adding the integers in the first and second column into 2 separate arrays
	// Resize the arrays to the length of the input data
	var firstColumn []int
	var secondColumn []int
	firstColumn = make([]int, len(inputData))
	secondColumn = make([]int, len(inputData))

	// Input data is a slice of strings and the columns are space separated
	for i, v := range inputData {
		// Split the string into 2 parts
		parts := strings.Split(v, " ")
		// spew.Dump(parts)
		// Convert the strings to integers
		firstColumn[i], _ = strconv.Atoi(parts[0])
		secondColumn[i], _ = strconv.Atoi(parts[3])
		// fmt.Println(firstColumn[i], secondColumn[i], parts[0], parts[3])
	}

	// Sort the arrays
	sort.Ints(firstColumn)
	sort.Ints(secondColumn)

	// Find the sum of the distance between each of the rows in the sorted columns
	var sum int
	for i := 0; i < len(firstColumn); i++ {
		sum += Difference(firstColumn[i], secondColumn[i])
		// fmt.Println(firstColumn[i], secondColumn[i], Difference(firstColumn[i], secondColumn[i]))
	}

	fmt.Println("Part 1: ", sum)

	fmt.Println("--------------------")

	// Part 2
	// Input data is the same as part 1
	// Build lookup table for the second column.  Create a map for each number present and record the sum for every number present
	// Loop through the input data and build the lookup table
	lookupTable := make(map[int]int)
	for i := 0; i < len(secondColumn); i++ {
		lookupTable[secondColumn[i]] += 1
	}

	// Find the similarity score by looking up the value in the first column and multiplying it by the value in the lookup table
	// and summing the results
	var similarityScore int
	for i := 0; i < len(firstColumn); i++ {
		similarityScore += firstColumn[i] * lookupTable[firstColumn[i]]
	}

	fmt.Println("Part 2: ", similarityScore)

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

// Abs returns the absolute value of an integer
func Difference(x, y int) int {
	if x == y {
		return 0
	} else if x > y {
		return x - y
	} else {
		return y - x
	}
}
