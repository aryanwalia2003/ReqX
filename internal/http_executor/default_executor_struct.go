package http_executor

import "net/http"

// DefaultExecutor is the standard implementation of RequestExecutor using net/http.
type DefaultExecutor struct {
	client *http.Client
}
