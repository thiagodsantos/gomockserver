package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/thiagodsantos/gomockserver/constants"
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

func SaveJSONFile(folder string, filename string, v interface{}) error {
	// Encode JSON data
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("Error encoding JSON data:", err)
		return err
	}

	outputFolder := filepath.Join(constants.OutputFolder, folder)

	if outputFolder != "" {
		if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
			if err := os.MkdirAll(outputFolder, 0755); err != nil {
				fmt.Println("Error creating folder:", err)
				return err
			}
		}
	}

	outputFolder = filepath.Join(outputFolder, filename)

	// Save JSON data to file
	err = SaveFile(outputFolder, data)
	if err != nil {
		fmt.Println("Error writing JSON data to file:", err)
		return err
	}

	return nil
}
