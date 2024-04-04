package server

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/thiagodsantos/gomockserver/config"
	"github.com/thiagodsantos/gomockserver/constants"
	handlers "github.com/thiagodsantos/gomockserver/server/handlers"
	"github.com/thiagodsantos/gomockserver/storage"
)

func handlerGenerate(w http.ResponseWriter, r *http.Request) {
	// Get host hostConfig from hosts.hostConfig.json
	hostConfig, err := config.GetHostConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting host config: %v", err), http.StatusInternalServerError)
		return
	}

	// Create URL from host and path
	url, err := hostConfig.GetHostURL(r.URL.Path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting host URL: %v", err), http.StatusInternalServerError)
		return
	}

	// Generate empty request file
	err = storage.GenerateEmptyRequestFile(url, hostConfig.OutputFolder)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating empty request file: %v", err), http.StatusInternalServerError)
		return
	}

	// Generate empty response file
	err = storage.GenerateEmptyResponseFile(url, hostConfig.OutputFolder)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating empty response file: %v", err), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Empty request and response files generated"))
}

// Handler function to handle requests
func Handler(w http.ResponseWriter, r *http.Request) {
	var err error

	// Get host config from hosts.config.json
	config, err := config.GetHostConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting host config: %v", err), http.StatusInternalServerError)
		return
	}

	// Check if request is for generating empty files
	if config.IsGeneratePath(r.URL.Path) {
		handlerGenerate(w, r)
		return
	}

	start := time.Now()

	// Get host URL
	hostURL, err := config.GetHostURL(r.URL.Path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting host URL: %v", err), http.StatusInternalServerError)
		return
	}

	// Get body from request
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var operationName string

	// Return mock response from file
	if config.EnableMock {

		// Get operation name for GraphQL
		if config.EnableGraphql {
			operationName, err = handlers.GetOperationName(w, r, requestBody)
			if err != nil {
				http.Error(w, "Error getting operation name", http.StatusInternalServerError)
				return
			}
		}

		// Get mock response from file
		responseFileJSON, statusCode, err := storage.GetMockResponse(hostURL, operationName, config.OutputFolder)
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
			w.Header().Set(constants.HeaderContentType, constants.JSONContentType)
			w.Write(responseFileJSON)
			return
		}
	}

	// Execute request (REST or GraphQL)
	resp, err := Execute(w, r, config, hostURL, requestBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing request: %v", err), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	// Start timer
	duration := time.Since(start)

	// Save request to file
	err = storage.SaveRequest(hostURL, r, duration.String(), requestBody, operationName, config.OutputFolder)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving request: %v", err), http.StatusInternalServerError)
		return
	}

	// Save response to file
	responseData, responseBody, err := storage.SaveResponse(hostURL, resp, duration.String(), operationName, config.OutputFolder)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving response: %v", err), http.StatusInternalServerError)
		return
	}

	// Return response when status code is 4xx or 5xx
	if responseData.StatusCode >= 400 {
		http.Error(w, string(responseBody), responseData.StatusCode)
		return
	}

	// Return response from endpoint
	w.Header().Set(constants.HeaderContentType, constants.JSONContentType)
	w.Write(responseBody)
}
