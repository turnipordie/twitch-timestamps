package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/turnipordie/twitch-ts-peaks/tsprocessor"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println(`usage: go run main.go "json-file-path" 2`)
		os.Exit(1)
	}

	jsonFilePath := os.Args[1]
	jsonData, err := os.ReadFile(jsonFilePath)
	if err != nil {
		fmt.Printf("reading JSON file: %v\n", err)
		os.Exit(1)
	}

	avgMult, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		fmt.Printf("reading avgMult: %v\n use a number like 2 or 3 :)", err)
		os.Exit(1)
	}

	var jsonCheck interface{}
	if err := json.Unmarshal(jsonData, &jsonCheck); err != nil {
		fmt.Printf("invalid json: %v\n", err)
		os.Exit(1)
	}

	timestamps := tsprocessor.Process(jsonData, avgMult)
	fmt.Printf("found %d possible highlights:\n", len(timestamps))
	for _, ts := range timestamps {
		fmt.Printf("%s\n", ts)
	}
}
