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

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func saveFile(filename string, data []byte) error {
	err := os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}

	return nil
}

func readFile(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}
	return data, nil
}

func readJSONFile(filename string, v interface{}) ([]byte, error) {
	data, err := readFile(filename)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return data, err
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		fmt.Println("Error unmarshal JSON file:", err)
		return data, err
	}

	return data, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	var configs []Config
	_, err := readJSONFile("hosts.config.json", &configs)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
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
		fileExists := fileExists(responseFilename)
		if fileExists {
			var responseFileInfo ResponseInfo
			responseFile, err := readJSONFile(responseFilename, &responseFileInfo)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error reading mock: %v", err), http.StatusInternalServerError)
				return
			}

			if len(responseFile) > 0 {
				durationMock := time.Since(start)
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
	}

	var resp *http.Response

	if r.Method == "GET" {
		resp, err = http.Get(url)
	}

	if r.Method == "POST" {
		resp, err = http.Post(url, "application/json", r.Body)
	}

	resp, err = http.Get(url)
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

	saveFile(responseFilename, responseInfoJSON)
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

	err = saveFile(requestFilename, requestInfoJSON)
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
