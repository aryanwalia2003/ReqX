package services

import (
	"time"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"reqx/internal/collection"
	"reqx/internal/environment"
	"reqx/internal/errs"
	"reqx/internal/metrics"
	"reqx/internal/planner"
	"reqx/internal/runner"
	"reqx/internal/storage"
)

// RunCollectionInput carries the collection to run (as returned by Open,
// round-tripped unedited for now) plus env vars for {{var}} substitution.
type RunCollectionInput struct {
	Collection   collection.Collection `json:"collection"`
	EnvVariables map[string]string     `json:"envVariables,omitempty"`
}

// RequestResult is one request's outcome — status/timing/error, no body
// (matches runner.RequestMetric, which drops the body to stay cheap across
// a whole collection; use RequestService.Send to inspect one response body).
type RequestResult struct {
	Name          string `json:"name"`
	Protocol      string `json:"protocol"`
	StatusCode    int    `json:"statusCode"`
	StatusString  string `json:"statusString"`
	DurationMs    int64  `json:"durationMs"`
	BytesReceived int64  `json:"bytesReceived"`
	ErrorMessage  string `json:"errorMessage,omitempty"`
}

// RunCollectionSummary is the aggregate view — HDR-histogram-derived
// percentiles via internal/metrics.Analyze, the same engine `reqx run` uses.
type RunCollectionSummary struct {
	TotalRequests   int     `json:"totalRequests"`
	TotalSuccess    int     `json:"totalSuccess"`
	TotalFailures   int     `json:"totalFailures"`
	SuccessRate     float64 `json:"successRate"`
	AvgLatencyMs    int64   `json:"avgLatencyMs"`
	P95LatencyMs    int64   `json:"p95LatencyMs"`
	TotalDurationMs int64   `json:"totalDurationMs"`
}

type RunCollectionOutput struct {
	Results []RequestResult      `json:"results"`
	Summary RunCollectionSummary `json:"summary"`
}

// PickFile opens the native OS file picker, filtered to .json files, and
// returns the chosen path — empty string (no error) if the user cancels.
// Used for both the collection and environment file inputs; title is the
// only thing that differs between them.
func (s *CollectionService) PickFile(title string) (string, error) {
	if appCtx == nil {
		return "", errs.Internal("app is still starting up — try again in a moment")
	}
	path, err := wailsruntime.OpenFileDialog(appCtx, wailsruntime.OpenDialogOptions{
		Title: title,
		Filters: []wailsruntime.FileFilter{
			{DisplayName: "JSON Files (*.json)", Pattern: "*.json"},
		},
	})
	if err != nil {
		return "", errs.Wrap(err, errs.KindInternal, "could not open file picker")
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

// Run executes every request in the collection, linearly or as a DAG
// (whichever BuildExecutionPlan detects from depends_on), through the same
// CollectionRunner/metrics.Analyze pipeline the CLI's `run` command uses.
func (s *CollectionService) Run(input RunCollectionInput) (RunCollectionOutput, error) {
	plan, err := planner.BuildExecutionPlan(&input.Collection, planner.PlanConfig{})
	if err != nil {
		return RunCollectionOutput{}, err
	}

	ctx := runner.NewRuntimeContext()
	for k, v := range input.EnvVariables {
		ctx.Environment.Set(k, v)
	}

	start := time.Now()
	metricsOut, err := s.runner.Run(plan, ctx)
	duration := time.Since(start)
	if err != nil {
		return RunCollectionOutput{}, errs.Wrap(err, errs.KindExternal, "collection run failed")
	}

	results := make([]RequestResult, len(metricsOut))
	for i, m := range metricsOut {
		results[i] = RequestResult{
			Name:          m.Name,
			Protocol:      m.Protocol,
			StatusCode:    m.StatusCode,
			StatusString:  m.StatusString,
			DurationMs:    m.Duration.Milliseconds(),
			BytesReceived: m.BytesReceived,
			ErrorMessage:  m.ErrorMsg,
		}
	}

	report := metrics.Analyze([][]runner.RequestMetric{metricsOut}, duration)

	return RunCollectionOutput{
		Results: results,
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
