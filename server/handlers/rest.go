package server

import (
	"fmt"
	"net/http"

	"github.com/thiagodsantos/gomockserver/constants"
)

func RESTHandler(r *http.Request, url string, requestBody []byte) (*http.Response, error) {
	// Only allow HTTP methods defined in constants
	if !constants.HTTPMethods[r.Method] {
		return nil, fmt.Errorf("method not allowed")
	}

	// Create new request
	req, err := http.NewRequest(r.Method, url, r.Body)

	// Return error when request creation fails
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Copy headers from original request to new request
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	var client = &http.Client{}
	var resp *http.Response

	// Make request to endpoint
	resp, err = client.Do(req)

	// Return error when request fails
	if err != nil {
		return nil, fmt.Errorf("error making request to endpoint: %v", err)
	}

	return resp, nil
}
