package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Position struct {
	x int
	y int
}

type Velocity struct {
	x int
	y int
}

type RobotState struct {
	p Position
	v Velocity
}

const (
	quadTL = iota
	quadTR
	quadBL
	quadBR
)

func main() {
	fmt.Println("Hello, day 14!")

	// Part 1
	// Load the input data
	inputData := LoadInputData("input1.txt")
	// inputData := LoadInputData("test1.txt")

	// Each line is the position and velocity of a robot
	// Each instruction is in the format:
	// p=0,4 v=3,-3

	// We need to use a regular expression to extract the instruction strings
	regexRobotState := `p=(-?\d+),(-?\d+) v=(-?\d+),(-?\d+)`
	tileArea := `width=(-?\d+),height=(-?\d+)`
	m1 := regexp.MustCompile(regexRobotState)
	m2 := regexp.MustCompile(tileArea)

	initialState := make([]RobotState, 0)
	tileWidth := 0
	tileHeight := 0

	// Loop through the input data and find all instances of Robot state
	for i, line := range inputData {
		// Read the width and height of the tiles from the first line
		if i == 0 {
			tile := m2.FindAllStringSubmatch(line, -1)
			tileWidth, _ = strconv.Atoi(tile[0][1])
			tileHeight, _ = strconv.Atoi(tile[0][2])
			fmt.Printf("width=%d, height=%d\n", tileWidth, tileHeight)
		} else {
			// Find the state strings by running the regular expression on the line
			state := m1.FindAllStringSubmatch(line, -1)

			// Convert the strings to integers
			px, _ := strconv.Atoi(state[0][1])
			py, _ := strconv.Atoi(state[0][2])
			vx, _ := strconv.Atoi(state[0][3])
			vy, _ := strconv.Atoi(state[0][4])
			// Append the state strings to the initialState slice
			initialState = append(initialState, RobotState{Position{px, py}, Velocity{vx, vy}})
		}
	}
	{
		// printState(initialState)
		// fmt.Println("Initial state:")
		// printStateMap(initialState, tileWidth, tileHeight)
		timeJump := 100
		// Project the state of the robots after a certain time
		newState := projectState(initialState, tileWidth, tileHeight, timeJump)
		fmt.Printf("Part 1: After %d seconds\n", timeJump)
		// printStateMap(newState, tileWidth, tileHeight)
		safetyFactor := calculateSafetyFactor(newState, tileWidth, tileHeight)
		fmt.Printf("Part 1: Safety factor=%d\n", safetyFactor)
	}

	fmt.Println("--------------------")

	// Part 2
	// We need to find the time when the Christmas tree is shown
	// We assume the safety factor will be the highest when the Christmas tree is shown
	// We need to find the time when the safety factor is the highest
	// sfMax := 0
	// sfMaxTime := 0
	newTime := 0
	timeStep := 1
	midY := tileHeight / 2
	midX := tileWidth / 2
	for step := 0; step < 10000; step++ {
		newTime += timeStep
		newState := projectState(initialState, tileWidth, tileHeight, newTime)
		// safetyFactor := calculateSafetyFactor(newState, tileWidth, tileHeight)
		// fmt.Printf("Time=%d, Safety Factor=%d\n", newTime, safetyFactor)
		// if safetyFactor > sfMax {
		// 	sfMax = safetyFactor
		// 	sfMaxTime = newTime
		// }
		// Look for a line of robots in the middle of the middle row
		// If we find a line of robots, we can assume that the Christmas tree is shown
		// We can stop the loop
		// Note: tried a few things barking up the wrong tree, looked for some hints on finding bur the tree was too small so didn't find anything
		//       followed sugfgestion to write out to text file and search with text editor that worked.  But then went back and tweaked the algorithm
		//       so it just finds it
		count := 0
		for _, s := range newState {
			// if s.p.y == 42 {
			if s.p.y == midY {
				if s.p.x >= midX-4 && s.p.x <= midX+4 {
					count++
				}
			}
		}
		if count >= 8 {
			fmt.Printf("**FOUND** Time=%d, count=%d\n", newTime, count)
			break
		}
		// fmt.Printf("Time=%d\n", newTime)
		// printStateMap(newState, tileWidth, tileHeight)

	}
	//TEMP
	// sfMaxTime = 104
	maxState := projectState(initialState, tileWidth, tileHeight, newTime)
	maxSafetyFactor := calculateSafetyFactor(maxState, tileWidth, tileHeight)
	fmt.Printf("Part 2: After %d seconds, Max Safety Factor %d\n", newTime, maxSafetyFactor)
	printStateMap(maxState, tileWidth, tileHeight)
}

func printState(state []RobotState) {
	for _, s := range state {
		fmt.Printf("p=(%2d, %2d) v=(%2d, %2d)\n", s.p.x, s.p.y, s.v.x, s.v.y)
	}
}

func printStateMap(state []RobotState, tileWidth, tileHeight int) {
	// Create a map of the state
	mapCounts := make(map[Position]int)
	for _, s := range state {
		mapCounts[s.p]++
	}
	// spew.Dump(mapCounts)
	// Create a state map
	stateMapRune := make([][]rune, 0)
	for y := 0; y < tileHeight; y++ {
		runeSlice := []rune(strings.Repeat(".", tileWidth))
		stateMapRune = append(stateMapRune, runeSlice)
	}
	// Fill in the state map
	for p, count := range mapCounts {
		if count > 0 && count < 10 {
			stateMapRune[p.y][p.x] = rune('0' + count)
		} else if count >= 10 {
			stateMapRune[p.y][p.x] = rune('A' + count - 10)
		} else {
			stateMapRune[p.y][p.x] = '#'
		}
	}

	// Print the state map
	for y := 0; y < tileHeight; y++ {
		fmt.Println(string(stateMapRune[y]))
	}
}

// Project the state of the robots after a certain time, wrap around the edges
func projectState(state []RobotState, tileWidth, tileHeight, time int) []RobotState {
	newState := make([]RobotState, len(state))
	for i, s := range state {
		// Calculate the new position
		newX := (s.p.x + s.v.x*time) % tileWidth
		newY := (s.p.y + s.v.y*time) % tileHeight
		if newX < 0 {
			newX += tileWidth
		}
		if newY < 0 {
			newY += tileHeight
		}
		newState[i] = RobotState{Position{newX, newY}, s.v}
	}
	return newState
}

func calculateSafetyFactor(state []RobotState, tileWidth, tileHeight int) int {
	quadrantCount := [4]int{0, 0, 0, 0}
	midX := tileWidth / 2
	midY := tileHeight / 2

	for _, s := range state {
		if s.p.x < midX {
			if s.p.y < midY {
				quadrantCount[quadTL]++
			} else if s.p.y > midY {
				quadrantCount[quadBL]++
			}
		} else if s.p.x > midX {
			if s.p.y < midY {
				quadrantCount[quadTR]++
			} else if s.p.y > midY {
				quadrantCount[quadBR]++
			}
		}
	}

	// fmt.Printf("Quadrant counts: TL=%d, TR=%d, BL=%d, BR=%d\n", quadrantCount[quadTL], quadrantCount[quadTR], quadrantCount[quadBL], quadrantCount[quadBR])
	// Calcuate the safety factor
	safetyFactor := 1
	for _, count := range quadrantCount {
		if count > 0 {
			safetyFactor *= count
		}
	}
	return safetyFactor
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
