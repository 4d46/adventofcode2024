package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Equation struct {
	result    int
	operands  []int
	operators []int
}

func main() {
	fmt.Println("Hello, day 07!")

	// Part 1
	// Load the input data
	inputData := LoadInputData("input1.txt")
	// inputData := LoadInputData("test1.txt")

	// Parse the input data into a equations
	regexpEquation := regexp.MustCompile(`(\d+):([ \d]+)`)
	var equations []Equation
	for _, line := range inputData {
		var equation Equation
		match := regexpEquation.FindStringSubmatch(line)
		if match == nil {
			log.Fatalf("failed to parse line: %s", line)
		}
		result, err := strconv.Atoi(match[1])
		if err != nil {
			log.Fatalf("failed to parse result: %s", match[1])
		}
		equation.result = result

		operands := strings.Split(strings.Trim(match[2], " "), " ")
		for _, operand := range operands {
			operandInt, err := strconv.Atoi(operand)
			if err != nil {
				log.Fatalf("failed to parse operand: %s", operand)
			}
			equation.operands = append(equation.operands, operandInt)
		}

		equations = append(equations, equation)
	}

	// spew.Dump(equations)

	// Loop through the equations and attempt to solve them
	for i := 0; i < len(equations); i++ {
		// for i := 2; i < 3; i++ {
		// Use binary to represent all operator combinations, calculate the maximum value
		maxOperatorCombo := 1 << (len(equations[i].operands) - 1)
		// formatString := "%" + fmt.Sprintf("0%db\n", len(equations[i].operands)-1)
		for opcombo := 0; opcombo < maxOperatorCombo; opcombo++ {
			// fmt.Printf(formatString, i)
			potentialAnswer := calculateEquation(equations[i], opcombo)
			if potentialAnswer == equations[i].result {
				// fmt.Printf("Found an answer: %d\n", potentialAnswer)
				equations[i].operators = append(equations[i].operators, opcombo)
				// spew.Dump(equations[i])
			}
		}
		// fmt.Println()
	}

	// Print the equations and calculate the valid equations sum
	validEquationsSum := 0
	for _, equation := range equations {
		if len(equation.operators) > 0 {
			equation.Print()
			validEquationsSum += equation.result
		}
	}

	fmt.Printf("Part 1: The sum of the valid equations is %d\n", validEquationsSum)

	fmt.Println("----------------------------------------")

	// Part 2
	// opTest1 := 12 & 0b11
	// fmt.Printf("opTest1: %b %d\n", opTest1, opTest1)
	// opTest2 := (12 >> (2 * 1)) & 0b11
	// fmt.Printf("opTest2: %b %d\n", opTest2, opTest2)

	// num1 := 1234
	// num2 := 5678
	// numInt := int(math.Pow10(int(math.Log10(float64(num2))) + 1))
	// fmt.Printf("numInt: %d\n", numInt)
	// num3 := (num1 * int(math.Pow10(int(math.Log10(float64(num2)))+1))) + num2
	// fmt.Printf("num3: %d\n", num3)

	// testEquation := Equation{
	// 	result:    71,
	// 	operands:  []int{1, 2, 3, 4, 5},
	// 	operators: []int{144},
	// }
	// testEquation.PrintP2()
	// for i := 0; i < 16; i++ {
	// 	operTest1 := validOperator(i, 2)
	// 	fmt.Printf("operTest1: %d %b %t\n", i, i, operTest1)
	// }
	// os.Exit(0)
	// Loop through the equations and attempt to solve them
	for i := 0; i < len(equations); i++ {
		// Clear previous operators
		equations[i].operators = []int{}

		numOperators := len(equations[i].operands) - 1
		// Use 2 bit binary to represent all operator combinations, calculate the maximum value
		maxOperatorCombo := 1 << (numOperators * 2)

		for opcombo := 0; opcombo < maxOperatorCombo; opcombo++ {
			// Check if this is a valid operator combination
			if !validOperator(opcombo, numOperators) {
				continue
			}
			potentialAnswer := calculateEquationP2(equations[i], opcombo)
			if potentialAnswer == equations[i].result {
				// fmt.Printf("Found an answer: %d\n", potentialAnswer)
				equations[i].operators = append(equations[i].operators, opcombo)
				// spew.Dump(equations[i])
			}
		}
		// fmt.Println()
	}

	// Print the equations and calculate the valid equations sum
	validEquationsSumP2 := 0
	for _, equation := range equations {
		if len(equation.operators) > 0 {
			equation.PrintP2()
			validEquationsSumP2 += equation.result
		}
	}

	fmt.Printf("Part 2: The sum of the valid equations is %d\n", validEquationsSumP2)

}

func calculateEquation(equation Equation, operatorCombo int) int {
	result := equation.operands[0]
	for i := 1; i < len(equation.operands); i++ {
		// Check if the operator is a 1 or a 0
		operator := operatorCombo & (1 << (i - 1))
		if operator != 0 {
			result += equation.operands[i]
		} else {
			result *= equation.operands[i]
		}
	}
	return result
}

func calculateEquationP2(equation Equation, operatorCombo int) int {
	result := equation.operands[0]
	// fmt.Printf("operatorCombo: %b %d\n", operatorCombo, operatorCombo)
	for i := 1; i < len(equation.operands); i++ {
		// Check if the operator is a 1 or a 0
		operator := (operatorCombo >> (2 * (i - 1))) & 0b11
		// fmt.Printf("operator: %b %d\n", operator, operator)
		switch operator {
		case 0: // +
			result += equation.operands[i]
		case 1: // *
			result *= equation.operands[i]
		case 2: // ||
			// Add the next operand to the result, but like 2 strings were being concatenated
			result = (result * int(math.Pow10(int(math.Log10(float64(equation.operands[i])))+1))) + equation.operands[i]
		case 3: // Invalid
			log.Fatalf("invalid operator: %d", operator)
		}
	}
	return result
}

// Print the equation
func (e Equation) Print() {
	for _, answer := range e.operators {
		fmt.Printf("%d = %d", e.result, e.operands[0])
		for i := 1; i < len(e.operands); i++ {
			// fmt.Printf(" %d %b %b\n", answer, (1 << (i - 1)), (answer & (1 << (i - 1))))
			if (answer & (1 << (i - 1))) != 0 {
				fmt.Printf(" + %d", e.operands[i])
			} else {
				fmt.Printf(" * %d", e.operands[i])
			}
		}
		fmt.Println()
	}
}

func (e Equation) PrintP2() {
	for _, answer := range e.operators {
		fmt.Printf("%d = %d", e.result, e.operands[0])
		for i := 1; i < len(e.operands); i++ {
			operator := (answer >> (2 * (i - 1))) & 0b11
			// fmt.Printf("operator: %b %d\n", operator, operator)
			switch operator {
			case 0: // +
				fmt.Printf(" + %d", e.operands[i])
			case 1: // *
				fmt.Printf(" * %d", e.operands[i])
			case 2: // ||
				fmt.Printf(" || %d", e.operands[i])
			case 3: // Invalid
				log.Fatalf("invalid operator: %d", operator)
			}
		}
		fmt.Println()
	}
}

// Check if the operators are valid
func validOperator(operator, length int) bool {
	// Check each pair of bits in the operator does not have both bits set
	for i := 0; i < length; i++ {
		if ((operator >> (2 * i)) & 0b11) == 0b11 {
			return false
		}
	}
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
