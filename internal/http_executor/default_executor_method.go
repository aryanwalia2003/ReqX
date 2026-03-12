package http_executor

import "net/http"

// Execute sends the HTTP request using the embedded standard library client.
func (e *DefaultExecutor) Execute(req *http.Request) (*http.Response, error) {
	return e.client.Do(req)
}

// EnableCookies re-enables cookie persistence for all subsequent requests.
func (e *DefaultExecutor) EnableCookies() { e.jar.Enable() }

// DisableCookies stops cookies from being sent or stored.
func (e *DefaultExecutor) DisableCookies() { e.jar.Disable() }

// ClearCookies discards all cookies currently stored in the jar.
func (e *DefaultExecutor) ClearCookies() { e.jar.Clear() }

