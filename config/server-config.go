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
		return structs.ServerConfig{}, fmt.Errorf("server config file %s does not exist", constants.ServerConfigFileName)
	}

	// Read server config file
	_, err := utils.ReadJSONFile(constants.ServerConfigFileName, &configs)
	if err != nil {
		return structs.ServerConfig{}, fmt.Errorf("error reading server config file: %v", err)
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
