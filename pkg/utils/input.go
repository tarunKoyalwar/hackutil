package utils

import (
	"bufio"
	"os"
	"strings"
)

// GetInputList processes input that can be either a file path or comma-separated values
func GetInputList(input string) ([]string, error) {
	// Check if input is empty
	if input == "" {
		return []string{}, nil
	}

	// Check if input is a file path
	if _, err := os.Stat(input); err == nil {
		return readFileLines(input)
	}

	// Process as comma-separated values
	return strings.Split(input, ","), nil
}

// readFileLines reads a file line by line and returns a slice of strings
func readFileLines(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
} 