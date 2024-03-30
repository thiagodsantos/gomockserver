package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/thiagodsantos/gomockserver/constants"
	"github.com/thiagodsantos/gomockserver/structs"
)

func GetGraphQLRequestBody(w http.ResponseWriter, r *http.Request, reqBody []byte) (*structs.GraphQLRequest, error) {
	var err error

	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil, fmt.Errorf("method not allowed")
	}

	// Decode request body
	var requestBody structs.GraphQLRequest
	if err = json.Unmarshal(reqBody, &requestBody); err != nil {
		http.Error(w, "Error decoding request body", http.StatusInternalServerError)
		return nil, err
	}

	return &requestBody, err
}

func GetOperationName(w http.ResponseWriter, r *http.Request, reqBody []byte) (string, error) {
	requestBody, err := GetGraphQLRequestBody(w, r, reqBody)
	if err != nil {
		return "", err
	}

	operationName, err := requestBody.GetOperationNameHashed()
	if err != nil {
		return "", err
	}

	return operationName, nil
}

func GraphqlHandler(w http.ResponseWriter, r *http.Request, graphqlUrl string, reqBody []byte) (*http.Response, error) {
	var err error

	// Get GraphQL request body
	requestBody, err := GetGraphQLRequestBody(w, r, reqBody)
	if err != nil {
		http.Error(w, "Error getting graphql request", http.StatusInternalServerError)
		return nil, err
	}

	var reqJSON []byte

	// Create query request
	if requestBody.Query != "" {
		reqJSON, err = json.Marshal(map[string]string{
			"query": requestBody.Query,
		})
	}

	// Create mutation request
	if requestBody.Mutation != "" {
		reqJSON, err = json.Marshal(map[string]string{
			"mutation": requestBody.Mutation,
		})
	}

	if err != nil {
		fmt.Println("Error encoding request body:", err)
		return nil, err
	}

	var req *http.Request

	// Create HTTP request
	req, err = http.NewRequest(constants.MethodPost, graphqlUrl, bytes.NewBuffer(reqJSON))
	if err != nil {
		fmt.Println("Error creating graphql request:", err)
		return nil, err
	}
	req.Header.Set(constants.HeaderContentType, constants.JSONContentType)

	var resp *http.Response

	// Make GraphQL request
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Error making graphql request:", err)
		return nil, err
	}

	return resp, err
}
