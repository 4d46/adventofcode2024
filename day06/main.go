package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type GuardPosition struct {
	x   int
	y   int
	dir Direction
}

type Position struct {
	x int
	y int
}

func main() {
	fmt.Println("Hello, day 06!")

	// Part 1
	// Load the input data
	inputData := LoadInputData("input1.txt")
	// inputData := LoadInputData("test1.txt")

	// Create a grid to store the obstructions
	grid := make([][]rune, len(inputData))
	var guardStart GuardPosition
	for i := range grid {
		grid[i] = []rune(inputData[i])
		// Check for guard start position
		for j := range grid[i] {
			if grid[i][j] != '.' && grid[i][j] != '#' {
				// Found a guard
				fmt.Printf("Guard found at (%d,%d)\n", j, i)
				guardStart.x = j
				guardStart.y = i
				switch grid[i][j] {
				case 'v':
					guardStart.dir = Down
				case '^':
					guardStart.dir = Up
				case '<':
					guardStart.dir = Left
				case '>':
					guardStart.dir = Right
				}
				// Replace the guard with a '.'
				grid[i][j] = '.'
			}
		}
	}

	// Print out the grid
	printGrid(grid)
	spew.Dump(guardStart)

	// Now we are going to walk the guard forward until they start to repeat their path
	steps := make(map[GuardPosition]bool)
	steps[guardStart] = true
	currPos := guardStart
REPEATINGPATH:
	for {
		// Try all directions, turn right if we can't go forward
		// First direction that works, take it
		for i := 0; i < 4; i++ {
			// Find the next position
			nextPos := calculateNextPosition(currPos)

			// Check if we have been here before (in the same direction)
			if steps[nextPos] {
				fmt.Println("Found a loop at", nextPos)
				break REPEATINGPATH
			}

			// If the next position is off the grid the guard has left the area
			if !validPosition(grid, nextPos) {
				fmt.Println("Guard has left the area")
				break REPEATINGPATH
			}
			// If the next position is an obstruction try next direction
			if grid[nextPos.y][nextPos.x] == '#' {
				currPos.dir = turnRight(currPos.dir)
				continue
			}
			currPos = nextPos
			steps[currPos] = true
			break
		}
		// spew.Dump(currPos)
	}

	// Print the grid with the steps
	printGridWithSteps(grid, steps)

	// Calculate the number of positions
	uniquePositions := make(map[GuardPosition]bool)
	for pos := range steps {
		pos.dir = Up
		uniquePositions[pos] = true
	}

	fmt.Println("Part 1: Number of steps before repeating path:", len(uniquePositions))

	fmt.Println("-----------------------------")

	// Part 2
	// Now we want to find places where the guard could be diverted into a loop
	// This would only be on one of the path points we have already visited in Part 1
	// Loop over all previous steps and check if we can divert the guard into a loop
	// by placing an obstruction at that point
	loopFormingObstructions := make(map[Position]bool)
	for pos := range uniquePositions {
		// Place an obstruction at this point
		extraObstruction := Position{x: pos.x, y: pos.y}
		// Now we are going to walk the guard forward until they start to repeat their path
		stepsP2 := make(map[GuardPosition]bool, len(steps))
		stepsP2[guardStart] = true
		currPos := guardStart
	REPEATINGPATH2:
		for {
			// Try all directions, turn right if we can't go forward
			// First direction that works, take it
			for i := 0; i < 4; i++ {
				// Find the next position
				nextPos := calculateNextPosition(currPos)

				// Check if we have been here before (in the same direction)
				if stepsP2[nextPos] {
					fmt.Println("Found a loop at", nextPos)
					loopFormingObstructions[extraObstruction] = true
					break REPEATINGPATH2
				}

				// If the next position is off the grid the guard has left the area
				// So not loop forming, try next obstruction location
				if !validPosition(grid, nextPos) {
					// fmt.Println("Guard has left the area, so not loop forming")
					break REPEATINGPATH2
				}
				// If the next position is an obstruction try next direction
				if grid[nextPos.y][nextPos.x] == '#' || (nextPos.x == extraObstruction.x && nextPos.y == extraObstruction.y) {
					currPos.dir = turnRight(currPos.dir)
					continue
				}
				currPos = nextPos
				stepsP2[currPos] = true
				break
			}
			// spew.Dump(currPos)
		}
	}
	fmt.Printf("Part 2: Number of loop forming obstructions: %d\n", len(loopFormingObstructions))
}

func turnRight(dir Direction) Direction {
	switch dir {
	case Up:
		return Right
	case Down:
		return Left
	case Left:
		return Up
	case Right:
		return Down
	}
	return -1
}

func calculateNextPosition(currPos GuardPosition) GuardPosition {
	nextPos := currPos
	switch currPos.dir {
	case Up:
		nextPos.y--
	case Down:
		nextPos.y++
	case Left:
		nextPos.x--
	case Right:
		nextPos.x++
	}
	return nextPos
}

func validPosition(grid [][]rune, pos GuardPosition) bool {
	if pos.y < 0 || pos.y >= len(grid) {
		return false
	}
	if pos.x < 0 || pos.x >= len(grid[pos.y]) {
		return false
	}
	return true
}

func printGrid(grid [][]rune) {
	for _, row := range grid {
		fmt.Println(string(row))
	}
}

func printGridWithSteps(grid [][]rune, steps map[GuardPosition]bool) {
	for step := range steps {
		grid[step.y][step.x] = 'X'
	}
	for _, row := range grid {
		fmt.Println(string(row))
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
