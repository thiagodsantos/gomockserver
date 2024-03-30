package config

import (
	"fmt"

	"github.com/thiagodsantos/gomockserver/storage"
)

func Setup() error {
	// Generate empty server config file
	err := storage.GenerateEmptyServerConfigFile()
	if err != nil {
		return fmt.Errorf("error generating empty server config file: %v", err)
	}

	// Generate empty hosts config file
	err = storage.GenerateEmptyHostsConfigFile()
	if err != nil {
		return fmt.Errorf("error generating empty hosts config file: %v", err)
	}

	return nil
}
