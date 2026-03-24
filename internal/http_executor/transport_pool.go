package http_executor

import (
	"net"
	"net/http"
	"time"
)

// globalTransport is a shared connection pool used by all workers.
// Separating transport from the per-worker client lets TCP connections
// (and TLS sessions) be reused across iterations, eliminating the
// per-request handshake overhead that kills performance at high concurrency.
var globalTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	ForceAttemptHTTP2:     false, // Disabled to prevent MAX_CONCURRENT_STREAMS exhaustion at ALB
	MaxIdleConns:          10000,
	MaxIdleConnsPerHost:   10000,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}
