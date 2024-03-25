package config

import (
	"fmt"

	"github.com/thiagodsantos/gomockserver/constants"
	"github.com/thiagodsantos/gomockserver/structs"
	"github.com/thiagodsantos/gomockserver/utils"
)

// Get host config from hosts.config.json
func GetHostConfig() (structs.Config, error) {
	var configs []structs.Config

	// Read JSON file data from hosts.config.json
	_, err := utils.ReadJSONFile(constants.HostsConfigFileName, &configs)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return structs.Config{}, err
	}

	// Get host config enabled
	var hostCount int
	var hostConfig structs.Config
	for _, config := range configs {
		if config.Enabled {
			hostConfig = config
			hostCount++
		}
	}

	// Return error if more than one host is enabled
	if hostCount > 1 {
		return structs.Config{}, fmt.Errorf("more than one host enabled in hosts.config")
	}

	return hostConfig, nil
}
