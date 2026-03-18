package http_executor

import (
	"net/http"
	"time"
)

// globalTransport is a shared connection pool used by all workers.
// Separating transport from the per-worker client lets TCP connections
// (and TLS sessions) be reused across iterations, eliminating the
// per-request handshake overhead that kills performance at high concurrency.
var globalTransport = &http.Transport{
	MaxIdleConns:        10000,
	MaxIdleConnsPerHost: 2000,
	IdleConnTimeout:     90 * time.Second,
	DisableCompression:  false,
	ForceAttemptHTTP2:   true,
}
