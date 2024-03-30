package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/thiagodsantos/gomockserver/config"
	"github.com/thiagodsantos/gomockserver/constants"
	"github.com/thiagodsantos/gomockserver/server"
	"github.com/thiagodsantos/gomockserver/storage"
	"github.com/thiagodsantos/gomockserver/utils"
)

func handlerGenerate(w http.ResponseWriter, r *http.Request) {
	// Get host hostConfig from hosts.hostConfig.json
	hostConfig, err := config.GetHostConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting host config: %v", err), http.StatusInternalServerError)
		return
	}

	// Create URL from host and path
	path := r.URL.Path
	urlParsed, err := url.Parse(hostConfig.Url + path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing URL: %v", err), http.StatusInternalServerError)
		return
	}

	url := urlParsed.String()

	// Generate empty request file
	err = storage.GenerateEmptyRequestFile(url)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating empty request file: %v", err), http.StatusInternalServerError)
		return
	}

	// Generate empty response file
	err = storage.GenerateEmptyResponseFile(url)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating empty response file: %v", err), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Empty request and response files generated"))
}

// Handler function to handle requests
func handler(w http.ResponseWriter, r *http.Request) {
	var err error

	// Get body from request
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Get host config from hosts.config.json
	config, err := config.GetHostConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting host config: %v", err), http.StatusInternalServerError)
		return
	}

	// Check if request is for generating empty files
	if strings.HasPrefix(r.URL.Path, config.GeneratePath) {
		handlerGenerate(w, r)
		return
	}

	start := time.Now()

	var hostURL string

	if config.EnableGraphql {
		urlGraphQL, err := url.Parse(config.Url + config.GraphQLPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing GraphQL URL: %v", err), http.StatusInternalServerError)
			return
		}
		hostURL = urlGraphQL.String()
	}

	if config.EnableREST {
		path := r.URL.Path
		urlParsed, err := url.Parse(config.Url + path)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing URL: %v", err), http.StatusInternalServerError)
			return
		}
		hostURL = urlParsed.String()
	}

	var operationNameHashed string

	// Return mock response from file
	if config.UseMock {
		if config.EnableGraphql {
			graphqlRequestBody, err := server.GetGraphQLRequestBody(w, r, requestBody)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error getting GraphQL request: %v", err), http.StatusInternalServerError)
				return
			}

			operationNameHashed, err = graphqlRequestBody.GetOperationNameHashed()
			if err != nil {
				http.Error(w, fmt.Sprintf("Error getting operation: %v", err), http.StatusInternalServerError)
				return
			}
		}

		responseFileJSON, statusCode, err := storage.GetMockResponse(hostURL, operationNameHashed)
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

	// Make GraphQL request
	if config.EnableGraphql {
		resp, err = server.GraphqlHandler(w, r, hostURL, requestBody)
	}

	// Make REST request
	if config.EnableREST {
		resp, err = server.RESTHandler(w, r, hostURL, requestBody)
	}

	defer resp.Body.Close()

	if err != nil {
		http.Error(w, fmt.Sprintf("Error making request to endpoint: %v", err), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)

	// Save request to file
	err = storage.SaveRequest(hostURL, r, duration.String(), requestBody, operationNameHashed)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving request: %v", err), http.StatusInternalServerError)
		return
	}

	// Save response to file
	responseData, responseBody, err := storage.SaveResponse(hostURL, resp, duration.String(), operationNameHashed)
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

func setup() {
	// Generate empty server config file
	err := storage.GenerateEmptyServerConfigFile()
	if err != nil {
		fmt.Println("Error generating empty server config file:", err)
		return
	}

	// Generate empty hosts config file
	err = storage.GenerateEmptyHostsConfigFile()
	if err != nil {
		fmt.Println("Error generating empty hosts config file:", err)
		return
	}
}

func main() {
	// Setup
	setup()

	// Get server config
	serverConfig, _ := config.GetServerConfig()
	port := fmt.Sprintf(":%d", serverConfig.Port)

	// Start server
	http.HandleFunc(serverConfig.Path, handler)
	fmt.Println(utils.Format(utils.PURPLE, "Mock server running on "+port+"\n"))
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
}

/**

http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
	server.GraphqlHandler(w, r, "https://api.spacex.land/graphql/")
})
*/
