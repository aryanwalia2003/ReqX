package runner

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"reqx/internal/collection"
)

func TestExecuteSingle(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer secret-token" {
			t.Errorf("Authorization header = %q, want %q", got, "Bearer secret-token")
		}
		if got := r.Header.Get("X-Custom"); got != "hello" {
			t.Errorf("X-Custom header = %q, want %q", got, "hello")
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) != `{"id":42}` {
			t.Errorf("body = %q, want %q", body, `{"id":42}`)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	cr := NewCollectionRunner(nil, nil, nil, nil)
	cr.SetVerbosity(VerbosityQuiet)

	ctx := NewRuntimeContext()
	ctx.Environment.Set("id", "42")

	req := collection.Request{
		Method:  "POST",
		URL:     srv.URL + "/items",
		Headers: map[string]string{"X-Custom": "hello"},
		Body:    `{"id":{{id}}}`,
		Auth:    &collection.Auth{Type: "bearer", Token: "secret-token"},
	}

	result, err := cr.ExecuteSingle(req, nil, ctx)
	if err != nil {
		t.Fatalf("ExecuteSingle() error = %v", err)
	}

	if result.StatusCode != http.StatusCreated {
		t.Errorf("StatusCode = %d, want %d", result.StatusCode, http.StatusCreated)
	}
	if result.Body != `{"ok":true}` {
		t.Errorf("Body = %q, want %q", result.Body, `{"ok":true}`)
	}
	if result.Headers["Content-Type"] != "application/json" {
		t.Errorf("Content-Type header = %q, want %q", result.Headers["Content-Type"], "application/json")
	}
	if result.BytesReceived != int64(len(`{"ok":true}`)) {
		t.Errorf("BytesReceived = %d, want %d", result.BytesReceived, len(`{"ok":true}`))
	}
}

func TestExecuteSingle_NetworkError(t *testing.T) {
	cr := NewCollectionRunner(nil, nil, nil, nil)
	ctx := NewRuntimeContext()

	req := collection.Request{Method: "GET", URL: "http://127.0.0.1:1"}

	_, err := cr.ExecuteSingle(req, nil, ctx)
	if err == nil {
		t.Fatal("expected an error dialing a closed port, got nil")
	}
}
