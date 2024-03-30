package structs

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

// GraphQLRequest struct based on the request body
type GraphQLRequest struct {
	Query    string `json:"query"`
	Mutation string `json:"mutation"`
}

func (g *GraphQLRequest) GetOperationNameHashed() (string, error) {
	// Returns query identifier
	if g.Query != "" {
		return "query_" + hashOperationName(g.Query), nil
	}

	// Returns mutation identifier
	if g.Mutation != "" {
		return "mutation_" + hashOperationName(g.Mutation), nil
	}

	return "", fmt.Errorf("no operation found")
}

func hashOperationName(operationName string) string {
	// Hash operation name using MD5
	hash := md5.Sum([]byte(operationName))

	// Convert hash to string
	hashedString := hex.EncodeToString(hash[:])

	return hashedString
}
