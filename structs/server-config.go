package structs

import "fmt"

// Server config struct from server.config.json
type ServerConfig struct {
	Port int    `json:"port"`
	Path string `json:"path"`
}

func (s *ServerConfig) GetPort() string {
	return fmt.Sprintf(":%d", s.Port)
}
