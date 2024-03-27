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

	serverConfig := structs.ServerConfig{
		Port: constants.ProxyServerPort,
		Path: "/",
	}

	return utils.SaveJSONFile(constants.ServerConfigFileName, serverConfig)
}

func GenerateEmptyHostsConfigFile() error {
	// Check if hosts.config.json exists
	if utils.FileExists(constants.HostsConfigFileName) {
		return nil
	}

	// generate empty HostConfig array
	hostsConfig := []structs.HostConfig{}
	hostsConfig = append(hostsConfig, structs.HostConfig{
		Host:         "",
		Enabled:      false,
		UseMock:      false,
		GeneratePath: constants.DefaultGeneratePath,
	})

	return utils.SaveJSONFile(constants.HostsConfigFileName, hostsConfig)
}
