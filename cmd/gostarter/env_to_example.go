package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func envToExample() {
	// Check if input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Printf("Error: Input file '%s' does not exist\n", inputFile)
		return
	}

	// Generate output filename
	outputFile := inputFile + ".example"
	if strings.HasSuffix(inputFile, ".env") {
		outputFile = strings.TrimSuffix(inputFile, ".env") + ".env.example"
	}

	// Open input file
	input, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		return
	}
	defer input.Close()

	// Create output file
	output, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer output.Close()

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines and comments
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			output.WriteString(line + "\n")
			continue
		}

		// Process variable assignments
		if idx := strings.Index(line, "="); idx != -1 {
			key := strings.TrimSpace(line[:idx])

			// Write the key with a placeholder value for the example file
			output.WriteString(key + "=changeme\n")
		} else {
			// If it's not a key=value pair, just write the line as is
			output.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		return
	}

	fmt.Printf("Successfully created '%s' from '%s'\n", outputFile, inputFile)
}

