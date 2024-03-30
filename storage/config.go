package storage

import (
	"github.com/thiagodsantos/gomockserver/constants"
	"github.com/thiagodsantos/gomockserver/structs"
	"github.com/thiagodsantos/gomockserver/utils"
)

func GenerateEmptyServerConfigFile() error {
	// Check if server.config.json exists
	if utils.FileExists(constants.ServerConfigFileName) {
		return nil
	}

	// Generate empty ServerConfig struct
	serverConfig := structs.ServerConfig{
		Port: constants.ProxyServerPort,
		Path: constants.ProxyPath,
	}

	// Save ServerConfig struct to server.config.json
	return utils.SaveJSONFile(constants.ServerConfigFileName, serverConfig)
}

func GenerateEmptyHostsConfigFile() error {
	// Check if hosts.config.json exists
	if utils.FileExists(constants.HostsConfigFileName) {
		return nil
	}

	// Generate empty HostConfig array
	hostsConfig := []structs.HostConfig{}
	hostsConfig = append(hostsConfig, structs.HostConfig{
		Url:          "",
		Enabled:      false,
		EnableMock:   false,
		GeneratePath: constants.GeneratePath,
	})

	// Save HostConfig array to hosts.config.json
	return utils.SaveJSONFile(constants.HostsConfigFileName, hostsConfig)
}
