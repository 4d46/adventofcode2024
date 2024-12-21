package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Position struct {
	x int
	y int
}

func main() {
	fmt.Println("Hello, day 08!")

	// Part 1
	// Load the input data
	inputData := LoadInputData("input1.txt")
	// inputData := LoadInputData("test1.txt")

	// Create a grid to store the obstructions
	grid := make([][]rune, len(inputData))
	for i := range grid {
		grid[i] = []rune(inputData[i])
	}

	// Print out the grid
	printGrid(grid)
	fmt.Println()

	// Store the antenna positions
	antennas := make(map[rune][]Position)
	for row := 0; row < len(grid); row++ {
		for column := 0; column < len(grid[0]); column++ {
			gridSymbol := grid[row][column]
			if gridSymbol != '.' {
				if _, ok := antennas[gridSymbol]; !ok {
					antennas[gridSymbol] = []Position{}
				}
				antennas[gridSymbol] = append(antennas[gridSymbol], Position{x: column, y: row})
			}

		}

	}
	// spew.Dump(antennas)
	{
		// Loop over all antenna types, calculating the antinodes for each pair
		antinodes := make(map[rune][]Position)
		for antennaType, positions := range antennas {
			// Check that this antenna type has an entry in the antinodes map
			if _, ok := antinodes[antennaType]; !ok {
				antinodes[antennaType] = []Position{}
			}
			// Loop over all combinations of position pairs
			for i := 0; i < len(positions); i++ {
				for j := i + 1; j < len(positions); j++ {
					// Calculate the delta from the first to the second node
					deltaX := positions[j].x - positions[i].x
					deltaY := positions[j].y - positions[i].y
					// Calculate the first antinode as negative delta from the first point
					antinodeX := positions[i].x - deltaX
					antinodeY := positions[i].y - deltaY
					// Calculate the second antinode as positive delta from the second point
					antinodeX2 := positions[j].x + deltaX
					antinodeY2 := positions[j].y + deltaY
					// Check the antinodes are within the grid and store if valid
					if antinodeX >= 0 && antinodeX < len(grid[0]) && antinodeY >= 0 && antinodeY < len(grid) {
						antinodes[antennaType] = append(antinodes[antennaType], Position{x: antinodeX, y: antinodeY})
					}
					if antinodeX2 >= 0 && antinodeX2 < len(grid[0]) && antinodeY2 >= 0 && antinodeY2 < len(grid) {
						antinodes[antennaType] = append(antinodes[antennaType], Position{x: antinodeX2, y: antinodeY2})
					}
				}
			}

		}

		// Print out the antinodes
		printAntinodes(grid, antinodes)

		// Find number of unique antinodes by adding to a map
		uniqueAntinodes := make(map[Position]bool)
		for _, positions := range antinodes {
			for _, position := range positions {
				uniqueAntinodes[position] = true
			}
		}

		fmt.Println("Part 1: Number of unique antinodes:", len(uniqueAntinodes))
	}
	fmt.Println("-----------------------------")

	// Part 2

	// printGrid(grid)
	{
		// Loop over all antenna types, calculating the antinodes for each pair
		// These antinodes now repeat until they are outside the grid
		antinodes := make(map[rune][]Position)
		for antennaType, positions := range antennas {
			// Check that this antenna type has an entry in the antinodes map
			if _, ok := antinodes[antennaType]; !ok {
				antinodes[antennaType] = []Position{}
			}
			// Loop over all combinations of position pairs
			for i := 0; i < len(positions); i++ {
				for j := i + 1; j < len(positions); j++ {
					// Calculate the delta from the first to the second node
					deltaX := positions[j].x - positions[i].x
					deltaY := positions[j].y - positions[i].y
					// Calculate the first antinode set as negative deltas from the first point
					for k := 1; ; k++ {
						antinodeX := positions[i].x - k*deltaX
						antinodeY := positions[i].y - k*deltaY
						// Check the antinode is within the grid and store if valid
						if antinodeX >= 0 && antinodeX < len(grid[0]) && antinodeY >= 0 && antinodeY < len(grid) {
							antinodes[antennaType] = append(antinodes[antennaType], Position{x: antinodeX, y: antinodeY})
						} else {
							break
						}
					}
					// Calculate the second antinode set as positive deltas from the second point
					for k := 1; ; k++ {
						antinodeX2 := positions[j].x + k*deltaX
						antinodeY2 := positions[j].y + k*deltaY
						// Check the antinode is within the grid and store if valid
						if antinodeX2 >= 0 && antinodeX2 < len(grid[0]) && antinodeY2 >= 0 && antinodeY2 < len(grid) {
							antinodes[antennaType] = append(antinodes[antennaType], Position{x: antinodeX2, y: antinodeY2})
						} else {
							break
						}
					}
				}
			}
		}

		// Print out the antinodes
		printAntinodes(grid, antinodes)

		// Find number of unique antinodes by adding to a map
		uniqueAntinodes := make(map[Position]bool)
		for _, positions := range antinodes {
			for _, position := range positions {
				uniqueAntinodes[position] = true
			}
		}
		// Also add the Antinodes located at the antenna positions
		for _, positions := range antennas {
			if len(positions) > 1 {
				for _, position := range positions {
					uniqueAntinodes[position] = true
				}
			}
		}

		fmt.Println("Part 2: Number of unique antinodes:", len(uniqueAntinodes))
	}

}

func printGrid(grid [][]rune) {
	for _, row := range grid {
		fmt.Println(string(row))
	}
}

func printAntinodes(grid [][]rune, antinodes map[rune][]Position) {
	newGrid := make([][]rune, len(grid))
	for i := range grid {
		newGrid[i] = make([]rune, len(grid[i]))
		copy(newGrid[i], grid[i])
	}

	for antinodetype, positions := range antinodes {
		for _, position := range positions {
			newGrid[position.y][position.x] = antinodetype
		}
	}

	// Print the grid
	printGrid(newGrid)
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
