package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/thiagodsantos/gomockserver/constants"
	"github.com/thiagodsantos/gomockserver/structs"
	"github.com/thiagodsantos/gomockserver/utils"
)

// Save response to file
func SaveResponse(url string, response *http.Response, responseTime string, suffix string) (structs.Response, []byte, error) {
	// Read response body data from response
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return structs.Response{}, nil, err
	}

	// Check if response content type is JSON
	if !strings.Contains(response.Header.Get(constants.HeaderContentType), constants.JSONContentType) {
		fmt.Println("Response content type allows only JSON format")
		return structs.Response{}, nil, nil
	}

	// Decode response body to JSON format
	body := map[string]interface{}{}
	json.Unmarshal(responseBody, &body)

	responseData := structs.Response{
		URL:          url,
		Method:       response.Request.Method,
		Headers:      response.Header,
		Body:         body,
		StatusCode:   response.StatusCode,
		ResponseTime: responseTime,
	}

	responseFilename := utils.FormatFilename(constants.ResponseFileName, url+suffix)

	err = utils.SaveJSONFile(responseFilename, responseData)
	if err != nil {
		fmt.Println("Error saving response to file:", err)
		return structs.Response{}, nil, err
	}
	fmt.Println(utils.Format(utils.BLUE, "Response saved to "+responseFilename))

	return responseData, responseBody, nil
}

func GenerateEmptyResponseFile(url string) error {
	responseData := structs.Response{
		URL:          url,
		Method:       "",
		Headers:      nil,
		Body:         map[string]interface{}{},
		StatusCode:   0,
		ResponseTime: "",
	}
	responseFilename := utils.FormatFilename(constants.ResponseFileName, url)

	err := utils.SaveJSONFile(responseFilename, responseData)
	if err != nil {
		fmt.Println("Error saving empty response to file:", err)
		return err
	}
	fmt.Println(utils.Format(utils.BLUE, "Response saved to "+responseFilename))

	return nil
}
