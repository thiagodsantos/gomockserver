package structs

// Host config struct from hosts.config.json
type HostConfig struct {
	Host         string `json:"host"`
	Enabled      bool   `json:"enabled"`
	UseMock      bool   `json:"use_mock"`
	GeneratePath string `json:"generate_path"`
}
