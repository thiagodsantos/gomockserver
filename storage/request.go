package storage

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/thiagodsantos/gomockserver/constants"
	"github.com/thiagodsantos/gomockserver/structs"
	"github.com/thiagodsantos/gomockserver/utils"
)

// Save request info to file in JSON format with filename request_<url>.json
func SaveRequest(url string, request *http.Request, responseTime string, requestBody []byte) error {
	// Decode response body to JSON format
	body := map[string]interface{}{}
	json.Unmarshal(requestBody, &body)

	requestData := structs.Request{
		URL:          url,
		Method:       request.Method,
		Headers:      request.Header,
		Body:         body,
		ResponseTime: responseTime,
	}

	// Encode request info to JSON format
	requestJSON, err := json.MarshalIndent(requestData, "", "  ")
	if err != nil {
		fmt.Println("Error encoding request info to JSON:", err)
		return err
	}

	// Format filename with prefix and URL
	requestFilename := utils.FormatFilename(constants.RequestFileName, url)

	// Save request info to file
	err = utils.SaveFile(requestFilename, requestJSON)
	if err != nil {
		fmt.Println("Error writing request info to file:", err)
		return err
	}
	fmt.Println(utils.Format(utils.GREEN, "Request info saved to "+requestFilename))

	return nil
}
