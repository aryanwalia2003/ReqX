package http_executor

import "net/http"

// Execute sends the HTTP request using the embedded standard library client.
func (e *DefaultExecutor) Execute(req *http.Request) (*http.Response, error) {
	return e.client.Do(req)
}
