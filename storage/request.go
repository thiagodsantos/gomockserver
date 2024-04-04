package storage

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/thiagodsantos/gomockserver/constants"
	"github.com/thiagodsantos/gomockserver/structs"
	"github.com/thiagodsantos/gomockserver/utils"
)

// Save request to file in JSON format with filename request_<url>.json
func SaveRequest(url string, request *http.Request, responseTime string, requestBody []byte, suffix string, folder string) error {
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

	requestFilename := utils.FormatFilename(constants.RequestFileName, url+suffix)

	err := utils.SaveJSONFile(folder, requestFilename, requestData)
	if err != nil {
		fmt.Println("Error saving request to file:", err)
		return err
	}
	fmt.Println(utils.Format(utils.GREEN, "Request saved to "+requestFilename))

	return nil
}

func GenerateEmptyRequestFile(url string, folder string) error {
	requestData := structs.Request{
		URL:          url,
		Method:       "",
		Headers:      nil,
		Body:         map[string]interface{}{},
		ResponseTime: "",
	}

	requestFilename := utils.FormatFilename(constants.RequestFileName, url)

	err := utils.SaveJSONFile(folder, requestFilename, requestData)
	if err != nil {
		fmt.Println("Error saving request to file:", err)
		return err
	}
	fmt.Println(utils.Format(utils.GREEN, "Request saved to "+requestFilename))

	return nil
}
