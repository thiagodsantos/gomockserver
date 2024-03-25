package storage

import (
	"encoding/json"
	"fmt"

	"github.com/thiagodsantos/gomockserver/structs"
	"github.com/thiagodsantos/gomockserver/utils"
)

// Get mock response from file
func GetMockResponse(url string) ([]byte, int, error) {
	// Format filename with prefix and URL
	responseFilename := utils.FormatFilename("response", url)

	// Check if file exists in current directory
	fileExists := utils.FileExists(responseFilename)

	if !fileExists {
		return defaultResponse()
	}

	var responseFileInfo structs.Response
	// Read JSON file data from response_<url>.json
	responseFile, err := utils.ReadJSONFile(responseFilename, &responseFileInfo)
	if err != nil {
		fmt.Println("Error reading mock:", err)
		return nil, 500, err
	}

	if len(responseFile) == 0 {
		return defaultResponse()
	}

	// Return mock response from file
	responseFileJSON, err := json.MarshalIndent(responseFileInfo.Body, "", "  ")
	if err != nil {
		fmt.Println("Error encoding mock response to JSON:", err)
		return nil, 500, err
	}

	return responseFileJSON, responseFileInfo.StatusCode, nil
}

func defaultResponse() ([]byte, int, error) {
	return nil, 200, nil
}
