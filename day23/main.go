package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

func main() {
	fmt.Println("Hello, day 22!")

	// Part 1
	// Load the input data
	inputData := LoadInputData("input1.txt")
	// inputData := LoadInputData("test1.txt")

	// The strings in the input data are a set of instructions
	// Each instruction is in the format "mul(1,20)", but there are lots of errors in the data too
	// We need to use a regular expression to extract the instruction strings
	regex := `(.+)-(.+)`
	m := regexp.MustCompile(regex)

	allConnections := make(map[string][]string, 0)
	tComputers := make(map[string]bool, 0)

	// Loop through the input data and find all connections, adding them to a map
	for _, line := range inputData {
		// Find the connection string by running the regular expression on the line
		connection := m.FindAllStringSubmatch(line, -1)
		// Add the connection in both directions to the map
		allConnections[connection[0][1]] = append(allConnections[connection[0][1]], connection[0][2])
		allConnections[connection[0][2]] = append(allConnections[connection[0][2]], connection[0][1])

		// Also find all the computers with a t in their name
		if connection[0][1][0] == 't' {
			tComputers[connection[0][1]] = true
		}
		if connection[0][2][0] == 't' {
			tComputers[connection[0][2]] = true
		}
	}
	// spew.Dump(allConnections)

	{
		// Starting from computers with a t in their name, find all the computers that form a loop of 3
		// Loop through the tComputers map
		// Some loops will be clockwise and some anti-clockwise
		// But also some loops will contains 2 entries of t so we will have duplicate loops, dedupe them
		foundLoops := make(map[string]bool, 0)
		for tComputer := range tComputers {
			// Find the connections for this computer
			connections := allConnections[tComputer]
			// Loop through the connections
			for _, connection := range connections {
				// Find the connections for this connection
				connections2 := allConnections[connection]
				// Loop through the connections
				for _, connection2 := range connections2 {
					// Find the connections for this connection
					connections3 := allConnections[connection2]
					// Loop through the connections
					for _, connection3 := range connections3 {
						// If the connection3 is the same as the tComputer, we have a loop
						if connection3 == tComputer {
							// fmt.Printf("Found loop: %s -> %s -> %s\n", tComputer, connection, connection2)
							// Add the loop to the foundLoops map
							// Create a key which is the nodes sorted in alphabetical order
							// This will allow us to dedupe the loops
							keyParts := []string{tComputer, connection, connection2}
							// Sort the keyParts
							sort.Strings(keyParts)
							key := strings.Join(keyParts, "-")
							foundLoops[key] = true
						}
					}
				}
			}
		}
		fmt.Printf("Part 1: Number of loops: %d\n", len(foundLoops))

	}

	fmt.Println("-----------------------------")

	// Part 2
	{
		// Loop over all computer names and sort connected computers
		for startNode := range allConnections {
			connections := allConnections[startNode]
			sort.Strings(connections)
			allConnections[startNode] = connections
			// fmt.Printf("%s: %v\n", startNode, connections)
		}

		// Remember which computers we have found the max set size for, so we don't have to do it again
		nodeMaxSetSize := make(map[string]int, 0)

		largestSetSize := 0
		largestSetKey := ""

		// Loop over all computer names
		for startNode := range allConnections {
			// If we have already found the max set size for this node, skip it
			if _, ok := nodeMaxSetSize[startNode]; ok {
				continue
			}

			// Create a map to hold the set of connected computers
			computerSets := make(map[string]bool, 0)
			// Add the associated computers to the set, including the startNode
			computerSets[startNode] = true
			// Loop over the connections for the startNode
			for _, connection := range allConnections[startNode] {
				// Add the connection to the set
				computerSets[connection] = true
			}

			intersectedNodes := make(map[string]bool, 0)
			// Loop until we have found all the computers in the set
			for {
				// Find the next node to perform an interesection on
				nextNode := ""
				for node := range computerSets {
					if _, ok := intersectedNodes[node]; !ok {
						nextNode = node
						break
					}
				}
				// If we have intersected all the nodes, break out of the loop
				if nextNode == "" {
					break
				}
				// Add the node to the intersectedNodes map
				intersectedNodes[nextNode] = true
				// Loop over the current connection set and remove any that are not in the nextNode connection set
				for node := range computerSets {
					// Don't bother checking the startNode or the nextNode
					if node == startNode || node == nextNode {
						continue
					}
					// If the node is not in the nextNode connection set, remove it from the computerSets map
					inSet := false
					for _, connection := range allConnections[nextNode] {
						if node == connection {
							inSet = true
							break
						}
					}
					if !inSet {
						delete(computerSets, node)
					}
				}
			}

			// This should be a set of interlinked computers, remember the size of the set for all the nodes in the set
			// Hmm, not sure if this is always the case or we were just lucky on the ordering of the nodes
			// After running doesn't really seem to speed things up so ignoring
			// for node := range computerSets {
			// 	nodeMaxSetSize[node] = len(computerSets)
			// }

			// Remember the largest set size and the key
			if len(computerSets) > largestSetSize {
				largestSetSize = len(computerSets)
				setNodes := make([]string, 0)
				for node := range computerSets {
					setNodes = append(setNodes, node)
				}
				sort.Strings(setNodes)
				largestSetKey = strings.Join(setNodes, ",")
			}
		}

		fmt.Printf("Part 2: Largest set: %3d %s\n", largestSetSize, largestSetKey)
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
