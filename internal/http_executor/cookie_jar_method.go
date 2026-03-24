package http_executor

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

// Enable turns cookie persistence back on.
func (j *ManagedCookieJar) Enable() {
	j.mu.Lock()
	j.enabled = true
	j.mu.Unlock()
}

// Disable turns cookie persistence off — cookies are neither sent nor stored.
func (j *ManagedCookieJar) Disable() {
	j.mu.Lock()
	j.enabled = false
	j.mu.Unlock()
}

// IsEnabled reports whether the jar is currently active.
func (j *ManagedCookieJar) IsEnabled() bool {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.enabled
}

// Clear discards all stored cookies by replacing the inner jar with a fresh
// one. Both the pointer swap and the enabled check are done under the write
// lock so no goroutine can observe a partially-replaced jar.
func (j *ManagedCookieJar) Clear() {
	j.mu.Lock()
	j.inner, _ = cookiejar.New(nil)
	j.mu.Unlock()
}

// SetCookies stores cookies — no-op when the jar is disabled.
// The enabled check and the inner.SetCookies call are separated by design:
// the stdlib jar is internally locked, so we only need our lock to safely
// read the enabled flag and capture the inner pointer before releasing it.
func (j *ManagedCookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	j.mu.RLock()
	enabled := j.enabled
	inner := j.inner
	j.mu.RUnlock()

	if enabled {
		inner.SetCookies(u, cookies)
	}
}

// Cookies returns stored cookies — returns nil when the jar is disabled.
func (j *ManagedCookieJar) Cookies(u *url.URL) []*http.Cookie {
	j.mu.RLock()
	enabled := j.enabled
	inner := j.inner
	j.mu.RUnlock()

	if !enabled {
		return nil
	}
	return inner.Cookies(u)
}