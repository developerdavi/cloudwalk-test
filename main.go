package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/developerdavi/cloudwalk-test/lib/parser"
)

func main() {
	matches, err := parser.Parse("input/qgames.log")

	if err != nil {
		log.Fatal(err)
	}

	data, err := json.MarshalIndent(matches, "", "  ")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data))
}
