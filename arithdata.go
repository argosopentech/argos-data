package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

func main() {
	fmt.Println("Welcome to Arithmetic Data Generator")
	fmt.Println("-----------------------------------")
	fmt.Println("Please enter the number of arithmetic expressions you want to generate:")

	var numExpressions int
	_, err := fmt.Scanln(&numExpressions)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// Configuration for maximum integer size and supported operators
	maxIntSize := 99
	// operators := []string{"+", "-", "*", "/"}
	operators := []string{"+"}

	fmt.Println("Generating arithmetic expressions...")

	// Create a folder named 'arith' to store the generated data
	err = os.Mkdir("arith", 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Println("Error creating directory:", err)
		return
	}

	sourceFile, err := os.Create("arith/source.txt")
	if err != nil {
		fmt.Println("Error creating source file:", err)
		return
	}
	defer sourceFile.Close()

	targetFile, err := os.Create("arith/target.txt")
	if err != nil {
		fmt.Println("Error creating target file:", err)
		return
	}
	defer targetFile.Close()

	// Write arithmetic expressions to source.txt and corresponding results to target.txt
	for i := 0; i < numExpressions; i++ {
		expression, result := generateArithmeticExpression(maxIntSize, operators)
		sourceFile.WriteString(expression + "\n")
		targetFile.WriteString(result + "\n")
	}

	fmt.Println("Arithmetic data generated successfully in 'arith' folder.")
}

func generateArithmeticExpression(maxIntSize int, operators []string) (string, string) {
	// Generate two random numbers between 1 and maxIntSize
	num1 := rand.Intn(maxIntSize) + 1
	num2 := rand.Intn(maxIntSize) + 1

	// Choose a random arithmetic operator
	operator := operators[rand.Intn(len(operators))]

	// Create the expression and calculate the result
	expression := strconv.Itoa(num1) + operator + strconv.Itoa(num2)
	result := calculateResult(num1, num2, operator)

	return expression, result
}

func calculateResult(num1, num2 int, operator string) string {
	var result int
	switch operator {
	case "+":
		result = num1 + num2
	case "-":
		result = num1 - num2
	case "*":
		result = num1 * num2
	case "/":
		if num2 != 0 {
			result = num1 / num2
		} else {
			// Avoid division by zero
			return calculateResult(num1, num2+1, operator)
		}
	}

	return strconv.Itoa(result)
}

