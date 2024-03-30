package main

import (
	"fmt"
	"net/http"

	"github.com/thiagodsantos/gomockserver/config"
	"github.com/thiagodsantos/gomockserver/server"
	"github.com/thiagodsantos/gomockserver/utils"
)

func main() {
	// Setup
	err := config.Setup()
	if err != nil {
		fmt.Println("Error setting up config:", err)
		return
	}

	// Get server config
	serverConfig, err := config.GetServerConfig()
	if err != nil {
		fmt.Println("Error getting server config:", err)
		return
	}

	// Get server port
	port := serverConfig.GetPort()
	if port == "" {
		fmt.Println("Error getting server port")
		return
	}

	// Add handler to server
	http.HandleFunc(serverConfig.Path, server.Handler)

	// Print message to console with server port
	fmt.Println(utils.Format(utils.PURPLE, "Mock server starting on port "+port))

	// Start server
	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
}
