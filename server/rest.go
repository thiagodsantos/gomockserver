package server

import (
	"fmt"
	"net/http"

	"github.com/thiagodsantos/gomockserver/constants"
)

func RESTHandler(w http.ResponseWriter, r *http.Request, url string, requestBody []byte) (*http.Response, error) {
	var err error
	var resp *http.Response

	// Make request to endpoint
	if r.Method == constants.MethodGet {
		resp, err = http.Get(url)
	}

	if r.Method == constants.MethodPost {
		resp, err = http.Post(url, constants.JSONContentType, r.Body)
	}

	// Only allow GET and POST requests
	if r.Method != constants.MethodGet && r.Method != constants.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil, err
	}

	// Return error when request fails
	if err != nil {
		http.Error(w, fmt.Sprintf("Error making request to endpoint: %v", err), http.StatusInternalServerError)
		return nil, err
	}

	return resp, nil
}
