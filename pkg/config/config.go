package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type ServiceConfig struct {
	Enabled        bool   `json:"enabled"`
	Host           string `json:"host"`
	ResponseFolder string `json:"response_folder"`
	Store          bool   `json:"store"`
}

type Config map[string]ServiceConfig

var serviceConfig Config

func LoadConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&serviceConfig); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return err
	}

	return nil
}

func GetAllServicesConfig() Config {
	return serviceConfig
}

func GetServiceConfig(host string) (ServiceConfig, error) {
	service := serviceConfig[host]
	if service.Host == "" {
		return ServiceConfig{}, fmt.Errorf("service configuration not found")
	}

	return service, nil
}
