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
)

type Config struct {
	Host    string   `json:"host"`
	Paths   []string `json:"paths"`
	Enabled bool     `json:"enabled"`
	UseMock bool     `json:"use_mock"`
}

type RequestInfo struct {
	URL          string              `json:"url"`
	Method       string              `json:"method"`
	Headers      map[string][]string `json:"headers"`
	Body         interface{}         `json:"body,omitempty"`
	ResponseTime string              `json:"response_time"`
}

type ResponseInfo struct {
	URL          string              `json:"url"`
	Method       string              `json:"method"`
	Headers      map[string][]string `json:"headers"`
	Body         interface{}         `json:"body,omitempty"`
	ResponseTime string              `json:"response_time"`
	StatusCode   int                 `json:"status_code"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	jsonData, err := os.ReadFile("hosts.config.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	var configs []Config
	err = json.Unmarshal(jsonData, &configs)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	var host string
	var hostCount int
	var useMock bool
	for _, config := range configs {
		if config.Enabled {
			useMock = config.UseMock
			host = config.Host
			hostCount++
		}
	}

	// TODO: Support when more than one host is enabled, getting the host by path
	if hostCount != 1 {
		http.Error(w, "Exactly one host must be enabled in the config", http.StatusInternalServerError)
		return
	}

	path := r.URL.Path

	urlParsed, err := url.Parse(host + path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing URL: %v", err), http.StatusInternalServerError)
		return
	}

	url := urlParsed.String()

	responseFilename := fmt.Sprintf("response_%s.json", strings.ReplaceAll(url, "/", "-"))
	responseFilename = strings.ReplaceAll(responseFilename, ":", "-")

	start := time.Now()

	if useMock {
		durationMock := time.Since(start)
		responseFile, _ := os.ReadFile(responseFilename)
		if len(responseFile) > 0 {
			var responseFileInfo ResponseInfo
			err := json.Unmarshal(responseFile, &responseFileInfo)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error reading mock: %v", err), http.StatusInternalServerError)
				return
			}

			fmt.Println("Reading from mockfile duration", durationMock)
			responseFileJSON, _ := json.Marshal(responseFileInfo.Body)

			if responseFileInfo.StatusCode >= 400 {
				http.Error(w, string(responseFileJSON), responseFileInfo.StatusCode)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(responseFileJSON)
			return
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error making request to endpoint: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	duration := time.Since(start)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading response body: %v", err), http.StatusInternalServerError)
		return
	}

	body := map[string]interface{}{}
	json.Unmarshal(responseBody, &body)

	responseInfo := ResponseInfo{
		URL:          url,
		Method:       resp.Request.Method,
		Headers:      resp.Header,
		Body:         body,
		StatusCode:   resp.StatusCode,
		ResponseTime: duration.String(),
	}

	responseInfoJSON, err := json.MarshalIndent(responseInfo, "", "  ")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response info to JSON: %v", err), http.StatusInternalServerError)
		return
	}

	err = os.WriteFile(responseFilename, responseInfoJSON, 0644)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error writing response to file: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Println("Response saved to", responseFilename)

	requestInfo := RequestInfo{
		URL:          url,
		Method:       r.Method,
		Headers:      r.Header,
		ResponseTime: duration.String(),
	}

	requestInfoJSON, err := json.MarshalIndent(requestInfo, "", "  ")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding request info to JSON: %v", err), http.StatusInternalServerError)
		return
	}

	requestFilename := fmt.Sprintf("request_%s.json", strings.ReplaceAll(url, "/", "-"))
	requestFilename = strings.ReplaceAll(requestFilename, ":", "-")

	err = os.WriteFile(requestFilename, requestInfoJSON, 0644)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error writing request info to file: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Println("Request info saved to", requestFilename)

	if responseInfo.StatusCode >= 400 {
		http.Error(w, string(responseBody), responseInfo.StatusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
