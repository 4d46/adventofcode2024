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

var lookupTable map[string]int

// Still having trouble getting my head round memoisation.  My first attempt a half way house, but not fast enough
// Researched and found this https://github.com/jimlawless/aoc2024/blob/main/day_11/day_11b.go
// Reworked to use that approach instead and hopefully get my head around it

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

	// Build a lookup table to support memoisation
	// One map for each level.  As the recursion dives to the bottom it will populate the number of children
	// at each level, reducing the number of lookups required
	lookupTable = make(map[string]int)

	// Loop over the first few Blinks to check how it works
	for blinks := 1; blinks <= 25; blinks++ {

		stoneCount := 0
		// Loop over the stones
		for _, s := range initialStones {
			// Follow the stone
			stoneCount += followStone(blinks, s)
		}
		// spew.Dump(lookupTable)

		fmt.Printf("Part 1: Blinks %d: %d\n", blinks, stoneCount)
	}

	fmt.Println("--------------------")

	// Part 2
	blinks := 75
	stoneCount := 0
	// Loop over the stones
	for _, s := range initialStones {
		// Follow the stone
		stoneCount += followStone(blinks, s)
	}

	fmt.Printf("Part 2: Blinks %d: %d\n", blinks, stoneCount)

}

func followStone(depth int, st stone) int {
	// If we have hit the bottom of the recursion, we are delaing with a single stone
	// Return 1
	if depth == 0 {
		return 1
	}

	// Lets see if we have already calculated the number of child stones for this stone at this depth
	stoneKey := makeKey(st, depth)
	if _, ok := lookupTable[stoneKey]; ok {
		// We do, so return the number of child stones
		return lookupTable[stoneKey]
	}

	// If stone is 0, continue with 1
	if st == 0 {
		stones := followStone(depth-1, 1)
		lookupTable[stoneKey] = stones
		return stones
	} else if (int(math.Log10(float64(st))+1) % 2) == 0 {
		// Stone has an even number of digits, so split it into 2 stones
		stoneStr := fmt.Sprintf("%d", st)
		var newStone [2]int
		newStone[0], _ = strconv.Atoi(stoneStr[:len(stoneStr)/2])
		newStone[1], _ = strconv.Atoi(stoneStr[len(stoneStr)/2:])
		// Calculate the number of stones for each of the new stones
		stones := followStone(depth-1, stone(newStone[0])) + followStone(depth-1, stone(newStone[1]))
		lookupTable[stoneKey] = stones
		return stones
	} else {
		// Otherwise, replace stone with a stone with 2024 times the value
		stones := followStone(depth-1, st*2024)
		lookupTable[stoneKey] = stones
		return stones
	}
}

func makeKey(stone stone, depth int) string {
	return fmt.Sprintf("%d-%d", stone, depth)
}

// func deriveNextStones(s stone) []stone {
// 	// Calculate the next stones
// 	nextStones := make([]stone, 0)

// 	if s == 0 {
// 		// Stone is 0 so replace with a stone with 1
// 		nextStones = append(nextStones, 1)
// 	} else if (int(math.Log10(float64(s))+1) % 2) == 0 {
// 		// Stone has an even number of digits, so split it into 2 stones
// 		stoneStr := strconv.Itoa(int(s))
// 		var newStone [2]int
// 		newStone[0], _ = strconv.Atoi(stoneStr[:len(stoneStr)/2])
// 		newStone[1], _ = strconv.Atoi(stoneStr[len(stoneStr)/2:])
// 		nextStones = append(nextStones, stone(newStone[0]), stone(newStone[1]))
// 	} else {
// 		// Otherwise, replace stone with a stone with 2024 times the value
// 		nextStones = append(nextStones, s*2024)
// 	}
// 	return nextStones
// }

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
