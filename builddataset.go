package main

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type DataPackage struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	FromCode  string   `json:"from_code"`
	ToCode    string   `json:"to_code"`
	Size      int      `json:"size"`
	Reference string   `json:"reference"`
	Links     []string `json:"links"`
}

func AppendDataPackageToDataset(dataPackage DataPackage) {
	// Download zip file at link and save to disk
	var zipPackagePath string = "dataPackage.argosdata"
	out, err := os.Create(zipPackagePath)
	if err != nil {
		fmt.Println(err)
	}
	defer out.Close()

	resp, err := http.Get(dataPackage.Links[0])
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)

	// Example data package format:
	// data-europarl-en_de/
	//   source (text file each line is a sentence)
	//   target (text file each line is a translated sentence corresponding to source)
	//   metadata.json

	// Open zip file
	r, err := zip.OpenReader(zipPackagePath)
	if err != nil {
		fmt.Println(err)
	}
	defer r.Close()

	// Example format for generated dataset:
	// source
	// __en__Hello World
	// __en__How are you?
	// target
	// __de__Hallo Welt
	// __de__Wie geht es dir?

	var sourcePath string = "source"
	var targetPath string = "target"

	// Extract source and target data from zip file
	for _, f := range r.File {
		if strings.Contains(f.Name, "source") {
			// Make buffered reader for source file
			fReader, err := f.Open()
			if err != nil {
				fmt.Println(err)
			}
			defer fReader.Close()

			out, err := os.Create(sourcePath)
			if err != nil {
				fmt.Println(err)
			}
			defer out.Close()

			var langCode string = dataPackage.FromCode
			var langPrefix string = "__" + langCode + "__"

			// Loop through each line until newline character in source file and add language prefix
			for {
				line, err := bufio.NewReader(fReader).ReadString('\n')
				if err != nil {
					break
				}
				_, err = out.WriteString(langPrefix + line)
				if err != nil {
					fmt.Println(err)
				}
			}
		} else if strings.Contains(f.Name, "target") {
			// Make buffered reader for target file
			fReader, err := f.Open()
			if err != nil {
				fmt.Println(err)
			}
			defer fReader.Close()

			out, err := os.Create(targetPath)
			if err != nil {
				fmt.Println(err)
			}
			defer out.Close()

			var langCode string = dataPackage.ToCode
			var langPrefix string = "__" + langCode + "__"

			// Loop through each line until newline character in target file and add language prefix
			for {
				line, err := bufio.NewReader(fReader).ReadString('\n')
				if err != nil {
					break
				}
				_, err = out.WriteString(langPrefix + line)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
	// Delete zip file
	err = os.Remove(zipPackagePath)
}

func main() {
	// Load JSON file from URL
	resp, err := http.Get("https://raw.githubusercontent.com/argosopentech/argos-train/master/data-index.json")

	if err != nil {
		fmt.Println(err)
	}

	// Parse JSON file
	var data []DataPackage
	jsonParser := json.NewDecoder(resp.Body)
	if err = jsonParser.Decode(&data); err != nil {
		fmt.Println(err)
	}

	/*
		// Select data package with smallest size value
		var dataPackage DataPackage
		for _, d := range data {
			if dataPackage.Size == 0 || d.Size < dataPackage.Size {
				dataPackage = d
			}
		}
		fmt.Println(dataPackage)
		// Append data package to dataset
		AppendDataPackageToDataset(dataPackage)
	*/

	// Sort data packages by size smallest to largest
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data)-1; j++ {
			if data[j].Size > data[j+1].Size {
				data[j], data[j+1] = data[j+1], data[j]
			}
		}
	}

	// Loop through all data packages and append to dataset
	for _, d := range data {
		fmt.Println(d.Name)
		AppendDataPackageToDataset(d)
	}
}
