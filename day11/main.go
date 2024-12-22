package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type stone int

func main() {
	fmt.Println("Hello, day 05!")

	// Part 1
	// Load the input data
	inputData := LoadInputData("input1.txt")
	// inputData := LoadInputData("test1.txt")

	// Parse the input data
	initialStones := make([]stone, 0)
	stoneStrings := strings.Split(inputData[0], " ")
	for _, s := range stoneStrings {
		value, _ := strconv.Atoi(s)
		initialStones = append(initialStones, stone(value))
	}

	// Build a lookup table for the stone numbers
	lookupTable := make(map[stone][]stone)

	// Loop over the first few Blinks to check how it works
	for blinks := 1; blinks <= 25; blinks++ {

		stoneCount := 0
		// Loop over the stones
		for _, s := range initialStones {
			// Follow the stone
			stoneCount += followStone(blinks-1, s, lookupTable)
		}

		fmt.Printf("Part 1: Blinks %d: %d\n", blinks, stoneCount)
	}

	fmt.Println("--------------------")

	// Part 2
	blinks := 75
	stoneCount := 0
	// Loop over the stones
	for _, s := range initialStones {
		// Follow the stone
		stoneCount += followStone(blinks-1, s, lookupTable)
	}

	fmt.Printf("Part 2: Blinks %d: %d\n", blinks, stoneCount)

}

func followStone(depth int, stone stone, lookupTable map[stone][]stone) int {
	// Check if the stone is in the lookup table, otherwise calculate it
	if _, ok := lookupTable[stone]; !ok {
		// Calculate the next stones
		calcStones := deriveNextStones(stone)
		// Add the stones to the lookup table
		lookupTable[stone] = calcStones
	}

	nextStones := lookupTable[stone]

	if depth == 0 {
		return len(nextStones)
	}

	// Follow the stones
	totalStones := 0
	for _, s := range nextStones {
		totalStones += followStone(depth-1, s, lookupTable)
	}
	return totalStones
}

func deriveNextStones(s stone) []stone {
	// Calculate the next stones
	nextStones := make([]stone, 0)

	if s == 0 {
		// Stone is 0 so replace with a stone with 1
		nextStones = append(nextStones, 1)
	} else if (int(math.Log10(float64(s))+1) % 2) == 0 {
		// Stone has an even number of digits, so split it into 2 stones
		stoneStr := strconv.Itoa(int(s))
		var newStone [2]int
		newStone[0], _ = strconv.Atoi(stoneStr[:len(stoneStr)/2])
		newStone[1], _ = strconv.Atoi(stoneStr[len(stoneStr)/2:])
		nextStones = append(nextStones, stone(newStone[0]), stone(newStone[1]))
	} else {
		// Otherwise, replace stone with a stone with 2024 times the value
		nextStones = append(nextStones, s*2024)
	}
	return nextStones
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
