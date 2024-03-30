package config

import (
	"fmt"

	"github.com/thiagodsantos/gomockserver/constants"
	"github.com/thiagodsantos/gomockserver/structs"
	"github.com/thiagodsantos/gomockserver/utils"
)

func GetServerConfig() (structs.ServerConfig, error) {
	var configs structs.ServerConfig

	// Check if server config file exists
	if !utils.FileExists(constants.ServerConfigFileName) {
		fmt.Printf("file %s not found", constants.ServerConfigFileName)
		panic(1)
	}

	// Read server config file
	_, err := utils.ReadJSONFile(constants.ServerConfigFileName, &configs)
	if err != nil {
		fmt.Printf("Error reading JSON file %s", constants.ServerConfigFileName)
		panic(err)
	}

	// Set default values
	if configs.Path == "" {
		configs.Path = constants.ProxyPath
	}

	if configs.Port == 0 {
		configs.Port = constants.ProxyServerPort
	}

	return configs, nil
}
