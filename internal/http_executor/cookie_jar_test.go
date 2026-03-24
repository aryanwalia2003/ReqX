package http_executor

import (
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"
)

// TestManagedCookieJar_ConcurrentAccess hammers every method of ManagedCookieJar
// from many goroutines simultaneously. Run with -race to verify no data races.
//
//	go test ./internal/http_executor/... -race -count=1
func TestManagedCookieJar_ConcurrentAccess(t *testing.T) {
	jar := NewManagedCookieJar()
	u, _ := url.Parse("http://example.com/")
	cookies := []*http.Cookie{{Name: "session", Value: "abc"}}

	const goroutines = 64
	const iters = 200

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iters; j++ {
				switch j % 5 {
				case 0:
					jar.SetCookies(u, cookies)
				case 1:
					_ = jar.Cookies(u)
				case 2:
					jar.Enable()
				case 3:
					jar.Disable()
				case 4:
					jar.Clear()
				}
			}
		}(i)
	}

	wg.Wait()
}

// TestManagedCookieJar_EnableDisable verifies that Enable/Disable gate
// cookie storage and retrieval correctly under sequential access.
func TestManagedCookieJar_EnableDisable(t *testing.T) {
	jar := NewManagedCookieJar()
	u, _ := url.Parse("http://example.com/")
	cookies := []*http.Cookie{
		{Name: "tok", Value: "xyz", Expires: time.Now().Add(time.Hour)},
	}

	// Enabled by default — cookies should be stored and returned.
	jar.SetCookies(u, cookies)
	got := jar.Cookies(u)
	if len(got) == 0 {
		t.Fatal("expected cookies when jar is enabled, got none")
	}

	// Disabled — SetCookies should be a no-op, Cookies should return nil.
	jar.Disable()
	if jar.IsEnabled() {
		t.Fatal("jar should be disabled")
	}
	jar.SetCookies(u, []*http.Cookie{{Name: "new", Value: "should_not_store"}})
	if got := jar.Cookies(u); got != nil {
		t.Fatalf("expected nil from disabled jar, got %v", got)
	}

	// Re-enable — previously stored cookie is still accessible.
	jar.Enable()
	if !jar.IsEnabled() {
		t.Fatal("jar should be enabled")
	}
	got = jar.Cookies(u)
	if len(got) == 0 {
		t.Fatal("expected cookie to survive disable/re-enable cycle")
	}
}

// TestManagedCookieJar_Clear verifies that Clear wipes all stored cookies.
func TestManagedCookieJar_Clear(t *testing.T) {
	jar := NewManagedCookieJar()
	u, _ := url.Parse("http://example.com/")
	cookies := []*http.Cookie{
		{Name: "tok", Value: "xyz", Expires: time.Now().Add(time.Hour)},
	}

	jar.SetCookies(u, cookies)
	if len(jar.Cookies(u)) == 0 {
		t.Fatal("expected cookies before clear")
	}

	jar.Clear()
	if len(jar.Cookies(u)) != 0 {
		t.Fatal("expected no cookies after clear")
	}
}

// TestManagedCookieJar_ClearWhileConcurrentReads verifies that Clear can
// replace the inner jar while other goroutines are reading from it without
// causing a panic or data race.
func TestManagedCookieJar_ClearWhileConcurrentReads(t *testing.T) {
	jar := NewManagedCookieJar()
	u, _ := url.Parse("http://example.com/")
	cookies := []*http.Cookie{
		{Name: "tok", Value: "xyz", Expires: time.Now().Add(time.Hour)},
	}
	jar.SetCookies(u, cookies)

	stop := make(chan struct{})
	var wg sync.WaitGroup

	// Continuous readers.
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					_ = jar.Cookies(u)
				}
			}
		}()
	}

	// Rapid clears.
	for i := 0; i < 50; i++ {
		jar.Clear()
		jar.SetCookies(u, cookies)
	}

	close(stop)
	wg.Wait()
}

// TestManagedCookieJar_ClearWhileConcurrentWrites verifies that SetCookies
// executing concurrently with Clear does not panic.
func TestManagedCookieJar_ClearWhileConcurrentWrites(t *testing.T) {
	jar := NewManagedCookieJar()
	u, _ := url.Parse("http://example.com/")
	cookies := []*http.Cookie{
		{Name: "tok", Value: "xyz", Expires: time.Now().Add(time.Hour)},
	}

	stop := make(chan struct{})
	var wg sync.WaitGroup

	// Continuous writers.
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					jar.SetCookies(u, cookies)
				}
			}
		}()
	}

	// Rapid clears.
	for i := 0; i < 50; i++ {
		jar.Clear()
		time.Sleep(time.Millisecond) // Yield slightly to let writers run
	}

	close(stop)
	wg.Wait()
}

// TestManagedCookieJar_URLFiltering verifies that the wrapper correctly
// delegates the URL constraints to the underlying stdlib jar.
func TestManagedCookieJar_URLFiltering(t *testing.T) {
	jar := NewManagedCookieJar()
	
	u1, _ := url.Parse("http://example.com/")
	u2, _ := url.Parse("http://other.com/")
	
	cookies1 := []*http.Cookie{{Name: "c1", Value: "v1"}}
	cookies2 := []*http.Cookie{{Name: "c2", Value: "v2"}}
	
	jar.SetCookies(u1, cookies1)
	jar.SetCookies(u2, cookies2)
	
	// Check example.com cookies
	got1 := jar.Cookies(u1)
	if len(got1) != 1 || got1[0].Name != "c1" {
		t.Fatalf("expected 1 cookie 'c1' for example.com, got %v", got1)
	}

	// Check other.com cookies
	got2 := jar.Cookies(u2)
	if len(got2) != 1 || got2[0].Name != "c2" {
		t.Fatalf("expected 1 cookie 'c2' for other.com, got %v", got2)
	}
}

// TestManagedCookieJar_ToggleContinuously tests rapidly enabling and disabling
// the cookie jar while reading from it.
func TestManagedCookieJar_ToggleContinuously(t *testing.T) {
	jar := NewManagedCookieJar()
	u, _ := url.Parse("http://example.com/")
	jar.SetCookies(u, []*http.Cookie{{Name: "c3", Value: "v3"}})

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			jar.Disable()
			jar.Enable()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			_ = jar.Cookies(u)
		}
	}()

	wg.Wait()
}