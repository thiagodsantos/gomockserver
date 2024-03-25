package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/thiagodsantos/gomockserver/config"
	"github.com/thiagodsantos/gomockserver/constants"
	"github.com/thiagodsantos/gomockserver/storage"
	"github.com/thiagodsantos/gomockserver/utils"
)

// Handler function to handle requests
func handler(w http.ResponseWriter, r *http.Request) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da solicitação", http.StatusInternalServerError)
		return
	}

	// Get host config from hosts.config.json
	config, err := config.GetHostConfig()
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
		responseFileJSON, statusCode, err := storage.GetMockResponse(url)
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
	err = storage.SaveRequest(url, r, duration.String(), requestBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving request info: %v", err), http.StatusInternalServerError)
		return
	}

	// Save response info to file
	responseInfo, responseBody, err := storage.SaveResponse(url, resp, duration.String())
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
	w.Header().Set(constants.HeaderContentType, constants.JSONContentType)
	w.Write(responseBody)
}

func main() {
	// Start server
	http.HandleFunc("/", handler)
	fmt.Println(utils.Format(utils.PURPLE, "Mock server running on "+constants.ProxyServerPort+"\n"))
	http.ListenAndServe(constants.ProxyServerPort, nil)
}
