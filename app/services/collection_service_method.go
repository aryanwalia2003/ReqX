package services

import (
	"log"
	"strings"
	"time"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"reqx/internal/collection"
	"reqx/internal/environment"
	"reqx/internal/errs"
	"reqx/internal/history"
	"reqx/internal/http_executor"
	"reqx/internal/metrics"
	"reqx/internal/personas"
	"reqx/internal/planner"
	"reqx/internal/runner"
	"reqx/internal/storage"
)

// RunCollectionInput carries the collection to run (as returned by Open,
// round-tripped unedited for now), env vars for {{var}} substitution, and
// the same load-testing knobs as the CLI's `reqx run` flags:
//   - Workers (-c): concurrency. <=1 runs single-threaded.
//   - Iterations (-n): how many times to run the collection. Ignored once
//     Duration or RPS is set — Scheduler mode runs by wall-clock time, not
//     a fixed count, exactly like the CLI.
//   - DurationMs (-d) / RPS (--rps): switch into Scheduler mode (duration-
//     and/or rate-limited load test) once either is > 0.
//   - Personas (--personas): CSV rows (as returned by OpenPersonas),
//     round-tripped back unedited. Each column becomes {{persona.<col>}}.
//     Scheduler/WorkerPool assign one persona per worker, round-robin by
//     worker ID; the sequential path always uses Personas[0], same as the
//     CLI's sequential loop (cmd/run_cmd_ctor.go).
//   - Stages (--stages): a ramp plan string, e.g. "10s:5,30s:20,10s:0" —
//     parsed via runner.ParseStages, same syntax the CLI accepts. Setting
//     it also switches into Scheduler mode, same as Duration/RPS.
//   - NoCookies/ClearCookies/GraphQL: the CLI's --no-cookies/--clear-cookies/
//     --graphql flags.
//   - ExportPath: the CLI's --export — writes raw metrics as
//     newline-delimited JSON there after the run (best-effort; a failure
//     is logged, not returned as a run error, same as the CLI).
//   - InjectIndex/Name/Method/URL/Body/Headers: the CLI's --inject-* flags
//     — a temporary request inserted into the plan without touching the
//     saved collection. InjectHeaders is "Key: Value" strings, same as
//     --inject-header.
type RunCollectionInput struct {
	Collection    collection.Collection `json:"collection"`
	EnvVariables  map[string]string     `json:"envVariables,omitempty"`
	Workers       int                   `json:"workers,omitempty"`
	Iterations    int                   `json:"iterations,omitempty"`
	DurationMs    int64                 `json:"durationMs,omitempty"`
	RPS           float64               `json:"rps,omitempty"`
	Personas      []personas.Persona    `json:"personas,omitempty"`
	Stages        string                `json:"stages,omitempty"`
	NoCookies     bool                  `json:"noCookies,omitempty"`
	ClearCookies  bool                  `json:"clearCookies,omitempty"`
	GraphQL       bool                  `json:"graphql,omitempty"`
	ExportPath    string                `json:"exportPath,omitempty"`
	InjectIndex   string                `json:"injectIndex,omitempty"`
	InjectName    string                `json:"injectName,omitempty"`
	InjectMethod  string                `json:"injectMethod,omitempty"`
	InjectURL     string                `json:"injectUrl,omitempty"`
	InjectBody    string                `json:"injectBody,omitempty"`
	InjectHeaders []string              `json:"injectHeaders,omitempty"`
}

// RequestStat is one named request's aggregated outcome across every
// iteration it ran in — the same per-request view `reqx run`'s summary
// prints, not a flat per-iteration list (which would be unreadable at
// load-test scale — hundreds of iterations of the same request).
type RequestStat struct {
	Name         string `json:"name"`
	TotalRuns    int    `json:"totalRuns"`
	Successes    int    `json:"successes"`
	Failures     int    `json:"failures"`
	AvgLatencyMs int64  `json:"avgLatencyMs"`
	P95LatencyMs int64  `json:"p95LatencyMs"`
	TopError     string `json:"topError,omitempty"`
}

// RunCollectionSummary is the aggregate view — HDR-histogram-derived
// percentiles via internal/metrics.AnalyzeSharded, the same engine
// `reqx run` uses.
type RunCollectionSummary struct {
	TotalRequests   int     `json:"totalRequests"`
	TotalSuccess    int     `json:"totalSuccess"`
	TotalFailures   int     `json:"totalFailures"`
	SuccessRate     float64 `json:"successRate"`
	AvgLatencyMs    int64   `json:"avgLatencyMs"`
	P95LatencyMs    int64   `json:"p95LatencyMs"`
	TotalDurationMs int64   `json:"totalDurationMs"`
}

// RunCollectionOutput.DagNodes is only populated when the collection has
// depends_on edges (BuildExecutionPlan detected a DAG) — nil/omitted for a
// plain linear collection. Reuses internal/history.DagNodeRow — the exact
// same shape HistoryService.GetDAGNodes returns for a past run.
type RunCollectionOutput struct {
	Stats    []RequestStat        `json:"stats"`
	Summary  RunCollectionSummary `json:"summary"`
	DagNodes []history.DagNodeRow `json:"dagNodes,omitempty"`
}

// PickFile opens the native OS file picker filtered to *.<extension>, and
// returns the chosen path — empty string (no error) if the user cancels.
// One generic picker for every file input (collection/environment .json,
// personas .csv) instead of one method per file type.
func (s *CollectionService) PickFile(title string, extension string) (string, error) {
	if appCtx == nil {
		return "", errs.Internal("app is still starting up — try again in a moment")
	}
	pattern := "*." + extension
	path, err := wailsruntime.OpenFileDialog(appCtx, wailsruntime.OpenDialogOptions{
		Title: title,
		Filters: []wailsruntime.FileFilter{
			{DisplayName: strings.ToUpper(extension) + " Files (" + pattern + ")", Pattern: pattern},
		},
	})
	if err != nil {
		return "", errs.Wrap(err, errs.KindInternal, "could not open file picker")
	}
	return path, nil
}

// PickSaveFile opens the native OS save dialog defaulted to defaultFilename,
// and returns the chosen path — empty string (no error) if the user
// cancels. Used for Run's ExportPath (--export).
func (s *CollectionService) PickSaveFile(title string, defaultFilename string) (string, error) {
	if appCtx == nil {
		return "", errs.Internal("app is still starting up — try again in a moment")
	}
	path, err := wailsruntime.SaveFileDialog(appCtx, wailsruntime.SaveDialogOptions{
		Title:           title,
		DefaultFilename: defaultFilename,
	})
	if err != nil {
		return "", errs.Wrap(err, errs.KindInternal, "could not open save dialog")
	}
	return path, nil
}

// Open reads and parses a collection file — the app-side equivalent of the
// CLI reading args[0] in cmd/run_cmd_ctor.go. Returns the full Collection so
// the frontend can list requests and round-trip it back into Run.
func (s *CollectionService) Open(path string) (collection.Collection, error) {
	data, err := storage.ReadJSONFile(path)
	if err != nil {
		return collection.Collection{}, errs.Wrap(err, errs.KindInvalidInput, "could not read collection file")
	}
	coll, err := storage.ParseCollection(data)
	if err != nil {
		return collection.Collection{}, errs.Wrap(err, errs.KindInvalidInput, "could not parse collection JSON")
	}
	return *coll, nil
}

// OpenEnvironment reads and parses an environment file — the app-side
// equivalent of the CLI's `-e` flag (see cmd/run_cmd_ctor.go). Returns the
// full Environment so the frontend can show what's loaded and round-trip
// its Variables back into Run's EnvVariables.
func (s *CollectionService) OpenEnvironment(path string) (environment.Environment, error) {
	data, err := storage.ReadJSONFile(path)
	if err != nil {
		return environment.Environment{}, errs.Wrap(err, errs.KindInvalidInput, "could not read environment file")
	}
	env, err := storage.ParseEnvironment(data)
	if err != nil {
		return environment.Environment{}, errs.Wrap(err, errs.KindInvalidInput, "could not parse environment JSON")
	}
	return *env, nil
}

// OpenPersonas reads and parses a personas CSV file — the app-side
// equivalent of the CLI's `--personas` flag (see cmd/run_cmd_ctor.go).
// Every column becomes {{persona.<col>}} once passed back into Run.
func (s *CollectionService) OpenPersonas(path string) ([]personas.Persona, error) {
	return personas.LoadCSV(path)
}

// sequentialRunner picks which CollectionRunner the sequential path uses.
// The fast path (s.runner) is shared across every call for the app's
// lifetime specifically so connections and cookies persist BETWEEN
// separate Run clicks — but NoCookies/ClearCookies/GraphQL are per-call
// settings mutated on the runner/executor itself (SetClearCookiesPerRequest,
// SetGraphQLErrorCheck, DisableCookies), so applying them to the shared
// runner would leak into every future call, including ones that never
// asked for them. When any is requested, build a one-off runner instead —
// same pattern the CLI's sequential loop already uses (a fresh
// http_executor.NewDefaultExecutor() every iteration, cmd/run_cmd_ctor.go).
func (s *CollectionService) sequentialRunner(noCookies, clearCookies, graphqlCheck bool) *runner.CollectionRunner {
	if !noCookies && !clearCookies && !graphqlCheck {
		return s.runner
	}
	exec := http_executor.NewDefaultExecutor()
	if noCookies {
		exec.DisableCookies()
	}
	cr := runner.NewCollectionRunner(exec, nil, nil, nil)
	cr.SetVerbosity(runner.VerbosityQuiet)
	if clearCookies {
		cr.SetClearCookiesPerRequest(true)
	}
	if graphqlCheck {
		cr.SetGraphQLErrorCheck(true)
	}
	return cr
}

// buildBaseEnv turns a flat var map into an Environment, or nil if there's
// nothing to set — the load-testing paths below treat "no BaseEnv" as
// "start from a blank environment per worker", same as the CLI.
func buildBaseEnv(vars map[string]string) *environment.Environment {
	if len(vars) == 0 {
		return nil
	}
	env := environment.NewEnvironment("app")
	for k, v := range vars {
		env.Set(k, v)
	}
	return env
}

// Run executes the collection through whichever of the CLI's three engines
// input's load-testing fields select — same decision the CLI's `reqx run`
// makes between its Phase 3 Scheduler, WorkerPool, and plain sequential
// paths (cmd/run_cmd_ctor.go):
//   - Duration, RPS, or Stages set → Scheduler (duration-/rate-/ramp-limited
//     load test).
//   - Workers > 1 → WorkerPool (fixed concurrency, fixed iteration count).
//   - Otherwise → sequential iterations through the shared CollectionRunner
//     (Iterations defaults to 1 — the common single-run case).
func (s *CollectionService) Run(input RunCollectionInput) (RunCollectionOutput, error) {
	plan, err := planner.BuildExecutionPlan(&input.Collection, planner.PlanConfig{
		InjIndex:   input.InjectIndex,
		InjName:    input.InjectName,
		InjMethod:  input.InjectMethod,
		InjURL:     input.InjectURL,
		InjBody:    input.InjectBody,
		InjHeaders: input.InjectHeaders,
	})
	if err != nil {
		return RunCollectionOutput{}, err
	}

	iterations := input.Iterations
	if iterations < 1 {
		iterations = 1
	}
	workers := input.Workers
	if workers < 1 {
		workers = 1
	}
	duration := time.Duration(input.DurationMs) * time.Millisecond
	baseEnv := buildBaseEnv(input.EnvVariables)

	var stages []runner.Stage
	if input.Stages != "" {
		stages, err = runner.ParseStages(input.Stages)
		if err != nil {
			return RunCollectionOutput{}, err
		}
	}

	// Scheduler's Duration-mode conductor only stops on ctx cancellation —
	// with Duration==0 (and no Stages, which terminate on their own once
	// every stage's timer fires) that context has no deadline, so
	// RPS-without-Duration-or-Stages would block this call forever (same
	// gap exists in the CLI; there's at least a Ctrl+C there, here there'd
	// be no way to recover the UI).
	if input.RPS > 0 && duration <= 0 && len(stages) == 0 {
		return RunCollectionOutput{}, errs.InvalidInput("rps requires a duration or stages")
	}

	var allMetrics [][]runner.RequestMetric
	start := time.Now()

	switch {
	case duration > 0 || input.RPS > 0 || len(stages) > 0:
		cfg := runner.SchedulerConfig{
			Plan:         plan,
			BaseEnv:      baseEnv,
			Verbosity:    runner.VerbosityQuiet,
			Duration:     duration,
			MaxWorkers:   workers,
			RPS:          input.RPS,
			Stages:       stages,
			Personas:     input.Personas,
			NoCookies:    input.NoCookies,
			ClearCookies: input.ClearCookies,
			GraphQL:      input.GraphQL,
		}
		for _, r := range runner.NewScheduler(cfg).Run() {
			if r.Metrics != nil {
				allMetrics = append(allMetrics, r.Metrics)
			}
		}

	case workers > 1:
		cfg := runner.WorkerConfig{
			Plan:         plan,
			BaseEnv:      baseEnv,
			Verbosity:    runner.VerbosityQuiet,
			Personas:     input.Personas,
			NoCookies:    input.NoCookies,
			ClearCookies: input.ClearCookies,
			GraphQL:      input.GraphQL,
		}
		for _, r := range runner.NewWorkerPool(workers).Run(cfg, iterations) {
			allMetrics = append(allMetrics, r.Metrics)
		}

	default:
		activeRunner := s.sequentialRunner(input.NoCookies, input.ClearCookies, input.GraphQL)
		for i := 0; i < iterations; i++ {
			ctx := runner.NewRuntimeContext()
			if baseEnv != nil {
				ctx.SetEnvironment(baseEnv.Clone())
			}
			if len(input.Personas) > 0 {
				runner.ApplyPersona(ctx, input.Personas[0])
			}
			m, runErr := activeRunner.Run(plan, ctx)
			if runErr != nil {
				return RunCollectionOutput{}, errs.Wrap(runErr, errs.KindExternal, "collection run failed")
			}
			allMetrics = append(allMetrics, m)
		}
	}

	elapsed := time.Since(start)
	report := metrics.AnalyzeSharded(allMetrics, elapsed, 0)

	if input.ExportPath != "" {
		if exportErr := metrics.ExportJSON(allMetrics, input.ExportPath); exportErr != nil {
			log.Printf("[WARN] export failed: %v\n", exportErr)
		}
	}

	if s.history != nil {
		collectionName := input.Collection.Name
		if collectionName == "" {
			collectionName = "Untitled collection"
		}
		if saveErr := s.history.SaveRunWithDAG(collectionName, report, plan, allMetrics); saveErr != nil {
			log.Printf("[WARN] history save failed: %v\n", saveErr)
		}
	}

	var dagNodes []history.DagNodeRow
	if plan.DAG != nil && len(allMetrics) > 0 {
		if nodes, dagErr := history.ComputeDAGNodes(plan, allMetrics[0]); dagErr == nil {
			dagNodes = nodes
		} else {
			log.Printf("[WARN] dag node computation failed: %v\n", dagErr)
		}
	}

	stats := make([]RequestStat, len(report.PerRequest))
	for i, stat := range report.PerRequest {
		var topError string
		if len(stat.TopErrors) > 0 {
			topError = stat.TopErrors[0].Message
		}
		stats[i] = RequestStat{
			Name:         stat.Name,
			TotalRuns:    stat.TotalRuns,
			Successes:    stat.Successes,
			Failures:     stat.Failures,
			AvgLatencyMs: stat.AvgDuration.Milliseconds(),
			P95LatencyMs: stat.P95.Milliseconds(),
			TopError:     topError,
		}
	}

	return RunCollectionOutput{
		Stats:    stats,
		DagNodes: dagNodes,
		Summary: RunCollectionSummary{
			TotalRequests:   report.TotalRequests,
			TotalSuccess:    report.TotalSuccess,
			TotalFailures:   report.TotalFailures,
			SuccessRate:     report.SuccessRate,
			AvgLatencyMs:    report.AvgLatency.Milliseconds(),
			P95LatencyMs:    report.P95.Milliseconds(),
			TotalDurationMs: report.TotalDuration.Milliseconds(),
		},
	}, nil
}
