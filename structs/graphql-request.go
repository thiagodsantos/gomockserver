package structs

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type GraphQLRequest struct {
	Query    string `json:"query"`
	Mutation string `json:"mutation"`
}

func (g *GraphQLRequest) GetOperationNameHashed() (string, error) {
	if g.Query != "" {
		return "query_" + compressOperaionName(g.Query), nil
	}

	if g.Mutation != "" {
		return "mutation_" + compressOperaionName(g.Mutation), nil
	}

	return "", fmt.Errorf("no operation found")
}

func compressOperaionName(operationName string) string {
	hash := sha256.Sum256([]byte(operationName))

	hashedString := hex.EncodeToString(hash[:])

	return hashedString
}
