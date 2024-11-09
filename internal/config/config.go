package config

import (
	"encoding/json"
	"log"
	"os"
)

func ReadFromFile(path string) (map[string][]byte, error) {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_SYNC, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	log.Printf("Opened file: %v\n", path)

	var results map[string][]byte
	if err := json.NewDecoder(file).Decode(&results); err != nil {
		return nil, err
	}
	log.Print("Config file parsed successfully\n")
	return results, nil
}
