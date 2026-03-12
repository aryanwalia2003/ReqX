package scripting

import (
	"encoding/json"
)

// ResponseAPI exposes the HTTP or Socket.IO response to the JS VM.
type ResponseAPI struct {
	BodyString string           `json:"body"`
	HeadersMap map[string]string `json:"headersMap"`
	Headers    *ResponseHeaders `json:"headers"`
}

// json attempts to parse the response body as JSON.
// It returns a generic map/interface that Goja automatically translates to JS Objects.
func (r *ResponseAPI) Json() (interface{}, error) {
	var result interface{}
	err := json.Unmarshal([]byte(r.BodyString), &result) //yeh response ki body ko parse karta hai
	return result, err
}

// text returns the raw response body as a string.
func (r *ResponseAPI) Text() string {
	return r.BodyString
}

// ResponseHeaders provides the `pm.response.headers.*` capability.
type ResponseHeaders struct {
	Headers map[string]string
}

// get retrieves a specific header by key.
func (h *ResponseHeaders) Get(key string) string {
	if h.Headers == nil {
		return ""
	}
	return h.Headers[key]
}
