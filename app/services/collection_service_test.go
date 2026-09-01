package services

import (
	"encoding/json"
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

func TestCollectionService_Save(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "new-collection.json")

	s := NewCollectionService(nil)
	coll := collection.Collection{
		Name: "New",
		Requests: []collection.Request{
			{Name: "Ping", Method: "GET", URL: "https://api.example.com/ping"},
		},
	}
	if err := s.Save(coll, path); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	got, err := s.Open(path)
	if err != nil {
		t.Fatalf("Open() after Save() error = %v", err)
	}
	if got.Name != "New" || len(got.Requests) != 1 || got.Requests[0].Name != "Ping" {
		t.Errorf("round-tripped collection = %+v, want New collection with 1 request named 'Ping'", got)
	}
}

func TestCollectionService_Save_InvalidPath(t *testing.T) {
	s := NewCollectionService(nil)
	err := s.Save(collection.Collection{Name: "X"}, filepath.Join(t.TempDir(), "missing-dir", "x.json"))
	if err == nil {
		t.Fatal("expected an error saving to a directory that doesn't exist")
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

func TestCollectionService_Run_RPSWithStages_Allowed(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s := NewCollectionService(nil)
	input := RunCollectionInput{
		Collection: collection.Collection{
			Requests: []collection.Request{{Name: "Ping", Method: "GET", URL: srv.URL}},
		},
		RPS:    10,
		Stages: "100ms:1",
	}

	// Stages terminate on their own (no ctx deadline needed), so RPS+Stages
	// without Duration must NOT be rejected the way RPS-alone is.
	if _, err := s.Run(input); err != nil {
		t.Fatalf("Run() error = %v, want nil (stages provide their own termination)", err)
	}
}

func TestCollectionService_Run_NoCookiesSequential(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/set" {
			http.SetCookie(w, &http.Cookie{Name: "sess", Value: "abc"})
			w.WriteHeader(http.StatusOK)
			return
		}
		// /check
		_, err := r.Cookie("sess")
		if err == nil {
			w.WriteHeader(http.StatusConflict) // cookie leaked through
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer srv.Close()

	s := NewCollectionService(nil)
	input := RunCollectionInput{
		Collection: collection.Collection{
			Requests: []collection.Request{
				{Name: "Set", Method: "GET", URL: srv.URL + "/set"},
				{Name: "Check", Method: "GET", URL: srv.URL + "/check"},
			},
		},
		NoCookies: true,
	}

	out, err := s.Run(input)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if len(out.Stats) != 2 || out.Stats[1].Failures != 0 {
		t.Errorf("Stats = %+v, want 'Check' to succeed (no cookie jar means nothing to leak)", out.Stats)
	}

	// A later call without NoCookies must NOT be permanently affected by the
	// one-off no-cookies runner built for the call above (sequentialRunner
	// must not have mutated the shared s.runner) — cookies should flow
	// normally again, so /check now DOES see the cookie /set left behind
	// (the handler reports that as a 409, by design — see above).
	out2, err := s.Run(RunCollectionInput{
		Collection: collection.Collection{
			Requests: []collection.Request{
				{Name: "Set", Method: "GET", URL: srv.URL + "/set"},
				{Name: "Check", Method: "GET", URL: srv.URL + "/check"},
			},
		},
	})
	if err != nil {
		t.Fatalf("second Run() error = %v", err)
	}
	if len(out2.Stats) != 2 || out2.Stats[1].Failures != 1 {
		t.Errorf("second run Stats = %+v, want 'Check' to fail with 409 (cookie jar enabled by default, so it sees the cookie)", out2.Stats)
	}
}

func TestCollectionService_Run_GraphQLErrorDetection(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // GraphQL: 200 OK even on a business-rule failure
		_, _ = w.Write([]byte(`{"errors":[{"message":"not authorized"}]}`))
	}))
	defer srv.Close()

	s := NewCollectionService(nil)

	withoutFlag, err := s.Run(RunCollectionInput{
		Collection: collection.Collection{
			Requests: []collection.Request{{Name: "Query", Method: "POST", URL: srv.URL}},
		},
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if withoutFlag.Summary.TotalFailures != 0 {
		t.Errorf("without GraphQL flag: TotalFailures = %d, want 0 (200 OK counts as success)", withoutFlag.Summary.TotalFailures)
	}

	withFlag, err := s.Run(RunCollectionInput{
		Collection: collection.Collection{
			Requests: []collection.Request{{Name: "Query", Method: "POST", URL: srv.URL}},
		},
		GraphQL: true,
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if withFlag.Summary.TotalFailures != 1 {
		t.Errorf("with GraphQL flag: TotalFailures = %d, want 1 (errors array detected)", withFlag.Summary.TotalFailures)
	}
}

func TestCollectionService_Run_Injection(t *testing.T) {
	var gotPaths []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPaths = append(gotPaths, r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s := NewCollectionService(nil)
	input := RunCollectionInput{
		Collection: collection.Collection{
			Requests: []collection.Request{{Name: "Existing", Method: "GET", URL: srv.URL + "/existing"}},
		},
		InjectIndex:  "1",
		InjectName:   "Injected",
		InjectMethod: "GET",
		InjectURL:    srv.URL + "/injected",
	}

	if _, err := s.Run(input); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if len(gotPaths) != 2 || gotPaths[0] != "/injected" || gotPaths[1] != "/existing" {
		t.Errorf("paths hit = %v, want [/injected, /existing] (injected at index 1)", gotPaths)
	}
}

func TestCollectionService_Run_ExportPath(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	exportPath := filepath.Join(t.TempDir(), "export.ndjson")
	s := NewCollectionService(nil)
	input := RunCollectionInput{
		Collection: collection.Collection{
			Requests: []collection.Request{{Name: "Ping", Method: "GET", URL: srv.URL}},
		},
		ExportPath: exportPath,
	}

	if _, err := s.Run(input); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	data, err := os.ReadFile(exportPath)
	if err != nil {
		t.Fatalf("export file not written: %v", err)
	}
	var rec map[string]any
	if err := json.Unmarshal(data, &rec); err != nil {
		t.Fatalf("export file isn't valid JSON: %v (content: %s)", err, data)
	}
	if rec["name"] != "Ping" {
		t.Errorf("exported record name = %v, want 'Ping'", rec["name"])
	}
}

func TestCollectionService_Run_DagNodes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s := NewCollectionService(nil)
	input := RunCollectionInput{
		Collection: collection.Collection{
			Requests: []collection.Request{
				{Name: "First", Method: "GET", URL: srv.URL + "/1"},
				{Name: "Second", Method: "GET", URL: srv.URL + "/2", DependsOn: []string{"First"}},
			},
		},
	}

	out, err := s.Run(input)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if len(out.DagNodes) != 2 {
		t.Fatalf("len(DagNodes) = %d, want 2 (collection has depends_on)", len(out.DagNodes))
	}
	byName := map[string]int{}
	for _, n := range out.DagNodes {
		byName[n.Name] = n.LevelIdx
	}
	if byName["First"] != 0 || byName["Second"] != 1 {
		t.Errorf("levels = %+v, want First=0, Second=1", byName)
	}
}

func TestCollectionService_Run_NoDagNodesForLinearCollection(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s := NewCollectionService(nil)
	out, err := s.Run(RunCollectionInput{
		Collection: collection.Collection{
			Requests: []collection.Request{{Name: "Ping", Method: "GET", URL: srv.URL}},
		},
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if len(out.DagNodes) != 0 {
		t.Errorf("DagNodes = %+v, want empty for a linear (no depends_on) collection", out.DagNodes)
	}
}

func TestCollectionService_PickSaveFile_BeforeStartup(t *testing.T) {
	appCtx = nil
	s := NewCollectionService(nil)
	if _, err := s.PickSaveFile("Export results", "results.ndjson"); err == nil {
		t.Fatal("expected an error calling PickSaveFile before OnStartup, got nil")
	}
}
