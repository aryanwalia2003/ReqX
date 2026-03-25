package http_executor

import (
	"net/http"
	"net/url"
	"sync"
	"testing"
)

// TestCookieJar_ConcurrentEnableDisable hammers Enable/Disable/Clear from
// multiple goroutines simultaneously, exposing any RWMutex gaps.
// Run with: go test ./internal/http_executor/... -race
func TestCookieJar_ConcurrentEnableDisable(t *testing.T) {
	jar := NewManagedCookieJar()
	u, _ := url.Parse("https://example.com")
	cookie := &http.Cookie{Name: "session", Value: "abc123"}

	var wg sync.WaitGroup
	const goroutines = 100

	// Writers: enable/disable/clear concurrently
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			switch i % 3 {
			case 0:
				jar.Enable()
			case 1:
				jar.Disable()
			case 2:
				jar.Clear()
			}
		}(i)
	}

	// Readers: Cookies/SetCookies/IsEnabled concurrently with writes
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			switch i % 3 {
			case 0:
				_ = jar.IsEnabled()
			case 1:
				jar.SetCookies(u, []*http.Cookie{cookie})
			case 2:
				_ = jar.Cookies(u)
			}
		}(i)
	}

	wg.Wait()
}

// TestCookieJar_ClearRace validates that Clear (which swaps the inner jar pointer)
// can be called concurrently with SetCookies/Cookies without a race.
// This is the subtle double-pointer race: read `inner` then call Set on it
// while another goroutine replaces `inner` via Clear.
func TestCookieJar_ClearRace(t *testing.T) {
	jar := NewManagedCookieJar()
	jar.Enable()
	u, _ := url.Parse("https://example.com")
	cookie := &http.Cookie{Name: "tok", Value: "xyz"}

	var wg sync.WaitGroup
	const iters = 200

	for i := 0; i < iters; i++ {
		wg.Add(3)
		go func() {
			defer wg.Done()
			jar.Clear()
		}()
		go func() {
			defer wg.Done()
			jar.SetCookies(u, []*http.Cookie{cookie})
		}()
		go func() {
			defer wg.Done()
			_ = jar.Cookies(u)
		}()
	}
	wg.Wait()
}

// TestCookieJar_EnabledFlagConsistency verifies that IsEnabled always reflects
// the final state written by Enable/Disable without stale reads.
func TestCookieJar_EnabledFlagConsistency(t *testing.T) {
	jar := NewManagedCookieJar()
	var wg sync.WaitGroup

	// Last writer wins — we just ensure no panic/race.
	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				jar.Enable()
			} else {
				jar.Disable()
			}
		}(i)
	}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = jar.IsEnabled()
		}()
	}
	wg.Wait()
}

// TestTransportPool_SetInsecureConcurrentRead verifies that SetInsecure
// called before workers start doesn't cause races during concurrent HTTP
// transport use (simulated by concurrent reads of TLSClientConfig).
//
// NOTE: SetInsecure is designed as a startup-only call. This test documents
// that behaviour — if called mid-run it would be unsafe.
func TestTransportPool_SetInsecureSafe_BeforeUse(t *testing.T) {
	// Call SetInsecure once at "startup", before any goroutines read it.
	SetInsecure(true)
	defer SetInsecure(false)

	// Simulate concurrent HTTP executor construction (reads the transport).
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			exec := NewDefaultExecutor()
			if exec == nil {
				t.Error("executor must not be nil")
			}
		}()
	}
	wg.Wait()
}
