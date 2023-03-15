package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
)

// Requires unzstd to be installed

// Example The Pile download URLs:
// https://the-eye.eu/public/AI/pile/train/00.jsonl.zst
// https://the-eye.eu/public/AI/pile/train/29.jsonl.zst

func downloadFile(url string, filename string) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
}

func buildDownloadURL(index int) string {
	var url string = "https://the-eye.eu/public/AI/pile/train/"
	var indexString string = fmt.Sprintf("%02d", index)
	var extension string = ".jsonl.zst"
	var downloadURL string = url + indexString + extension
	return downloadURL
}

func unzstd(filename string) {
	// Example shell command
	// unzstd thepile.jsonl.zst

	// Create command
	cmd := exec.Command("unzstd", filename)

	// Run command
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func getThePile() {
	// Select random URL index from 0 to 29 inclusive
	index := rand.Intn(30)
	downloadURL := buildDownloadURL(index)
	fmt.Println(downloadURL)

	// Download file from URL
	filename := "thepile.jsonl.zst"
	downloadFile(downloadURL, filename)

	// Unzip file
	unzstd(filename)
}

func main() {
	getThePile()
	fmt.Println("Done")
}
