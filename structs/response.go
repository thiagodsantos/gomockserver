package structs

// Response sSruct to save response to file in JSON format with filename response_<url>.json
type Response struct {
	URL          string              `json:"url"`
	Method       string              `json:"method"`
	Headers      map[string][]string `json:"headers"`
	Body         interface{}         `json:"body,omitempty"`
	StatusCode   int                 `json:"status_code"`
	ResponseTime string              `json:"response_time"`
}
