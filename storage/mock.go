package storage

import (
	"encoding/json"
	"fmt"

	"github.com/thiagodsantos/gomockserver/constants"
	"github.com/thiagodsantos/gomockserver/structs"
	"github.com/thiagodsantos/gomockserver/utils"
)

// Get mock response from file
func GetMockResponse(url string) ([]byte, int, error) {
	// Format filename with prefix and URL
	responseFilename := utils.FormatFilename(constants.ResponseFileName, url)

	// Check if file exists in current directory
	if !utils.FileExists(responseFilename) {
		return defaultResponse()
	}

	var response structs.Response
	// Read JSON file data from response_<url>.json
	responseJSON, err := utils.ReadJSONFile(responseFilename, &response)
	if err != nil {
		fmt.Println("Error reading mock:", err)
		return nil, 500, err
	}

	if len(responseJSON) == 0 {
		return defaultResponse()
	}

	// Return mock response from file
	responseFileJSON, err := json.MarshalIndent(response.Body, "", "  ")
	if err != nil {
		fmt.Println("Error encoding mock response to JSON:", err)
		return nil, 500, err
	}

	return responseFileJSON, response.StatusCode, nil
}

func defaultResponse() ([]byte, int, error) {
	return nil, 200, nil
}
