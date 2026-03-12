package http_executor

import "net/http/cookiejar"

// NewManagedCookieJar creates a cookie jar that is enabled by default.
func NewManagedCookieJar() *ManagedCookieJar {
	jar, _ := cookiejar.New(nil)
	return &ManagedCookieJar{inner: jar, enabled: true}
}
