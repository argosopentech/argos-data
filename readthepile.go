package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// # thepile.jsonl
//
// {"text": "Heart Team: Joint Position of the Swiss Society of Cardiology and the Swiss Society of Cardiac Surgery.", "meta": {"pile_set_name": "PubMed Abstracts"}}

func main() {
	// Open thepile.jsonl
	filename := "thepile.jsonl"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// Create a new Reader and read the file line by line
	reader := bufio.NewReader(file)
	for {
		// Read line
		line, err := reader.ReadString('\n')
		if err != nil {
			// Skip line if it fails to read
			continue
		}

		// Print line
		fmt.Println(line)
		fmt.Println("\n")

		// Parse JSON
		var result map[string]interface{}
		json.Unmarshal([]byte(line), &result)
		fmt.Println(result["text"])
		fmt.Println("\n")

		// Print meta
		meta := result["meta"].(map[string]interface{})
		fmt.Println(meta["pile_set_name"])

		// Print a blank line
		fmt.Println("\n")

	}
}
