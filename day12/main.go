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
	Right
	Down
	Left
	UpLeft
	UpRight
	DownLeft
	DownRight
	Undefined
)

type region struct {
	plant     rune
	start     position
	clockdir  direction
	single    bool
	perimeter int
	area      int
	sides     int
	// innerRegions []region
}

func main() {
	fmt.Println("Hello, day 12!")

	// Part 1
	// Load the input data
	inputData := LoadInputData("input1.txt")
	// inputData := LoadInputData("test1.txt")
	// inputData := LoadInputData("test2.txt")
	// inputData := LoadInputData("test3.txt")
	// inputData := LoadInputData("test4.txt")
	// inputData := LoadInputData("test5.txt")

	// The strings in the input data are a word search grid
	// Convert the strings to a 2D slice of runes
	var garden [][]rune
	for _, line := range inputData {
		garden = append(garden, []rune(line))
	}
	printGarden(garden)

	regions := make(map[position]region, 0)

	// Start from all points in  the garden finding regions
	for y := 0; y < len(garden); y++ {
		for x := 0; x < len(garden[0]); x++ {
			pos := position{x, y}
			// Check in case we already know about this region
			if _, ok := regions[pos]; ok {
				continue
			}
			validRegion, nextRegion := walkPerimeter(garden, pos)
			if validRegion {
				// Add the region to the regions map if it doesn't already exist
				if _, ok := regions[nextRegion.start]; !ok {
					regions[nextRegion.start] = nextRegion
				}
			}
			fmt.Print(".")
		}
	}
	fmt.Println()

	printRegions(regions)
	{
		totalFencePrice := 0
		for _, reg := range regions {
			// Calculate the price of the fence
			fencePrice := reg.perimeter * reg.area
			totalFencePrice += fencePrice
		}
		fmt.Printf("Part 1: Total fence price: %d\n", totalFencePrice)
	}
	fmt.Println("--------------------")

	// Part 2
	{
		totalFencePrice := 0
		for _, reg := range regions {
			// Calculate the price of the fence
			fencePrice := reg.sides * reg.area
			totalFencePrice += fencePrice
		}
		fmt.Printf("Part 2: Total fence price: %d\n", totalFencePrice)
	}
}

// Walk permimeter of region, return whether it is valid and the region details
func walkPerimeter(garden [][]rune, pos position) (bool, region) {
	plant := garden[pos.y][pos.x]

	// // Check if this a single plant region by searching all adjacent cells for the same plant
	// // If we find a different plant, this is not a single plant region
	// singleRegion := true
	// for dir := Up; dir <= DownRight; dir++ {
	// 	adjPos := iteratePosition(pos, dir)
	// 	// Check if the adjacent position is within the bounds of the garden
	// 	if !(adjPos.x < 0 || adjPos.x >= len(garden[0]) || adjPos.y < 0 || adjPos.y >= len(garden)) {
	// 		if garden[adjPos.y][adjPos.x] == plant {
	// 			singleRegion = false
	// 			break
	// 		}
	// 	}
	// }
	// if singleRegion {
	// 	return region{plant: plant, start: pos, single: true, perimeter: 4, area: 1}
	// }

	walkPos := pos
	var walkDir direction

	// Not a single plant region, walk the perimeter
	// Remember the top, then leftmost position position of the boundary
	topLeft := pos
	// Now we need to find which side we are on to start walking the parmieter
	// Check the edges of the plot to find which edge to start walking in a clockwise direction
	foundDir := Undefined
	for dir := Up; dir <= Left; dir++ {
		// Check if the adjacent position is within the bounds of the garden
		checkPos := iteratePosition(pos, dir)
		if checkPos.x < 0 || checkPos.x >= len(garden[0]) || checkPos.y < 0 || checkPos.y >= len(garden) {
			foundDir = dir
			break
		}
		if garden[checkPos.y][checkPos.x] != plant {
			foundDir = dir
			break
		}
	}
	switch foundDir {
	case Up:
		walkDir = Right
	case Right:
		walkDir = Down
	case Down:
		walkDir = Left
	case Left:
		walkDir = Up
	case Undefined:
		// Not on an edge, so we must be in the middle of a region, skip this position
		return false, region{}
	}

	startPos := pos
	startDir := walkDir

	perimeterNodes := make(map[position][]direction)
	// Add the starting position to the perimeter nodes
	perimeterNodes[pos] = append(perimeterNodes[pos], walkDir)
	sideCount := 0

	// Walk the perimeter
	for {
		// If this position is more top or left than the current topLeft, update topLeft
		if walkPos.y < topLeft.y || (walkPos.y == topLeft.y && walkPos.x < topLeft.x) {
			topLeft = walkPos
		}
		// Get the adjacent position in the direction of travel
		adjPos := iteratePosition(walkPos, walkDir)
		// Check if the adjacent position in the direction of travel is within the bounds of the garden
		if adjPos.x < 0 || adjPos.x >= len(garden[0]) || adjPos.y < 0 || adjPos.y >= len(garden) {
			// We've reached the edge of the garden, change direction clockwise
			walkDir = (walkDir + 1) % 4
			sideCount++
		} else if garden[adjPos.y][adjPos.x] == plant {
			// Check if the adjacent position in the direction of travel is the same plant
			// It is but it could be a left hand bend (right hand bends already handled)
			// Need to check the left hand plot from direction of travel to see if it is the same plant
			leftPos := iteratePosition(adjPos, (walkDir+3)%4)
			isLeftTurn := false
			// Check if the left position is within the bounds of the garden and is the same plant
			if leftPos.x >= 0 && leftPos.x < len(garden[0]) && leftPos.y >= 0 && leftPos.y < len(garden) {
				if garden[leftPos.y][leftPos.x] == plant {
					// It is the same plant and in the garden, so this is a left turn
					isLeftTurn = true
				}
			}
			if isLeftTurn {
				// Turn left
				walkDir = (walkDir + 3) % 4
				walkPos = leftPos
				sideCount++
			} else {
				// Continue walking in the current direction
				walkPos = adjPos
			}
		} else {
			// We've reached the edge of the region, change direction clockwise
			walkDir = (walkDir + 1) % 4
			sideCount++
		}
		// Continue walking until we reach the same point and direction we started from
		if walkPos == startPos && walkDir == startDir {
			// We're done walking the perimeter
			break
		}
		perimeterNodes[walkPos] = append(perimeterNodes[walkPos], walkDir)
	}

	// spew.Dump(perimeterNodes)

	// Calculate the total permimeter length
	perimeter := 0
	for _, dirs := range perimeterNodes {
		perimeter += len(dirs)
	}

	// Start creating return region
	newRegion := region{plant: plant,
		start:     topLeft,
		clockdir:  startDir,
		single:    (perimeter == 4),
		perimeter: perimeter,
		sides:     sideCount}

	// Calcluate the area of the region and process inner regions
	newRegion.floodFillRegion(garden, perimeterNodes)

	return true, newRegion
}

// Flood fill a region to find the area and detect inner regions
func (reg *region) floodFillRegion(garden [][]rune, parimeterNodes map[position][]direction) {
	// Create a map of visited positions
	visited := make(map[position]bool, 0)
	// // Create a map of inner regions
	// innerRegions := make(map[position]region, 0)

	// Create a stack for the flood fill
	stack := make([]position, 0)
	stack = append(stack, reg.start)

	// Iterate until the stack is empty
	for len(stack) > 0 {
		// Pop the top position from the stack
		pos := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// Check if the position has already been visited
		if _, ok := visited[pos]; ok {
			// Skip this position
			continue
		}

		// Mark the position as visited
		visited[pos] = true

		// // Check if the position is part of the perimeter
		// if _, ok := parimeterNodes[pos]; ok {
		// 	// Skip this position
		// 	continue
		// }

		// Check if the position is the same plant
		if garden[pos.y][pos.x] != reg.plant {
			// This must be an internal region
			// TODO: Need to process internal region and add to innerRegions
			continue
		}

		// Increment the area of the region
		reg.area++

		// Add the adjacent positions to the stack
		for dir := Up; dir <= Left; dir++ {
			if isBoundaryEdge(garden, parimeterNodes[pos], dir) {
				// Skip this direction
				continue
			}
			adjPos := iteratePosition(pos, dir)
			// Check if the adjacent position is within the bounds of the garden
			if !(adjPos.x < 0 || adjPos.x >= len(garden[0]) || adjPos.y < 0 || adjPos.y >= len(garden)) {
				stack = append(stack, adjPos)
			}
		}
	}
}

// This function is to support flood filling.  It checks whether we should cross the border in the direction
func isBoundaryEdge(garden [][]rune, parimeterNode []direction, dir direction) bool {
	// Funnily we record the direction of the border, but the crossing direction is 90 degrees to the right
	effectiveDir := (dir + 1) % 4
	// Loop over the recorded border directions for this node
	for _, d := range parimeterNode {
		if d == effectiveDir {
			// We found a matching border direction, so we should not cross the border
			return true
		}
	}
	return false
}

// Print the garden
func printGarden(garden [][]rune) {
	for _, row := range garden {
		fmt.Println(string(row))
	}
}

// Print the regions
func printRegions(regions map[position]region) {
	for pos, reg := range regions {
		fmt.Printf("Region: %c, Start: %v, Perimeter: %5d, Area: %5d, Sides: %5d\n", reg.plant, pos, reg.perimeter, reg.area, reg.sides)
	}
}

// Iterate position in a direction
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
		pos.y--
		pos.x--
	case UpRight:
		pos.y--
		pos.x++
	case DownLeft:
		pos.y++
		pos.x--
	case DownRight:
		pos.y++
		pos.x++
	}
	return pos
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
