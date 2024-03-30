package structs

// Host config struct from hosts.config.json
type HostConfig struct {
	Url           string `json:"url"`
	Enabled       bool   `json:"enabled"`
	UseMock       bool   `json:"use_mock"`
	EnableGraphql bool   `json:"enable_graphql"`
	EnableREST    bool   `json:"enable_rest"`
	GeneratePath  string `json:"generate_path"`
	GraphQLPath   string `json:"graphql_path"`
}
