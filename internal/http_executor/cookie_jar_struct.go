package http_executor

import "net/http/cookiejar"

// ManagedCookieJar wraps the stdlib cookie jar with enable/disable/clear support.
// When disabled, cookies are neither sent nor stored — the jar behaves as if empty.
type ManagedCookieJar struct {
	inner   *cookiejar.Jar
	enabled bool
}
