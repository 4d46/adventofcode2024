package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type position struct {
	x int
	y int
}

type direction int

// Word directions
const (
	Up direction = iota
	Down
	Left
	Right
	UpLeft
	UpRight
	DownLeft
	DownRight
)

type wordLocation struct {
	start position
	dir   direction
}

func main() {
	fmt.Println("Hello, day 04!")

	// Part 1
	// Load the input data
	inputData := LoadInputData("input1.txt")
	// inputData := LoadInputData("test1.txt")

	// The strings in the input data are a word search grid
	// Convert the strings to a 2D slice of runes
	var grid [][]rune
	for _, line := range inputData {
		grid = append(grid, []rune(line))
	}
	{
		// Print the grid
		// spew.Dump(grid)
		// PrintGrid(grid)

		// Find all the first letters of the word "XMAS"
		startLocs := make([]position, 0)
		for y, row := range grid {
			for x, cell := range row {
				if cell == 'X' {
					startLocs = append(startLocs, position{x, y})
				}
			}
		}

		// spew.Dump(startLocs)

		wordLocations := make([]wordLocation, 0)

		// For each starting location, check all 8 directions for the word "XMAS"
		for _, startLoc := range startLocs {
			// fmt.Printf("Start location: %v\n", startLoc)
			for dir := Up; dir <= DownRight; dir++ {
				// fmt.Printf("Direction: %d\n", dir)
				// Check the direction
				if hasWordXmas(grid, startLoc, dir) {
					wordLocations = append(wordLocations, wordLocation{startLoc, dir})
				}
			}
		}

		// spew.Dump(wordLocations)

		fmt.Printf("Part 1: found %d occurrences of the word 'XMAS'\n", len(wordLocations))
	}

	// Part 2
	fmt.Println("-----------------------------")
	{
		// We are now looking for the word "MAS"
		// Find all the first letters of the word "MAS"
		startLocs := make([]position, 0)
		for y, row := range grid {
			for x, cell := range row {
				if cell == 'M' {
					startLocs = append(startLocs, position{x, y})
				}
			}
		}

		// spew.Dump(startLocs)

		wordLocations := make([]wordLocation, 0)
		for _, startLoc := range startLocs {
			// fmt.Printf("Start location: %v\n", startLoc)
			// We only want to check the 4 diagonal directions
			for dir := UpLeft; dir <= DownRight; dir++ {
				// fmt.Printf("Direction: %d\n", dir)
				if hasWordMas(grid, startLoc, dir) {
					wordLocations = append(wordLocations, wordLocation{startLoc, dir})
				}
			}
		}

		// Now we have all the locations of the diagonal word "MAS" we need to find where they intersect
		masCenters := make(map[position]([]wordLocation))
		for _, loc := range wordLocations {
			// Calculate the center of the word "MAS"
			center := iteratePosition(loc.start, loc.dir)
			masCenters[center] = append(masCenters[center], loc)
		}

		// Now we have all the intersections, we need to find ones that have 2 intersections
		// These are the number of X-MAS we have
		x_masCount := 0
		for _, locs := range masCenters {
			if len(locs) == 2 {
				x_masCount++
			} else if len(locs) > 2 {
				log.Fatalf("Found more than 2 intersections for MAS center: %v", locs)
			}
		}

		fmt.Printf("Part 2: found %d occurrences of the word 'X-MAS'\n", x_masCount)
	}

}

// Check if the word "XMAS" is present in the grid
func hasWordXmas(grid [][]rune, start position, dir direction) bool {
	// The word we are looking for
	word := []rune("XMAS")

	// The current position
	pos := start

	// Grid dimensions
	width := len(grid[0])
	height := len(grid)

	// Iterate through the word
	// We know the first letter is already present
	for i := 1; i < len(word); i++ {
		// Iterate the position based on the direction
		pos = iteratePosition(pos, dir)
		// Check if the position is within the grid, if not nothing found
		if !isPositionValid(width, height, pos) {
			return false
		}
		// Check if the letter at the position doesn't match the word
		if grid[pos.y][pos.x] != word[i] {
			return false
		}
	}

	// Reached the end of the word, found it
	return true
}

// Check if the word "MAS" is present in the grid
func hasWordMas(grid [][]rune, start position, dir direction) bool {
	// The word we are looking for
	word := []rune("MAS")

	// The current position
	pos := start

	// Grid dimensions
	width := len(grid[0])
	height := len(grid)

	// Iterate through the word
	// We know the first letter is already present
	for i := 1; i < len(word); i++ {
		// Iterate the position based on the direction
		pos = iteratePosition(pos, dir)
		// Check if the position is within the grid, if not nothing found
		if !isPositionValid(width, height, pos) {
			return false
		}
		// Check if the letter at the position doesn't match the word
		if grid[pos.y][pos.x] != word[i] {
			return false
		}
	}

	// Reached the end of the word, found it
	return true
}

// Iterate the position based on the direction
func iteratePosition(pos position, dir direction) position {
	switch dir {
	case Up:
		pos.y--
	case Down:
		pos.y++
	case Left:
		pos.x--
	case Right:
		pos.x++
	case UpLeft:
		pos.x--
		pos.y--
	case UpRight:
		pos.x++
		pos.y--
	case DownLeft:
		pos.x--
		pos.y++
	case DownRight:
		pos.x++
		pos.y++
	}
	return pos
}

// Check if the position is within the grid
func isPositionValid(width, height int, pos position) bool {
	if pos.x < 0 || pos.x >= width {
		return false
	}
	if pos.y < 0 || pos.y >= height {
		return false
	}
	return true
}

// PrintGrid prints the grid to the console
func PrintGrid(grid [][]rune) {
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
