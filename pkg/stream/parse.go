package stream

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Choice struct {
	Delta struct {
		Content string `json:"content"`
	} `json:"delta"`
}

type Data struct {
	Choices []Choice `json:"choices"`
}

func ParseFile(filename string) error {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var contentBuilder strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line has "data: " prefix
		if strings.HasPrefix(line, "data: ") {
			// Remove the "data: " prefix
			line = strings.TrimPrefix(line, "data: ")
		} else {
			continue // skip lines without "data: "
		}

		// Parse the JSON line into our `Data` struct
		var data Data
		err := json.Unmarshal([]byte(line), &data)
		if err != nil {
			// Skip this line if JSON is incomplete or malformed
			continue
		}

		// Extract delta.content and concatenate it
		for _, choice := range data.Choices {
			contentBuilder.WriteString(choice.Delta.Content)
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Print the final concatenated result
	result := contentBuilder.String()
	fmt.Println(result)

	return nil
}
