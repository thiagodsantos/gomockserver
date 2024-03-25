package structs

// Config struct from hosts.config.json
type Config struct {
	Host    string   `json:"host"`
	Paths   []string `json:"paths"`
	Enabled bool     `json:"enabled"`
	UseMock bool     `json:"use_mock"`
}
