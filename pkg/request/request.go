package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/thiagodsantos/gomockserver/pkg/config"
)

type ResponseFile struct {
	Headers    map[string]string `json:"headers"`
	Response   interface{}       `json:"body,omitempty"`
	StatusCode int               `json:"status_code"`
	Url        string            `json:"url"`
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	var err error

	endpoint := r.URL.Path
	originalHost := getOriginalHost(r)

	serviceConfig, err := config.GetServiceConfig(originalHost)
	if err != nil {
		http.Error(w, "Service configuration not found", http.StatusBadGateway)
		return
	}

	url := fmt.Sprintf("%s%s", serviceConfig.Host, endpoint)

	var responseData []byte

	if serviceConfig.Enabled {
		response, err := loadResponse(w, url, serviceConfig.ResponseFolder)

		if err == nil && response.Response != nil {
			responseData, err = json.Marshal(response.Response)
			if err == nil {
				w.WriteHeader(response.StatusCode)
				w.Write(responseData)
				return
			}
		}
	}

	req, err := http.NewRequestWithContext(r.Context(), r.Method, url, r.Body)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error sending request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	if serviceConfig.Store {
		response, err := saveResponse(url, resp, serviceConfig.ResponseFolder)
		if err != nil {
			fmt.Println("Error saving response:", err)
		}

		responseData, err = json.Marshal(response.Response)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}
	}

	w.Write(responseData)
}

func getOriginalHost(r *http.Request) string {
	scheme := "http"

	if r.TLS != nil {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s", scheme, r.Host)
}

func snakeCase(input string) string {
	output := strings.ReplaceAll(input, "/", "_")
	output = strings.ReplaceAll(output, "://", "-")
	output = strings.ToLower(output)

	return output
}

func getFilepath(folder string, url string) string {
	return filepath.Join(".output", folder, snakeCase(url)+".json")
}

func saveResponse(url string, resp *http.Response, folder string) (ResponseFile, error) {
	filePath := getFilepath(folder, url)

	var err error

	err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		fmt.Println("Error creating response folder:", err)
		return ResponseFile{}, err
	}

	var responseBody []byte

	responseBody, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ResponseFile{}, err
	}

	headers := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	response := map[string]interface{}{}

	if err := json.Unmarshal(responseBody, &response); err != nil {
		fmt.Println("Erro ao decodificar JSON:", err)
		return ResponseFile{}, err
	}

	responseFile := ResponseFile{
		Headers:    headers,
		Response:   response,
		StatusCode: resp.StatusCode,
		Url:        url,
	}

	data, err := json.MarshalIndent(responseFile, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return ResponseFile{}, err
	}

	err = os.WriteFile(filePath, data, os.ModePerm)
	if err != nil {
		fmt.Println("Error writing response file:", err)
		return ResponseFile{}, err
	}

	return responseFile, nil
}

func loadResponse(w http.ResponseWriter, url string, folder string) (ResponseFile, error) {
	filePath := getFilepath(folder, url)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return ResponseFile{}, nil
	}

	response, err := os.ReadFile(filePath)

	if err != nil {
		fmt.Println("Error reading response file:", err)
		return ResponseFile{}, nil
	}

	res := ResponseFile{}
	if err := json.Unmarshal(response, &res); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return ResponseFile{}, nil
	}

	for key, value := range res.Headers {
		w.Header().Add(key, value)
	}

	return res, nil
}
