package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/thiagodsantos/gomockserver/constants"
	"github.com/thiagodsantos/gomockserver/structs"
)

func GetGraphQLRequestBody(r *http.Request, reqBody []byte) (*structs.GraphQLRequest, error) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		return nil, fmt.Errorf("method not allowed")
	}

	// Decode request body
	var requestBody structs.GraphQLRequest
	err := json.Unmarshal(reqBody, &requestBody)

	// Return error when decoding request body fails
	if err != nil {
		return nil, fmt.Errorf("error decoding request body: %v", err)
	}

	return &requestBody, err
}

func GetOperationName(w http.ResponseWriter, r *http.Request, reqBody []byte) (string, error) {
	requestBody, err := GetGraphQLRequestBody(r, reqBody)
	if err != nil {
		return "", err
	}

	operationName, err := requestBody.GetOperationNameHashed()
	if err != nil {
		return "", err
	}

	return operationName, nil
}

func GraphqlHandler(r *http.Request, graphqlUrl string, reqBody []byte) (*http.Response, error) {
	// Get GraphQL request body
	requestBody, err := GetGraphQLRequestBody(r, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error getting graphql request")
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

	// Return error when encoding request body fails
	if err != nil {
		return nil, fmt.Errorf("error encoding request body: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest(constants.MethodPost, graphqlUrl, bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set(constants.HeaderContentType, constants.JSONContentType)

	// Make GraphQL request
	var resp *http.Response
	var client = &http.Client{}

	// Make graphql request
	resp, err = client.Do(req)

	// Return error when request fails
	if err != nil {
		return nil, fmt.Errorf("error making graphql request: %v", err)
	}

	return resp, err
}
