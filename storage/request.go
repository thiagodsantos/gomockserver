package storage

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/thiagodsantos/gomockserver/structs"
	"github.com/thiagodsantos/gomockserver/utils"
)

// Save request info to file in JSON format with filename request_<url>.json
func SaveRequest(url string, request *http.Request, responseTime string) error {
	requestInfo := structs.Request{
		URL:          url,
		Method:       request.Method,
		Headers:      request.Header,
		Body:         request.Body,
		ResponseTime: responseTime,
	}

	// Encode request info to JSON format
	requestInfoJSON, err := json.MarshalIndent(requestInfo, "", "  ")
	if err != nil {
		fmt.Println("Error encoding request info to JSON:", err)
		return err
	}

	// Format filename with prefix and URL
	requestFilename := utils.FormatFilename("request", url)

	// Save request info to file
	err = utils.SaveFile(requestFilename, requestInfoJSON)
	if err != nil {
		fmt.Println("Error writing request info to file:", err)
		return err
	}
	fmt.Println("Request info saved to", requestFilename)

	return nil
}
