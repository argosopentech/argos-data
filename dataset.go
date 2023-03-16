package main

import (
	"fmt"
)

type Datum struct {
	source string
	target string
}

type Dataset struct {
	data []Datum
}

func main() {
	fmt.Println("Hello, playground")
}
