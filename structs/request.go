package structs

// RequestInfo struct to save request info to file in JSON format with filename request_<url>.json
type Request struct {
	URL          string              `json:"url"`
	Method       string              `json:"method"`
	Headers      map[string][]string `json:"headers"`
	Body         interface{}         `json:"body,omitempty"`
	ResponseTime string              `json:"response_time"`
}
