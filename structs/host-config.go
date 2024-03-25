package structs

// HostConfig struct from hosts.config.json
type HostConfig struct {
	Host    string   `json:"host"`
	Paths   []string `json:"paths"`
	Enabled bool     `json:"enabled"`
	UseMock bool     `json:"use_mock"`
}
