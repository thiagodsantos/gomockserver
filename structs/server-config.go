package structs

// Server config struct from server.config.json
type ServerConfig struct {
	Port int    `json:"port"`
	Path string `json:"path"`
}
