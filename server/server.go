package server

import (
	"net/http"

	handlers "github.com/thiagodsantos/gomockserver/server/handlers"
	"github.com/thiagodsantos/gomockserver/structs"
)

func Execute(w http.ResponseWriter, r *http.Request, config structs.HostConfig, url string, requestBody []byte) (*http.Response, error) {
	var err error
	var resp *http.Response

	// Make GraphQL request
	if config.EnableGraphql {
		resp, err = handlers.GraphqlHandler(w, r, url, requestBody)
	}

	// Make REST request
	if config.EnableREST {
		resp, err = handlers.RESTHandler(w, r, url, requestBody)
	}

	return resp, err
}
