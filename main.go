package main

import (
	"encoding/json"
	"fmt"

	"github.com/developerdavi/cloudwalk-test/lib/parser"
)

func main() {
	matches := parser.Parse()

	data, error := json.MarshalIndent(matches, "", "  ")

	if error != nil {
		fmt.Println(error)
	}

	fmt.Println(string(data))
}
