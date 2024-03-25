package utils

import (
	"encoding/json"
	"fmt"
)

// Read JSON file data from file
func ReadJSONFile(filename string, v interface{}) ([]byte, error) {
	// Read file data
	data, err := ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return data, err
	}

	// Parse JSON data
	err = json.Unmarshal(data, v)
	if err != nil {
		fmt.Println("Error unmarshal JSON file:", err)
		return data, err
	}

	return data, nil
}
