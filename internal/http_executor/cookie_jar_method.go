package http_executor

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

// Enable turns cookie persistence back on.
func (j *ManagedCookieJar) Enable() { j.enabled = true }

// Disable turns cookie persistence off — cookies are neither sent nor stored.
func (j *ManagedCookieJar) Disable() { j.enabled = false }

// IsEnabled reports whether the jar is currently active.
func (j *ManagedCookieJar) IsEnabled() bool { return j.enabled }

// Clear discards all stored cookies by replacing the inner jar with a fresh one.
func (j *ManagedCookieJar) Clear() { j.inner, _ = cookiejar.New(nil) }

// SetCookies stores cookies — no-op when the jar is disabled.
func (j *ManagedCookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	if j.enabled {
		j.inner.SetCookies(u, cookies)
	}
}

// Cookies returns stored cookies — returns nil when the jar is disabled.
func (j *ManagedCookieJar) Cookies(u *url.URL) []*http.Cookie {
	if !j.enabled {
		return nil
	}
	return j.inner.Cookies(u)
}
