package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type BuySequence [4]int8
type PriceData struct {
	Price uint8
	Delta int8
}

func main() {
	fmt.Println("Hello, day 22!")

	// Part 1
	// Load the input data
	inputData := LoadInputData("input1.txt")
	// inputData := LoadInputData("test1.txt")
	// inputData := LoadInputData("test2.txt")

	// Parse the input data
	buyerSecrets := make([]uint32, 0)
	for _, line := range inputData {
		if line == "" {
			continue
		}
		secret, _ := strconv.ParseUint(line, 10, 32)
		buyerSecrets = append(buyerSecrets, uint32(secret))
	}

	// {
	// 	// Generate the secrets
	// 	initialValue := 123
	// 	value := initialValue
	// 	fmt.Println("Initial value:", value)
	// 	for i := 0; i < 10; i++ {

	// 		// XOR value and value multiplied by 64, then modulo operation with 16777216, all as a binary operations
	// 		value = (value ^ (value << 6)) & 16777215

	// 		// XOR value and value divided by 32, then modulo operation with 16777216, as a binary operation
	// 		value = (value ^ (value >> 5)) & 16777215

	// 		// XOR value and value divided by 32, then modulo operation with 16777216, as a binary operation
	// 		value = (value ^ (value << 11)) & 16777215

	// 		// Print the result
	// 		fmt.Printf("Step %d: %d\n", i+1, value)
	// 	}
	// 	// Now test the function
	// 	fmt.Printf("Step %d: %d\n", 10, generateSecret(initialValue, 10))
	// }

	// Generate the 2000th iteration of each secret number, and sum the results
	// Execute each function on a separate goroutine
	var part1Sum uint32
	for _, initialSecret := range buyerSecrets {
		secret := generateSecret(initialSecret, 2000)
		part1Sum += secret
		fmt.Printf("%d: %d\n", initialSecret, secret)
	}
	fmt.Println("Part 1: Sum of secrets:", part1Sum)

	// // Attempt at using goroutines
	// var part1GoSum atomic.Uint32
	// var wg sync.WaitGroup
	// for _, initialSecret := range buyerSecrets {
	// 	wg.Add(1)
	// 	go func(initialSecret uint32) {
	// 		defer wg.Done()
	// 		secret := generateSecret(initialSecret, 2000)
	// 		part1GoSum.Add(secret)
	// 		fmt.Printf("%d: %d\n", initialSecret, secret)
	// 	}(initialSecret)
	// }
	// wg.Wait()
	// fmt.Println("Part 1: Sum of secrets:", part1GoSum.Load())

	fmt.Println("----------------------------------------------------")

	// Part 2
	// For a given 4 digit crib, parse the first 2000 iterations of the secret number and find the maximum sum of the last digits
	// Generate the potential buy sequences and add them to a channel

	// Generate the pricing data
	pricingData := make([][]PriceData, 0)
	for _, initialSecret := range buyerSecrets {
		// initialSecret := uint32(123)
		// pricingData = append(pricingData, generatePriceData(initialSecret, 12))
		pricingData = append(pricingData, generatePriceData(initialSecret, 2000))
	}
	// spew.Dump(pricingData)

	// Create map of possible buy sequences by looping over generated price data
	buyOrders := make(map[BuySequence]int)
	for _, priceData := range pricingData {
		for i := 1; i < len(priceData)-3; i++ {
			buyOrders[BuySequence{int8(priceData[i].Delta), int8(priceData[i+1].Delta), int8(priceData[i+2].Delta), int8(priceData[i+3].Delta)}] = 0
		}
	}
	fmt.Printf("Number of buy orders: %d\n", len(buyOrders))

	// Calculate the sum of the buy orders for each buy sequence
	for buyOrder := range buyOrders {
		buyOrders[buyOrder] = calculateBuyOrderSum(pricingData, buyOrder)
	}

	// // Print the first 20 buy orders
	// i := 0
	// for buyOrder, sum := range buyOrders {
	// 	fmt.Printf("BuyOrder %v: %d\n", buyOrder, sum)
	// 	if i >= 20 {
	// 		break
	// 	}
	// 	i++
	// }

	// Find the maximum sum
	maxSum := 0
	maxCrib := BuySequence{}
	for crib, sum := range buyOrders {
		if sum > maxSum {
			// fmt.Printf("New max sum: %d for crib: %v\n", sum, crib)
			maxSum = sum
			maxCrib = crib
		}
	}
	fmt.Println("Part 2: Maximum sum:", maxSum, "for crib:", maxCrib)
}

// Function that generates the a new secret number based defined number of cycles and the initial secret number
func generateSecret(secret uint32, cycles int) uint32 {
	for i := 0; i < cycles; i++ {

		// XOR value and value multiplied by 64, then modulo operation with 16777216, all as a binary operations
		secret = (secret ^ (secret << 6)) & 16777215

		// XOR value and value divided by 32, then modulo operation with 16777216, as a binary operation
		secret = (secret ^ (secret >> 5)) & 16777215

		// XOR value and value divided by 32, then modulo operation with 16777216, as a binary operation
		secret = (secret ^ (secret << 11)) & 16777215
	}
	return secret
}

// Function that generates price data for a given secret number initial value and a given number of cycles
func generatePriceData(secret uint32, cycles int) []PriceData {
	priceData := make([]PriceData, 0, cycles)
	// Add the initial price data
	priceData = append(priceData, PriceData{uint8(secret % 10), 0})
	for i := 0; i < cycles; i++ {

		// XOR value and value multiplied by 64, then modulo operation with 16777216, all as a binary operations
		secret = (secret ^ (secret << 6)) & 16777215

		// XOR value and value divided by 32, then modulo operation with 16777216, as a binary operation
		secret = (secret ^ (secret >> 5)) & 16777215

		// XOR value and value divided by 32, then modulo operation with 16777216, as a binary operation
		secret = (secret ^ (secret << 11)) & 16777215

		// Calculate the price as the last binary digit
		price := uint8(secret % 10)
		// Note: because we added the initial secret, i is the index of that last price
		delta := int8(price) - int8(priceData[i].Price)

		// if i < 20 {
		// 	fmt.Printf("Step %d: Price %d Delta %2d Secret %d\n", i, price, delta, secret)
		// }
		priceData = append(priceData, PriceData{price, delta})
	}
	return priceData
}

// Function that calculates the sum of the prices that match the crib
func calculateBuyOrderSum(priceData [][]PriceData, crib BuySequence) int {
	sum := 0
	for _, priceData := range priceData {
		for i := 1; i < len(priceData)-3; i++ {
			if priceData[i].Delta == crib[0] && priceData[i+1].Delta == crib[1] && priceData[i+2].Delta == crib[2] && priceData[i+3].Delta == crib[3] {
				// if crib[0] == -2 && crib[1] == 1 && crib[2] == -1 && crib[3] == 3 {
				// 	fmt.Printf("Match at %4d: %v %2d %2d %2d %2d Price %d\n", i, crib, priceData[i].Delta, priceData[i+1].Delta, priceData[i+2].Delta, priceData[i+3].Delta, priceData[i+3].Price)
				// }
				// if crib[0] == -2 && crib[1] == 2 && crib[2] == -1 && crib[3] == -1 {
				// 	fmt.Printf("Match at %4d: %v %2d %2d %2d %2d Price %d\n", i, crib, priceData[i].Delta, priceData[i+1].Delta, priceData[i+2].Delta, priceData[i+3].Delta, priceData[i+3].Price)
				// }
				sum += int(priceData[i+3].Price)
				break
			}
		}
	}
	return sum
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
