package structs

import (
	"fmt"
	"net/url"
	"strings"
)

// Host config struct from hosts.config.json
type HostConfig struct {
	Url           string `json:"url"`
	Enabled       bool   `json:"enabled"`
	EnableMock    bool   `json:"enable_mock"`
	EnableGraphql bool   `json:"enable_graphql"`
	EnableREST    bool   `json:"enable_rest"`
	GeneratePath  string `json:"generate_path"`
	GraphQLPath   string `json:"graphql_path"`
}

func (h *HostConfig) GetHostURL(path string) (string, error) {
	var err error
	var hostURL string

	if h.EnableGraphql {
		hostURL, err = h.getGraphQLURL()
		if err != nil {
			return "", err
		}

		return hostURL, nil
	}

	if h.EnableREST {
		hostURL, err = h.getRESTURL(path)
		if err != nil {
			return "", err
		}

		return hostURL, nil
	}

	return "", fmt.Errorf("no enabled services")
}

func (h *HostConfig) getRESTURL(path string) (string, error) {
	url, err := url.Parse(h.Url + path)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}

func (h *HostConfig) getGraphQLURL() (string, error) {
	url, err := url.Parse(h.Url + h.GraphQLPath)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}

func (h *HostConfig) IsGeneratePath(path string) bool {
	return strings.HasPrefix(path, h.GeneratePath)
}
