package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Hello, day 02!")

	// Part 1
	// Load the input data
	inputData := LoadInputData("input1.txt")
	// inputData := LoadInputData("test1.txt")

	// The strings  in the input data are Reports, with a list of space separated integers
	// Loop through each report and determine is they are safe, counting the number of safe reports
	var safeReports int
	for _, report := range inputData {
		if IsSafeReport(report) {
			safeReports++
		}
	}

	fmt.Printf("Part 1: There are %d safe reports\n", safeReports)

	fmt.Println("--------------------")

	// Part 2
	// Loop through each report and determine is they are safe, counting the number of safe reports with the damped check
	var safeReports2 int
	for _, report := range inputData {
		if IsSafeReport(report) {
			safeReports2++
		} else {
			// Report isn't as safe as we thought, so we need to check the damped versions
			// Generate all possible damped reports
			// Note: Tried being clever about this first, but my approach the corner cases were becoming a pain
			//       So given the data set is small, I brute forced it instead
			dampedReports := generateDampedReports(report)
			// Loop through each damped report and determine any are safe
			for _, dampedReport := range dampedReports {
				// If we find a safe report, increment the count and break out of the loop
				if IsSafeReport(dampedReport) {
					safeReports2++
					break
				}
			}
		}
	}

	fmt.Printf("Part 2: There are %d safe reports\n", safeReports2)
}

// generateDampedReports takes a report string and returns a slice of damped reports
func generateDampedReports(report string) []string {
	var dampedReports []string
	// Create a number of damped reports by removing a level from the reports
	// First, split the report into a slice of integers (levels) strings
	levels := strings.Split(report, " ")
	// This should produce as many reports as there are levels in the original report
	for i := 0; i < len(levels); i++ {
		// Create a damped report by removing a level from the report
		var delimiter string
		if i == 0 || i == len(levels)-1 {
			delimiter = ""
		} else {
			delimiter = " "
		}
		dampedReport := strings.Join(levels[:i], " ") + delimiter + strings.Join(levels[i+1:], " ")
		// spew.Dump(dampedReport)
		// Append the damped report to the dampedReports slice
		dampedReports = append(dampedReports, dampedReport)
	}
	return dampedReports
}

// IsSafeReport takes a report string and returns true if the report is safe
func IsSafeReport(report string) bool {
	// spew.Dump(report)
	// Split the report into a slice of integers (levels)
	levels := strings.Split(report, " ")

	// Loop through each level, starting at the second level
	var previousLevel, levelDirection int
	for i, x := range levels {
		if i == 0 {
			// Convert the first level to an integer
			previousLevel, _ = strconv.Atoi(x)
			continue
		}
		// Convert the level to an integer
		level, _ := strconv.Atoi(x)
		// Check if the level is safe by being between 1 or 2 above or below previous level
		difference := Difference(previousLevel, level)
		if difference < 1 || difference > 3 {
			return false
		}
		// Check if the level is consistantly increasing or decreasing
		if level-previousLevel > 0 {
			levelDirection += 1
		} else {
			levelDirection -= 1
		}
		// fmt.Printf("Previous Level: %d, Level: %d, Difference: %d, Level Direction: %d\n", previousLevel, level, difference, levelDirection)
		if abs(levelDirection) != i {
			return false
		}
		// Set the previous level to the current level
		previousLevel = level
	}
	// fmt.Println("Report is safe")
	// If we get here, the report is safe
	return true
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

// Difference returns the difference between 2 integers
func Difference(x, y int) int {
	if x == y {
		return 0
	} else if x > y {
		return x - y
	} else {
		return y - x
	}
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
