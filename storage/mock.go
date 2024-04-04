package storage

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/thiagodsantos/gomockserver/constants"
	"github.com/thiagodsantos/gomockserver/structs"
	"github.com/thiagodsantos/gomockserver/utils"
)

// Get mock response from file
func GetMockResponse(url string, operationName string, folder string) ([]byte, int, error) {
	// Format filename with prefix and URL
	responseFilename := utils.FormatFilename(constants.ResponseFileName, url+operationName)

	outputFile := filepath.Join(constants.OutputFolder, folder, responseFilename)

	// Check if file exists in current directory
	if !utils.FileExists(outputFile) {
		return defaultResponse()
	}

	var response structs.Response
	// Read JSON file data from response_<url>.json
	responseJSON, err := utils.ReadJSONFile(outputFile, &response)
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

	duration, err := time.ParseDuration(response.ResponseTime)
	if err != nil {
		return responseFileJSON, response.StatusCode, nil
	}

	time.Sleep(duration)

	return responseFileJSON, response.StatusCode, nil
}

func defaultResponse() ([]byte, int, error) {
	return nil, 200, nil
}
