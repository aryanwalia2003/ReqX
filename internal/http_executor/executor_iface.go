package http_executor

import "net/http"

// RequestExecutor defines the interface for executing HTTP requests.
type RequestExecutor interface {
	Execute(req *http.Request) (*http.Response, error)
	EnableCookies()
	DisableCookies()
	ClearCookies()
}
