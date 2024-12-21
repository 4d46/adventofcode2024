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

type updateFormat struct {
	list      []int
	positions map[int]int
}

func main() {
	fmt.Println("Hello, day 05!")

	// Part 1
	// Load the input data
	inputData := LoadInputData("input1.txt")
	// inputData := LoadInputData("test1.txt")

	rules := make(map[int][]int)
	updates := make([]updateFormat, 0)

	ruleReg := regexp.MustCompile(`(\d+)\|(\d+)`)
	updateReg := regexp.MustCompile(`,?(\d+)`)

	// Parse input into rules and updates
	for _, line := range inputData {
		if len(line) == 0 {
			continue
		} else if strings.Contains(line, "|") {
			// This is a rule
			matches := ruleReg.FindAllStringSubmatch(line, -1)
			// Should be only one match
			if len(matches) != 1 {
				log.Fatalf("invalid rule: %s", line)
			}
			ruleFirst, _ := strconv.Atoi(matches[0][1])
			ruleSecond, _ := strconv.Atoi(matches[0][2])
			rules[ruleFirst] = append(rules[ruleFirst], ruleSecond)
		} else {
			// This is an update
			matches := updateReg.FindAllStringSubmatch(line, -1)

			var update updateFormat
			update.positions = make(map[int]int)

			// Firstly just record the list of update entries
			for _, match := range matches {
				num, _ := strconv.Atoi(match[1])
				update.list = append(update.list, num)
			}

			// Now record the positions of each number in the list
			for i, num := range update.list {
				update.positions[num] = i
			}

			updates = append(updates, update)
		}
	}
	// spew.Dump(rules)
	// spew.Dump(updates)

	// Now we need to check if the updates match the rules
	// And record the sum of the centre pages for valid updates
	sum := 0
	invalidUpdates := make([]updateFormat, 0)

UPDATE:
	for _, update := range updates {
		// Check if the update matches the rules
		// Loop over the rules
		for ruleFirst, ruleSeconds := range rules {
			// Check the positions of pages the rule defines in the update
			// If the pages are in the correct order, the update is valid
			// First check if we have a page for this rule
			firstPos, ok := update.positions[ruleFirst]
			if !ok {
				// No page for this rule, move on to the next rule
				continue
			}
			// There may be multiple second pages for this rule, check them all
			for _, ruleSecond := range ruleSeconds {
				secondPos, ok := update.positions[ruleSecond]
				if !ok {
					// No page for this rule, move on to the next second rule
					continue
				}
				// Check if the second page is after the first page
				// If it isn't the update is invalid
				if secondPos < firstPos {
					// The second page is before the first page
					// continue to the next update
					// fmt.Printf("Invalid update: %v (breaks rule %d|%d)\n", update.list, ruleFirst, ruleSecond)
					// Add the update to the invalid list for later fixing
					invalidUpdates = append(invalidUpdates, update)
					continue UPDATE
				}

			}
		}
		// If we get here, the update is valid
		// fmt.Printf("Valid update:   %v\n", update.list)
		// Add the centre page to the sum
		centrePos := (len(update.list) - 1) / 2
		sum += update.list[centrePos]
	}

	fmt.Printf("Part 1: Sum of centre pages: %d\n", sum)

	fmt.Println("-----------------------------")

	// Part 2
	// Fix the invalid updates
	fmt.Printf("Number of invalid updates: %d\n", len(invalidUpdates))

	sumOfFixed := 0

	// Loop over the invalid updates
	for _, update := range invalidUpdates {
		// Loop over the rules to find the issue and swap the page positions when invalid rule is found
		// Repeat until update passes all rules
		attempt := 0
	RETRYAA:
		for {
			attempt++

			if attempt > 100 {
				log.Fatalf("Too many attempts to fix update: %v", update.list)
			}

			// Check if the update matches the rules
			// Loop over the rules
			for ruleFirst, ruleSeconds := range rules {
				// Check the positions of pages the rule defines in the update
				// If the pages are in the correct order, the update is valid
				// First check if we have a page for this rule
				firstPos, ok := update.positions[ruleFirst]
				if !ok {
					// No page for this rule, move on to the next rule
					continue
				}
				// There may be multiple second pages for this rule, check them all
				for _, ruleSecond := range ruleSeconds {
					secondPos, ok := update.positions[ruleSecond]
					if !ok {
						// No page for this rule, move on to the next second rule
						continue
					}
					// Check if the second page is after the first page
					// If it isn't the update is invalid
					if secondPos < firstPos {
						// The second page is before the first page, let's swap them and then recheck
						// fmt.Printf("Invalid update: %v (breaks rule %d|%d)\n", update.list, ruleFirst, ruleSecond)
						// Swap the pages
						update.list[firstPos], update.list[secondPos] = update.list[secondPos], update.list[firstPos]
						update.positions[ruleFirst], update.positions[ruleSecond] = update.positions[ruleSecond], update.positions[ruleFirst]
						// Now let's retry the update
						continue RETRYAA
					}

				}
			}
			// If we get here, the update is now valid, calculate the sum and move onto the next invalid update
			// fmt.Printf("Valid update:   %v\n", update.list)
			// Add the centre page to the sum
			centrePos := (len(update.list) - 1) / 2
			sumOfFixed += update.list[centrePos]
			break
		}
	}

	fmt.Printf("Part 2: Sum of previously invalid centre pages: %d\n", sumOfFixed)

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
