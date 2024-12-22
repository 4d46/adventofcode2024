package main

import (
	"bufio"
	"fmt"
	"log"
	"maps"
	"os"
	"strconv"
	"strings"
)

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type Position struct {
	x int
	y int
}

func main() {
	fmt.Println("Hello, day 10!")

	// Part 1
	// Load the input data
	inputData := LoadInputData("input1.txt")
	// inputData := LoadInputData("test1.txt")
	// inputData := LoadInputData("test3.txt")
	// inputData := LoadInputData("test4.txt")
	// inputData := LoadInputData("test5.txt")
	// inputData := LoadInputData("test6.txt")

	// Create a grid to store the obstructions
	grid := make([][]byte, len(inputData))
	for i := range grid {
		grid[i] = make([]byte, len(inputData[i]))
		for j := 0; j < len(inputData[i]); j++ {
			if inputData[i][j] == '.' {
				grid[i][j] = 255
			} else {
				height, err := strconv.ParseInt(string(inputData[i][j]), 10, 8)
				if err != nil {
					log.Fatalf("failed to convert height to int: %s", err)
				}
				grid[i][j] = byte(height)
			}
		}
	}

	// Print out the grid
	printGrid(grid)

	// Find starting points, thos are the positions with height 0
	startingPoints := make([]Position, 0)
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			if grid[i][j] == 0 {
				startingPoints = append(startingPoints, Position{x: j, y: i})
			}
		}
	}

	fmt.Printf("Starting points: %v\n", startingPoints)

	// Walk through the grid, starting at the starting points
	validPaths := 0
	for _, startingPoint := range startingPoints {
		trailHeads := step(grid, startingPoint)
		validPaths += len(trailHeads)
	}

	fmt.Printf("Part 1: Number of valid paths: %d\n", validPaths)

	fmt.Println("-----------------------------")

	// Part 2

	// Walk through the grid, starting at the starting points, but this time finding all possible paths
	validPaths = 0
	for _, startingPoint := range startingPoints {
		validPaths += stepDistinct(grid, startingPoint)
	}

	fmt.Printf("Part 2: Number of valid paths: %d\n", validPaths)
}

// Function that steps along a path, returning the number of complete paths from this point
func step(grid [][]byte, current Position) map[Position]bool {
	// fmt.Printf("Height: %d, Position: %v\n", grid[current.x][current.y], current)

	// If we are off the grid exit
	if current.x < 0 || current.x >= len(grid[0]) || current.y < 0 || current.y >= len(grid) {
		log.Fatalf("off the grid")
	}
	// If we are at the end of a path, return 1
	if grid[current.y][current.x] == 9 {
		return map[Position]bool{current: true}
	}

	// Find next steps
	nextSteps := nextSteps(grid, current)
	// fmt.Printf("Position %v, next steps: %v\n", current, nextSteps)
	// Attempt to follow each step
	foundTrailHeads := make(map[Position]bool)
	for _, nextStep := range nextSteps {
		trailHeads := step(grid, nextStep)
		maps.Copy(foundTrailHeads, trailHeads)
	}
	return foundTrailHeads
}

// Function that steps along a path, returning the number of complete paths from this point
func stepDistinct(grid [][]byte, current Position) int {
	// fmt.Printf("Height: %d, Position: %v\n", grid[current.x][current.y], current)

	// If we are off the grid exit
	if current.x < 0 || current.x >= len(grid[0]) || current.y < 0 || current.y >= len(grid) {
		log.Fatalf("off the grid")
	}
	// If we are at the end of a path, return 1
	if grid[current.y][current.x] == 9 {
		return 1
	}

	// Find next steps
	nextSteps := nextSteps(grid, current)
	// fmt.Printf("Position %v, next steps: %v\n", current, nextSteps)
	// Attempt to follow each step
	validPaths := 0
	for _, nextStep := range nextSteps {
		validPaths += stepDistinct(grid, nextStep)
	}
	return validPaths
}

// Given a current point, retun all possible next steps
func nextSteps(grid [][]byte, current Position) []Position {
	var steps []Position
	nextHeight := grid[current.y][current.x] + 1
	// fmt.Printf("This pos: %v, This height: %d\n", current, grid[current.y][current.x])
	// fmt.Printf("Next height: %d\n", nextHeight)
	// Look for adjacent squares that are one higher
	// Check if we can move left
	if (current.x > 0) && (grid[current.y][current.x-1] == nextHeight) {
		steps = append(steps, Position{x: current.x - 1, y: current.y})
	}
	// Check if we can move right
	if (current.x < len(grid[0])-1) && (grid[current.y][current.x+1] == nextHeight) {
		steps = append(steps, Position{x: current.x + 1, y: current.y})
	}
	// Check if we can move up
	if (current.y > 0) && (grid[current.y-1][current.x] == nextHeight) {
		steps = append(steps, Position{x: current.x, y: current.y - 1})
	}
	// Check if we can move down
	if (current.y < len(grid)-1) && (grid[current.y+1][current.x] == nextHeight) {
		steps = append(steps, Position{x: current.x, y: current.y + 1})
	}

	return steps
}

func printGrid(grid [][]byte) {
	for _, row := range grid {
		fmt.Printf("%v\n", row)
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
