package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/thiagodsantos/gomockserver/constants"
)

// Config struct from hosts.config.json
type Config struct {
	Host    string   `json:"host"`
	Paths   []string `json:"paths"`
	Enabled bool     `json:"enabled"`
	UseMock bool     `json:"use_mock"`
}

// RequestInfo struct to save request info to file in JSON format with filename request_<url>.json
type RequestInfo struct {
	URL          string              `json:"url"`
	Method       string              `json:"method"`
	Headers      map[string][]string `json:"headers"`
	Body         interface{}         `json:"body,omitempty"`
	ResponseTime string              `json:"response_time"`
}

// ResponseInfo struct to save response info to file in JSON format with filename response_<url>.json
type ResponseInfo struct {
	URL          string              `json:"url"`
	Method       string              `json:"method"`
	Headers      map[string][]string `json:"headers"`
	Body         interface{}         `json:"body,omitempty"`
	StatusCode   int                 `json:"status_code"`
	ResponseTime string              `json:"response_time"`
}

// Check if file exists in current directory
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// Save data to file
func saveFile(filename string, data []byte) error {
	err := os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}

	return nil
}

// Read file data from file
func readFile(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}
	return data, nil
}

// Read JSON file data from file
func readJSONFile(filename string, v interface{}) ([]byte, error) {
	// Read file data
	data, err := readFile(filename)
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

// Format filename with prefix and URL
func formatFilename(prefix string, url string) string {
	// Replace / and : with -
	filename := fmt.Sprintf("%s_%s.%s", prefix, strings.ReplaceAll(url, "/", "-"), constants.JSONExtension)
	filename = strings.ReplaceAll(filename, ":", "-")

	return filename
}

// Get host config from hosts.config.json
func getHostConfig() (Config, error) {
	var configs []Config

	// Read JSON file data from hosts.config.json
	_, err := readJSONFile("hosts.config.json", &configs)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return Config{}, err
	}

	// Get host config enabled
	var hostCount int
	var hostConfig Config
	for _, config := range configs {
		if config.Enabled {
			hostConfig = config
			hostCount++
		}
	}

	// Return error if more than one host is enabled
	if hostCount > 1 {
		return Config{}, fmt.Errorf("more than one host enabled in hosts.config")
	}

	return hostConfig, nil
}

// Save request info to file in JSON format with filename request_<url>.json
func saveRequestInfo(url string, request *http.Request, responseTime string) error {
	requestInfo := RequestInfo{
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
	requestFilename := formatFilename("request", url)

	// Save request info to file
	err = saveFile(requestFilename, requestInfoJSON)
	if err != nil {
		fmt.Println("Error writing request info to file:", err)
		return err
	}
	fmt.Println("Request info saved to", requestFilename)

	return nil
}

// Save response info to file
func saveResponseInfo(url string, response *http.Response, responseTime string) (ResponseInfo, []byte, error) {
	// Read response body data from response
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ResponseInfo{}, nil, err
	}

	// Decode response body to JSON format
	body := map[string]interface{}{}
	json.Unmarshal(responseBody, &body)

	responseInfo := ResponseInfo{
		URL:          url,
		Method:       response.Request.Method,
		Headers:      response.Header,
		Body:         body,
		StatusCode:   response.StatusCode,
		ResponseTime: responseTime,
	}

	// Encode response info to JSON format
	responseInfoJSON, err := json.MarshalIndent(responseInfo, "", "  ")
	if err != nil {
		fmt.Println("Error encoding response info to JSON:", err)
		return ResponseInfo{}, nil, err
	}

	responseFilename := formatFilename("response", url)

	// Save response info to file
	err = saveFile(responseFilename, responseInfoJSON)
	if err != nil {
		fmt.Println("Error writing response info to file:", err)
		return ResponseInfo{}, nil, err
	}
	fmt.Println("Response info saved to", responseFilename)

	return responseInfo, responseBody, nil
}

// Get mock response from file
func getMockResponse(url string) ([]byte, int, error) {
	// Format filename with prefix and URL
	responseFilename := formatFilename("response", url)

	// Check if file exists in current directory
	fileExists := fileExists(responseFilename)

	if fileExists {
		var responseFileInfo ResponseInfo
		// Read JSON file data from response_<url>.json
		responseFile, err := readJSONFile(responseFilename, &responseFileInfo)
		if err != nil {
			fmt.Println("Error reading mock:", err)
			return nil, 500, err
		}

		// Return mock response from file
		if len(responseFile) > 0 {
			responseFileJSON, err := json.Marshal(responseFileInfo.Body)
			if err != nil {
				fmt.Println("Error encoding mock response to JSON:", err)
				return nil, 500, err
			}

			return responseFileJSON, responseFileInfo.StatusCode, nil
		}
	}

	return nil, 200, nil
}

// Handler function to handle requests
func handler(w http.ResponseWriter, r *http.Request) {
	// Get host config from hosts.config.json
	config, err := getHostConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting host config: %v", err), http.StatusInternalServerError)
		return
	}

	// Create URL from host and path
	path := r.URL.Path
	urlParsed, err := url.Parse(config.Host + path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing URL: %v", err), http.StatusInternalServerError)
		return
	}

	start := time.Now()
	url := urlParsed.String()

	// Return mock response from file
	if config.UseMock {
		responseFileJSON, statusCode, err := getMockResponse(url)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting mock response: %v", err), http.StatusInternalServerError)
			return
		}

		// Return response when status code is 4xx or 5xx
		if statusCode >= 400 {
			http.Error(w, string(responseFileJSON), statusCode)
			return
		}

		// Return mock response
		if responseFileJSON != nil {
			w.Header().Set("Content-Type", constants.JSONContentType)
			w.Write(responseFileJSON)
			return
		}
	}

	var resp *http.Response

	// Make request to endpoint
	if r.Method == constants.MethodGet {
		resp, err = http.Get(url)
	}

	if r.Method == constants.MethodPost {
		resp, err = http.Post(url, constants.JSONContentType, r.Body)
	}

	// Return error when request fails
	if err != nil {
		http.Error(w, fmt.Sprintf("Error making request to endpoint: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	duration := time.Since(start)

	// Save request info to file
	err = saveRequestInfo(url, r, duration.String())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving request info: %v", err), http.StatusInternalServerError)
		return
	}

	// Save response info to file
	responseInfo, responseBody, err := saveResponseInfo(url, resp, duration.String())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving response info: %v", err), http.StatusInternalServerError)
		return
	}

	// Return response when status code is 4xx or 5xx
	if responseInfo.StatusCode >= 400 {
		http.Error(w, string(responseBody), responseInfo.StatusCode)
		return
	}

	// Return response from endpoint
	w.Header().Set("Content-Type", constants.JSONContentType)
	w.Write(responseBody)
}

func main() {
	// Start server
	http.HandleFunc("/", handler)
	fmt.Println("Server running on", constants.ProxyServerPort)
	http.ListenAndServe(constants.ProxyServerPort, nil)
}
