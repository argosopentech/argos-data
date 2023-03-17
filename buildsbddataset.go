package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// # thepile.jsonl
//
// {"text": "Heart Team: Joint Position of the Swiss Society of Cardiology and the Swiss Society of Cardiac Surgery.", "meta": {"pile_set_name": "PubMed Abstracts"}}

func RulesBasedSBD(text string) []string {
	// bin/sbd.py
	// Reads text from stdin and writes to stdout
	// with one sentence per line.

	pythonEnv := "env/bin/python"
	pythonScript := "bin/sbd.py"
	cmd := exec.Command(pythonEnv, pythonScript)
	cmd.Stdin = strings.NewReader(text)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(out.String(), "\n")
}

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

		// Split text into sentences
		sentences := RulesBasedSBD(result["text"].(string))
		for _, sentence := range sentences {
			fmt.Println(sentence)
		}

		// Print a blank line
		fmt.Println("\n")

	}
}
