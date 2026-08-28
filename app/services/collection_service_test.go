package services

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"reqx/internal/collection"
)

func TestCollectionService_PickFile_BeforeStartup(t *testing.T) {
	// appCtx is package-level and unset in tests (no wails.Run/OnStartup) —
	// PickFile must fail cleanly instead of nil-pointer-panicking on
	// runtime.OpenFileDialog.
	appCtx = nil
	s := NewCollectionService()
	if _, err := s.PickFile("Open collection"); err == nil {
		t.Fatal("expected an error calling PickFile before OnStartup, got nil")
	}
}

func TestCollectionService_Open(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "collection.json")
	body := `{"name":"Demo","requests":[{"name":"Get user","method":"GET","url":"https://api.example.com/user"}]}`
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	s := NewCollectionService()
	got, err := s.Open(path)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	if got.Name != "Demo" || len(got.Requests) != 1 || got.Requests[0].Name != "Get user" {
		t.Errorf("Open() = %+v, want Demo collection with 1 request named 'Get user'", got)
	}
}

func TestCollectionService_Open_MissingFile(t *testing.T) {
	s := NewCollectionService()
	if _, err := s.Open(filepath.Join(t.TempDir(), "does-not-exist.json")); err == nil {
		t.Fatal("expected an error opening a missing file, got nil")
	}
}

func TestCollectionService_OpenEnvironment(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "env.json")
	body := `{"name":"dev","variables":{"base_url":"https://api.example.com"}}`
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	s := NewCollectionService()
	got, err := s.OpenEnvironment(path)
	if err != nil {
		t.Fatalf("OpenEnvironment() error = %v", err)
	}
	if got.Name != "dev" || got.Variables["base_url"] != "https://api.example.com" {
		t.Errorf("OpenEnvironment() = %+v, want dev env with base_url set", got)
	}
}

func TestCollectionService_OpenEnvironment_MissingFile(t *testing.T) {
	s := NewCollectionService()
	if _, err := s.OpenEnvironment(filepath.Join(t.TempDir(), "does-not-exist.json")); err == nil {
		t.Fatal("expected an error opening a missing file, got nil")
	}
}

func TestCollectionService_Run(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ok":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"ok":true}`))
		case "/broken":
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer srv.Close()

	s := NewCollectionService()
	input := RunCollectionInput{
		Collection: collection.Collection{
			Name: "Demo",
			Requests: []collection.Request{
				{Name: "Healthy", Method: "GET", URL: srv.URL + "/ok"},
				{Name: "Broken", Method: "GET", URL: srv.URL + "/broken"},
			},
		},
	}

	out, err := s.Run(input)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if len(out.Results) != 2 {
		t.Fatalf("len(Results) = %d, want 2", len(out.Results))
	}
	if out.Results[0].StatusCode != 200 || out.Results[0].ErrorMessage != "" {
		t.Errorf("Results[0] = %+v, want a clean 200", out.Results[0])
	}
	if out.Results[1].StatusCode != 500 || out.Results[1].ErrorMessage == "" {
		t.Errorf("Results[1] = %+v, want a 500 with an error message", out.Results[1])
	}

	if out.Summary.TotalRequests != 2 || out.Summary.TotalSuccess != 1 || out.Summary.TotalFailures != 1 {
		t.Errorf("Summary = %+v, want 2 total, 1 success, 1 failure", out.Summary)
	}
}

func TestCollectionService_Run_SubstitutesEnvVariables(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s := NewCollectionService()
	input := RunCollectionInput{
		Collection: collection.Collection{
			Requests: []collection.Request{{Name: "Get", Method: "GET", URL: srv.URL + "/{{id}}"}},
		},
		EnvVariables: map[string]string{"id": "42"},
	}

	if _, err := s.Run(input); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if gotPath != "/42" {
		t.Errorf("request path = %q, want /42 ({{id}} substituted)", gotPath)
	}
}
