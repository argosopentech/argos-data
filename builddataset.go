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

func (d DataPackage) String() string {
	return fmt.Sprintf("Name: %s, Type: %s, FromCode: %s, ToCode: %s, Size: %d, Reference: %s, Links: %s", d.Name, d.Type, d.FromCode, d.ToCode, d.Size, d.Reference, d.Links)
}

func WriteDataToFile(dataPackage DataPackage, dataPath string, langCode string, f *zip.File) {
	// Make buffered reader for source file
	fReader, err := f.Open()
	if err != nil {
		fmt.Println(err)
	}
	defer fReader.Close()

	out, err := os.OpenFile(dataPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer out.Close()

	var langPrefix string = "__" + langCode + "__"

	var i int = 0

	// Read zipped file line by line and write to output file
	scanner := bufio.NewScanner(fReader)

	scanner.Split(bufio.ScanLines)
	var text []string

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	for _, each_ln := range text {
		var line string = langPrefix + each_ln + "\n"

		_, err = out.WriteString(line)
		if err != nil {
			fmt.Println(err)
		}

		i++
	}

	if i != len(text) {
		fmt.Println("Error: i != len(text)")
	}
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
			WriteDataToFile(dataPackage, sourcePath, dataPackage.FromCode, f)
		} else if strings.Contains(f.Name, "target") {
			WriteDataToFile(dataPackage, targetPath, dataPackage.ToCode, f)
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

	// Sort data packages by size smallest to largest
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data)-1; j++ {
			if data[j].Size > data[j+1].Size {
				data[j], data[j+1] = data[j+1], data[j]
			}
		}
	}

	// Limit total dataset size to 100 million lines
	maxDataSize := 100 * 1000000
	var cummulativeDataSize int = 0
	for i := 0; i < len(data); i++ {
		cummulativeDataSize += data[i].Size
		if cummulativeDataSize > maxDataSize {
			data = data[:i]
			break
		}
	}

	// Loop through all data packages and append to dataset
	for _, d := range data {
		fmt.Println(d)
		AppendDataPackageToDataset(d)
	}
}
