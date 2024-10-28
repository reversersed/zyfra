package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func Read(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fmt.Printf("Opened file: %v\n", path)

	var results map[string]string
	if err := json.NewDecoder(file).Decode(&results); err != nil {
		return nil, err
	}
	fmt.Print("Config file parsed successfully\n")
	return results, nil
}
