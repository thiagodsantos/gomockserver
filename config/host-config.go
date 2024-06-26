package config

import (
	"fmt"

	"github.com/thiagodsantos/gomockserver/constants"
	"github.com/thiagodsantos/gomockserver/structs"
	"github.com/thiagodsantos/gomockserver/utils"
)

// Get host config from hosts.config.json
func GetHostConfig() (structs.HostConfig, error) {
	var configs []structs.HostConfig

	// Check if hosts.config.json exists
	if !utils.FileExists(constants.HostsConfigFileName) {
		return structs.HostConfig{}, fmt.Errorf("file %s not found", constants.HostsConfigFileName)
	}

	// Read JSON file data from hosts.config.json
	_, err := utils.ReadJSONFile(constants.HostsConfigFileName, &configs)
	if err != nil {
		return structs.HostConfig{}, fmt.Errorf("error reading JSON file %s", constants.HostsConfigFileName)
	}

	// Get host config enabled
	var hostCount int
	var hostConfig structs.HostConfig
	for _, config := range configs {
		if config.Enabled {
			hostConfig = config
			hostCount++
		}
	}

	// Return error if more than one host is enabled
	if hostCount > 1 {
		return structs.HostConfig{}, fmt.Errorf("more than one host enabled in hosts.config")
	}

	// Return error if no host is enabled
	if hostCount == 0 {
		return structs.HostConfig{}, fmt.Errorf("no host enabled in hosts.config")
	}

	// Set default values
	if hostConfig.GeneratePath == "" {
		hostConfig.GeneratePath = constants.GeneratePath
	}

	return hostConfig, nil
}
