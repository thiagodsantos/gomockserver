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

	// Read JSON file data from hosts.config.json
	_, err := utils.ReadJSONFile(constants.HostsConfigFileName, &configs)
	if err != nil {
		fmt.Printf("Error reading JSON file %s", constants.HostsConfigFileName)
		return structs.HostConfig{}, err
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

	return hostConfig, nil
}
