package http_executor

import "net/http"

// RequestExecutor defines the interface for executing HTTP requests.
type RequestExecutor interface {
	Execute(req *http.Request) (*http.Response, error) //execute method will take a pointer to the http request and will return a pointer to the http response and an error
}
