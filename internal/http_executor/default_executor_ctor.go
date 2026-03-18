package http_executor

import (
	"net/http"
	"time"
)

// NewDefaultExecutor constructs a new HTTP executor with cookie jar enabled.
// All executors share globalTransport so TCP connections are reused across
// workers, eliminating the per-request handshake overhead at high concurrency.
func NewDefaultExecutor() *DefaultExecutor {
	jar := NewManagedCookieJar()
	return &DefaultExecutor{
		jar: jar,
		client: &http.Client{
			Timeout:   60 * time.Second,
			Jar:       jar,
			Transport: globalTransport,
		},
	}
}
