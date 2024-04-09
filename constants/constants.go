package constants

// Server constants
const (
	ProxyPath       = "/"
	ProxyServerPort = 8080
	GeneratePath    = "/generate"
)

// HTTP constants
const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodDelete = "DELETE"
)

var HTTPMethods = map[string]bool{
	MethodGet:    true,
	MethodPost:   true,
	MethodPut:    true,
	MethodDelete: true,
}

// Header constants
const (
	HeaderContentType = "Content-Type"
)

// JSON constants
const (
	JSONExtension   = "json"
	JSONContentType = "application/json"
)

// File name constants
const (
	HostsConfigFileName  = "hosts.config.json"
	ServerConfigFileName = "server.config.json"
	ResponseFileName     = "response"
	RequestFileName      = "request"
)

// Folder name constants
const (
	OutputFolder = ".output"
)
