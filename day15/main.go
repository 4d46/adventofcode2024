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

type Direction int

// Word directions
const (
	Up Direction = iota
	Right
	Down
	Left
	Undefined
)

const (
	Wall      = '#'
	Box       = 'O'
	Robot     = '@'
	Space     = '.'
	DoubleBox = "[]"
	BoxL      = '['
	BoxR      = ']'
)

func main() {
	fmt.Println("Hello, day 15!")

	// Part 1
	// Load the input data
	// inputData := LoadInputData("input1.txt")
	// inputData := LoadInputData("test1.txt")
	// inputData := LoadInputData("test2.txt")
	inputData := LoadInputData("test3.txt")
	{
		// The strings in the input data are a warehouse layout
		// Convert the strings to a 2D slice of runes
		var warehouse [][]rune
		var moveInstructions []Direction
		for _, line := range inputData {
			if len(line) < 1 {
				continue
			} else if line[0] == Wall {
				warehouse = append(warehouse, []rune(line))
			} else {
				// This must be a set of move instructions
				// Loop over the string and convert to directions
				for _, r := range line {
					switch r {
					case '^':
						moveInstructions = append(moveInstructions, Up)
					case '>':
						moveInstructions = append(moveInstructions, Right)
					case 'v':
						moveInstructions = append(moveInstructions, Down)
					case '<':
						moveInstructions = append(moveInstructions, Left)
					default:
						log.Fatalf("Invalid move instruction: %c", r)
					}
				}
			}
		}
		printWarehouse(warehouse)

		// Find the starting position of the robot
		var robotPos Position
		for y, row := range warehouse {
			for x, r := range row {
				if r == Robot {
					robotPos = Position{x, y}
					break
				}
			}
		}
		fmt.Printf("Robot starting position: %v\n", robotPos)

		// Move the robot in the specified directions
		for _, dir := range moveInstructions {
			robotPos = moveRobot(warehouse, robotPos, dir)
			// fmt.Printf("--------------------\n")
			// printWarehouse(warehouse)
		}
		printWarehouse(warehouse)
		fmt.Printf("Part 1: %d\n", sumBoxLocations(warehouse))
	}
	fmt.Println("--------------------")
	{
		// Part 2
		// Load the input data, but this time the map should be created with walls, boxes and space twice as wide as the original map
		// inputData := LoadInputData("input2.txt")
		var warehouse [][]rune
		var moveInstructions []Direction
		for _, line := range inputData {
			if len(line) < 1 {
				continue
			} else if line[0] == Wall {
				doubleWidth := makeMapDoubleWidth(line)
				warehouse = append(warehouse, doubleWidth)
			} else {
				// This must be a set of move instructions
				// Loop over the string and convert to directions
				for _, r := range line {
					switch r {
					case '^':
						moveInstructions = append(moveInstructions, Up)
					case '>':
						moveInstructions = append(moveInstructions, Right)
					case 'v':
						moveInstructions = append(moveInstructions, Down)
					case '<':
						moveInstructions = append(moveInstructions, Left)
					default:
						log.Fatalf("Invalid move instruction: %c", r)
					}
				}
			}
		}
		printWarehouse(warehouse)

		// Find the starting position of the robot
		var robotPos Position
		for y, row := range warehouse {
			for x, r := range row {
				if r == Robot {
					robotPos = Position{x, y}
					break
				}
			}
		}
		fmt.Printf("Robot starting position: %v\n", robotPos)

		// Move the robot in the specified directions, but this time boxes are twice as wide and consist of 2 characters "[]"
		for _, dir := range moveInstructions {
			robotPos = moveRobotWide(warehouse, robotPos, dir)
			fmt.Printf("--------------------\n")
			printWarehouse(warehouse)
		}
	}
}

// Create a new map that is twice as wide as the original map
func makeMapDoubleWidth(line string) []rune {
	doubleWidth := make([]rune, 0)
	for _, r := range line {
		switch r {
		case Wall:
			doubleWidth = append(doubleWidth, Wall, Wall)
		case Box:
			doubleWidth = append(doubleWidth, '[', ']')
		case Robot:
			doubleWidth = append(doubleWidth, Robot, '.')
		case Space:
			doubleWidth = append(doubleWidth, Space, Space)
		default:
			log.Fatalf("Invalid character in map: %c", r)
		}
	}
	return doubleWidth
}

// Attempt to move the robot in the specificed direction in the warehouse.  If the robot impacts a Box then if there is an available space
// in the direction of movement, the Box will be moved in that direction.  Do this recursively until the robot is able to move in the
// specified direction. If the robot is unable to move in the specified direction, return that same position, otherwise return the new position
// But boxes ar now twice as wide and consist of 2 characters "[]", Boxes are moved as a single unit and boxes that overlap at all will be moved
// as a single unit
func moveRobotWide(warehouse [][]rune, robotPos Position, dir Direction) Position {
	// Get the new position
	newPos := iteratePosition(robotPos, dir)

	// Check if the new position is a wall
	if warehouse[newPos.y][newPos.x] == Wall {
		return robotPos
	}

	// Check if the new position is a Box
	if warehouse[newPos.y][newPos.x] == BoxL || warehouse[newPos.y][newPos.x] == BoxR {
		// Attempt to move the Box
		boxPos := moveBoxWide(warehouse, newPos, dir)
		// If the Box was unable to move, return the current position
		if boxPos == newPos {
			return robotPos
		}
	}

	// Move the robot to the new position
	warehouse[robotPos.y][robotPos.x] = '.'
	warehouse[newPos.y][newPos.x] = Robot
	return newPos
}

// Move the Box in the specified direction.  If the Box is unable to move in the specified direction, return the current position, otherwise
// return the new position
// But boxes ar now twice as wide and consist of 2 characters "[]", Boxes are moved as a single unit and boxes that overlap at all will be moved
// as a single unit
func moveBoxWide(warehouse [][]rune, boxPos Position, dir Direction) Position {
	// Get the 2 new positions the box will occupy
	newPosL := iteratePosition(boxPos, dir)
	newPosR := iteratePosition(newPosL, dir)

	// Check if either new position is a wall
	if warehouse[newPosL.y][newPosL.x] == Wall || warehouse[newPosR.y][newPosR.x] == Wall {
		return boxPos
	}

	// Check if either new position is a Box
	if warehouse[newPosL.y][newPosL.x] == BoxL || warehouse[newPosL.y][newPosL.x] == BoxR ||
		warehouse[newPosR.y][newPosR.x] == BoxL || warehouse[newPosR.y][newPosR.x] == BoxR {
		// Attempt to move the Box
		newBoxPos := moveBoxWide(warehouse, newPosL, dir)
		// If the Box was unable to move, return the current position
		if newBoxPos == newPosL {
			return boxPos
		}
	}

	// Move the Box to the new position
	warehouse[boxPos.y][boxPos.x] = '.'
	warehouse[newPosL.y][newPosL.x] = BoxL
	warehouse[newPosR.y][newPosR.x] = BoxR
	return newPosL

}

// Attempt to move the robot in the specificed direction in the warehouse.  If the robot impacts a Box then if there is an available space
// in the direction of movement, the Box will be moved in that direction.  Do this recursively until the robot is able to move in the
// specified direction. If the robot is unable to move in the specified direction, return that same position, otherwise return the new position
func moveRobot(warehouse [][]rune, robotPos Position, dir Direction) Position {
	// Get the new position
	newPos := iteratePosition(robotPos, dir)

	// Check if the new position is a wall
	if warehouse[newPos.y][newPos.x] == Wall {
		return robotPos
	}

	// Check if the new position is a Box
	if warehouse[newPos.y][newPos.x] == Box {
		// Attempt to move the Box
		boxPos := moveBox(warehouse, newPos, dir)
		// If the Box was unable to move, return the current position
		if boxPos == newPos {
			return robotPos
		}
	}

	// Move the robot to the new position
	warehouse[robotPos.y][robotPos.x] = '.'
	warehouse[newPos.y][newPos.x] = Robot
	return newPos
}

// Move the Box in the specified direction.  If the Box is unable to move in the specified direction, return the current position, otherwise
// return the new position
func moveBox(warehouse [][]rune, boxPos Position, dir Direction) Position {
	// Get the new position
	newPos := iteratePosition(boxPos, dir)

	// Check if the new position is a wall
	if warehouse[newPos.y][newPos.x] == Wall {
		return boxPos
	}

	// Check if the new position is a Box
	if warehouse[newPos.y][newPos.x] == Box {
		// Attempt to move the Box
		newBoxPos := moveBox(warehouse, newPos, dir)
		// If the Box was unable to move, return the current position
		if newBoxPos == newPos {
			return boxPos
		}
	}

	// Move the Box to the new position
	warehouse[boxPos.y][boxPos.x] = '.'
	warehouse[newPos.y][newPos.x] = Box
	return newPos
}

// Create a function that sums all the box locations
// Each box location is calculated as 100*y + x
// Ignore the robot location and walls
// Return the total sum of all box locations
func sumBoxLocations(warehouse [][]rune) int {
	sum := 0
	for y, row := range warehouse {
		for x, r := range row {
			if r == Box {
				sum += 100*y + x
			}
		}
	}
	return sum
}

// Print the warehouse
func printWarehouse(warehouse [][]rune) {
	for _, row := range warehouse {
		fmt.Println(string(row))
	}
}

// Iterate position in a direction
func iteratePosition(pos Position, dir Direction) Position {
	switch dir {
	case Up:
		pos.y--
	case Down:
		pos.y++
	case Left:
		pos.x--
	case Right:
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
