package metrics

import "time"

// RequestStat holds aggregated performance data for a single named request.
type RequestStat struct {
	Name        string
	TotalRuns   int
	Successes   int
	Failures    int
	Durations   []time.Duration // all HTTP durations, sorted after Analyze
	P50         time.Duration
	P90         time.Duration
	P95         time.Duration
	P99         time.Duration
	AvgDuration time.Duration
	TopErrors   []ErrorGroup
}

// ErrorGroup tracks how many times a specific error message occurred.
type ErrorGroup struct {
	Message string
	Count   int
}

// Report is the complete output of Analyze — everything needed to render any summary.
type Report struct {
	TotalRequests int
	TotalSuccess  int
	TotalFailures int
	SuccessRate   float64
	AvgLatency    time.Duration
	P50           time.Duration
	P90           time.Duration
	P95           time.Duration
	P99           time.Duration
	RPS           float64
	TotalDuration time.Duration
	PerRequest    []RequestStat
}