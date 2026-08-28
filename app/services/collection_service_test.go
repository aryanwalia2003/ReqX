package services

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"reqx/internal/collection"
	"reqx/internal/personas"
)

func TestCollectionService_PickFile_BeforeStartup(t *testing.T) {
	// appCtx is package-level and unset in tests (no wails.Run/OnStartup) —
	// PickFile must fail cleanly instead of nil-pointer-panicking on
	// runtime.OpenFileDialog.
	appCtx = nil
	s := NewCollectionService(nil)
	if _, err := s.PickFile("Open collection", "json"); err == nil {
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

	s := NewCollectionService(nil)
	got, err := s.Open(path)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	if got.Name != "Demo" || len(got.Requests) != 1 || got.Requests[0].Name != "Get user" {
		t.Errorf("Open() = %+v, want Demo collection with 1 request named 'Get user'", got)
	}
}

func TestCollectionService_Open_MissingFile(t *testing.T) {
	s := NewCollectionService(nil)
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

	s := NewCollectionService(nil)
	got, err := s.OpenEnvironment(path)
	if err != nil {
		t.Fatalf("OpenEnvironment() error = %v", err)
	}
	if got.Name != "dev" || got.Variables["base_url"] != "https://api.example.com" {
		t.Errorf("OpenEnvironment() = %+v, want dev env with base_url set", got)
	}
}

func TestCollectionService_OpenEnvironment_MissingFile(t *testing.T) {
	s := NewCollectionService(nil)
	if _, err := s.OpenEnvironment(filepath.Join(t.TempDir(), "does-not-exist.json")); err == nil {
		t.Fatal("expected an error opening a missing file, got nil")
	}
}

func TestCollectionService_OpenPersonas(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "personas.csv")
	body := "name,role\nAlice,admin\nBob,viewer\n"
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	s := NewCollectionService(nil)
	got, err := s.OpenPersonas(path)
	if err != nil {
		t.Fatalf("OpenPersonas() error = %v", err)
	}
	if len(got) != 2 || got[0]["name"] != "Alice" || got[0]["role"] != "admin" || got[1]["name"] != "Bob" {
		t.Errorf("OpenPersonas() = %+v, want [Alice/admin, Bob/viewer]", got)
	}
}

func TestCollectionService_OpenPersonas_MissingFile(t *testing.T) {
	s := NewCollectionService(nil)
	if _, err := s.OpenPersonas(filepath.Join(t.TempDir(), "does-not-exist.csv")); err == nil {
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

	s := NewCollectionService(nil)
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

	if len(out.Stats) != 2 {
		t.Fatalf("len(Stats) = %d, want 2", len(out.Stats))
	}
	if out.Stats[0].Successes != 1 || out.Stats[0].Failures != 0 {
		t.Errorf("Stats[0] (Healthy) = %+v, want 1 success, 0 failures", out.Stats[0])
	}
	if out.Stats[1].Successes != 0 || out.Stats[1].Failures != 1 || out.Stats[1].TopError == "" {
		t.Errorf("Stats[1] (Broken) = %+v, want 0 successes, 1 failure, a TopError", out.Stats[1])
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

	s := NewCollectionService(nil)
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

func TestCollectionService_Run_WorkerPool(t *testing.T) {
	var hits atomic.Int64
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s := NewCollectionService(nil)
	input := RunCollectionInput{
		Collection: collection.Collection{
			Requests: []collection.Request{{Name: "Ping", Method: "GET", URL: srv.URL}},
		},
		Workers:    4,
		Iterations: 8,
	}

	out, err := s.Run(input)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if hits.Load() != 8 {
		t.Errorf("server saw %d hits, want 8 (workers=4 x iterations=8, one request each)", hits.Load())
	}
	if len(out.Stats) != 1 || out.Stats[0].TotalRuns != 8 || out.Stats[0].Successes != 8 {
		t.Errorf("Stats = %+v, want 1 entry with TotalRuns=8, Successes=8", out.Stats)
	}
	if out.Summary.TotalRequests != 8 {
		t.Errorf("Summary.TotalRequests = %d, want 8", out.Summary.TotalRequests)
	}
}

func TestCollectionService_Run_Scheduler_Duration(t *testing.T) {
	var hits atomic.Int64
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s := NewCollectionService(nil)
	input := RunCollectionInput{
		Collection: collection.Collection{
			Requests: []collection.Request{{Name: "Ping", Method: "GET", URL: srv.URL}},
		},
		Workers:    2,
		DurationMs: 150,
	}

	start := time.Now()
	out, err := s.Run(input)
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if elapsed > time.Second {
		t.Errorf("Run() took %v, want roughly the 150ms duration (scheduler may be hanging)", elapsed)
	}
	if hits.Load() == 0 {
		t.Error("server saw 0 hits during a 150ms duration run with 2 workers")
	}
	if out.Summary.TotalRequests == 0 {
		t.Errorf("Summary.TotalRequests = 0, want > 0")
	}
}

func TestCollectionService_Run_SubstitutesPersonaSequential(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s := NewCollectionService(nil)
	input := RunCollectionInput{
		Collection: collection.Collection{
			Requests: []collection.Request{{Name: "Get", Method: "GET", URL: srv.URL + "/{{persona.name}}"}},
		},
		Personas: []personas.Persona{{"name": "Alice"}, {"name": "Bob"}},
	}

	// Sequential path always uses Personas[0], matching the CLI's loop.
	if _, err := s.Run(input); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if gotPath != "/Alice" {
		t.Errorf("request path = %q, want /Alice ({{persona.name}} from Personas[0])", gotPath)
	}
}

func TestCollectionService_Run_PersonasRoundRobinAcrossWorkers(t *testing.T) {
	seen := make(chan string, 4)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen <- r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s := NewCollectionService(nil)
	input := RunCollectionInput{
		Collection: collection.Collection{
			Requests: []collection.Request{{Name: "Get", Method: "GET", URL: srv.URL + "/{{persona.name}}"}},
		},
		Workers:    2,
		Iterations: 4,
		Personas:   []personas.Persona{{"name": "Alice"}, {"name": "Bob"}},
	}

	if _, err := s.Run(input); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	close(seen)

	names := map[string]bool{}
	for p := range seen {
		names[p] = true
	}
	if !names["/Alice"] || !names["/Bob"] {
		t.Errorf("paths seen = %v, want both /Alice and /Bob (round-robin across 2 workers)", names)
	}
}

func TestCollectionService_Run_RPSWithoutDuration_Rejected(t *testing.T) {
	s := NewCollectionService(nil)
	input := RunCollectionInput{
		Collection: collection.Collection{
			Requests: []collection.Request{{Name: "Ping", Method: "GET", URL: "http://127.0.0.1:1"}},
		},
		RPS: 10,
	}

	if _, err := s.Run(input); err == nil {
		t.Fatal("expected an error for rps without a duration (would hang the Scheduler forever), got nil")
	}
}
