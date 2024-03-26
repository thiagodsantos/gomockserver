package config

import (
	"fmt"

	"github.com/thiagodsantos/gomockserver/constants"
	"github.com/thiagodsantos/gomockserver/structs"
	"github.com/thiagodsantos/gomockserver/utils"
)

func GetServerConfig() (structs.ServerConfig, error) {
	var configs structs.ServerConfig

	if !utils.FileExists(constants.ServerConfigFileName) {
		return structs.ServerConfig{}, fmt.Errorf("file %s not found", constants.ServerConfigFileName)
	}

	_, err := utils.ReadJSONFile(constants.ServerConfigFileName, &configs)
	if err != nil {
		fmt.Printf("Error reading JSON file %s", constants.ServerConfigFileName)
		panic(err)
	}

	if configs.Path == "" {
		configs.Path = "/"
	}

	if configs.Port == 0 {
		configs.Port = constants.ProxyServerPort
	}

	return configs, nil
}
