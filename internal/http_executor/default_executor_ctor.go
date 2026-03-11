package http_executor

import (
	"net/http"
	"time"
)

// NewDefaultExecutor constructs a new HTTP executor.
func NewDefaultExecutor() RequestExecutor {
	return &DefaultExecutor{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}
