package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/thiagodsantos/gomockserver/constants"
)

// Format filename with prefix and URL
func FormatFilename(prefix string, url string) string {
	// Replace / and : with -
	filename := fmt.Sprintf("%s_%s.%s", prefix, strings.ReplaceAll(url, "/", "-"), constants.JSONExtension)
	filename = strings.ReplaceAll(filename, ":", "-")

	return filename
}

// Check if file exists in current directory
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// Save data to file
func SaveFile(filename string, data []byte) error {
	err := os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}

	return nil
}

// Read file data from file
func ReadFile(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}
	return data, nil
}
