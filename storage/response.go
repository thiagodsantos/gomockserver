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

// Save response info to file
func SaveResponse(url string, response *http.Response, responseTime string) (structs.Response, []byte, error) {
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

	// Encode response info to JSON format
	responseJSON, err := json.MarshalIndent(responseData, "", "  ")
	if err != nil {
		fmt.Println("Error encoding response info to JSON:", err)
		return structs.Response{}, nil, err
	}

	responseFilename := utils.FormatFilename(constants.ResponseFileName, url)

	// Save response info to file
	err = utils.SaveFile(responseFilename, responseJSON)
	if err != nil {
		fmt.Println("Error writing response info to file:", err)
		return structs.Response{}, nil, err
	}
	//fmt.Println("Response info saved to", responseFilename)
	fmt.Println(utils.Format(utils.BLUE, "Response info saved to "+responseFilename))

	return responseData, responseBody, nil
}
