package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Schematic struct {
	heights [5]int8
}

func main() {
	fmt.Println("Hello, day 25!")

	// Part 1
	// Load the input data
	inputData := LoadInputData("input1.txt")
	// inputData := LoadInputData("test1.txt")

	// Each block is a key or a lock schematic.  Parse each block to capture their heights
	locks := make([]Schematic, 0)
	keys := make([]Schematic, 0)

	// Loop through the input data and find all schematics
	schematicStrings := []string{}
	for _, line := range inputData {
		if line == "" {
			// We have reached the end of the schematic
			// Parse the schematic
			if strings.Contains(schematicStrings[0], ".....") {
				// This is a key schematic
				keys = append(keys, ParseKey(schematicStrings))
			} else if strings.Contains(schematicStrings[0], "#####") {
				// This is a lock schematic
				locks = append(locks, ParseLock(schematicStrings))
			} else {
				log.Fatalf("Unknown schematic type: %s", schematicStrings[0])
			}
			// Clear the schematicStrings slice
			schematicStrings = []string{}
		} else {
			// Append the line to the schematicStrings slice
			schematicStrings = append(schematicStrings, line)
		}
	}

	// fmt.Println("Locks:")
	// spew.Dump(locks)
	// fmt.Println("Keys:")
	// spew.Dump(keys)

	// Loop through the keys and locks and check if they fit
	keyLockFitCombinations := 0
	for _, lock := range locks {
		for _, key := range keys {
			// Check if the key fits the lock
			if key.fits(lock) {
				fmt.Printf("Lock %v fitted by key %v\n", lock, key)
				keyLockFitCombinations++
			}
		}
	}

	fmt.Printf("Part 1: There are %d key-lock fit combinations\n", keyLockFitCombinations)
}

// fits checks if the key fits the lock
func (key Schematic) fits(lock Schematic) bool {
	// Loop through the heights and check if the key fits the lock
	for i, height := range key.heights {
		// Check if the key height is less than the lock gap
		if height+lock.heights[i] > 0 {
			// The key does not fit the lock
			return false
		}
	}
	// The key fits the lock
	return true
}

// ParseKey parses a key schematic and returns a Schematic struct
func ParseKey(schematicStrings []string) Schematic {
	var schematic Schematic
	for i, line := range schematicStrings {
		// Ignore top and bottom lines
		if !(i == 0 || i == len(schematicStrings)-1) {
			// Loop through the line and capture the heights
			for j, char := range line {
				if char == '#' {
					// This is a height
					schematic.heights[j]++
				}
			}
		}
	}
	return schematic
}

// ParseLock parses a lock schematic and returns a Schematic struct
func ParseLock(schematicStrings []string) Schematic {
	var schematic Schematic
	for i, line := range schematicStrings {
		// Ignore top and bottom lines
		if !(i == 0 || i == len(schematicStrings)-1) {
			// Loop through the line and capture the gap as a negative height
			for j, char := range line {
				if char == '.' {
					// This is a gap
					schematic.heights[j]--
				}
			}
		}
	}
	return schematic
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
