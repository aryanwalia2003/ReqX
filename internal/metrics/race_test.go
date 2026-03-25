package metrics

import (
	"errors"
	"sync"
	"testing"
	"time"

	"reqx/internal/runner"
)

// TestAnalyzeSharded_ConcurrentProducers fires N goroutines that all produce
// metrics simultaneously and feeds them to AnalyzeSharded.
// Run with: go test ./internal/metrics/... -race
//
// This validates that the channel-based dispatcher and shard goroutines
// are free of data races when many goroutines produce results concurrently.
func TestAnalyzeSharded_ConcurrentProducers(t *testing.T) {
	const (
		workers   = 50
		requests  = 20
		shardCnt  = 4
	)

	// Build allMetrics: one slice per worker, all written concurrently.
	allMetrics := make([][]runner.RequestMetric, workers)
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(idx int) {
			defer wg.Done()
			slice := make([]runner.RequestMetric, requests)
			for j := 0; j < requests; j++ {
				slice[j] = runner.RequestMetric{
					Name:       "request/endpoint",
					Protocol:   "HTTP",
					StatusCode: 200,
					Duration:   time.Duration(j+1) * time.Millisecond,
					BytesSent:  100,
					BytesReceived: 200,
				}
			}
			allMetrics[idx] = slice
		}(i)
	}
	wg.Wait()

	// AnalyzeSharded must not race internally despite receiving a large dataset.
	report := AnalyzeSharded(allMetrics, 5*time.Second, shardCnt)

	total := workers * requests
	if report.TotalRequests != total {
		t.Errorf("expected %d total requests, got %d", total, report.TotalRequests)
	}
	if report.TotalSuccess != total {
		t.Errorf("expected %d successes, got %d", total, report.TotalSuccess)
	}
}

// TestAnalyzeSharded_MixedErrorsAndSuccess verifies that failure accounting
// is still race-free when errors and successes are interleaved across workers.
func TestAnalyzeSharded_MixedErrorsAndSuccess(t *testing.T) {
	const workers = 20
	allMetrics := make([][]runner.RequestMetric, workers)
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(idx int) {
			defer wg.Done()
			allMetrics[idx] = []runner.RequestMetric{
				{Name: "A", Protocol: "HTTP", StatusCode: 200, Duration: 50 * time.Millisecond},
				{Name: "A", Protocol: "HTTP", StatusCode: 500, Duration: 10 * time.Millisecond, Error: errors.New("internal server error")},
				{Name: "B", Protocol: "HTTP", StatusCode: 201, Duration: 30 * time.Millisecond},
			}
		}(i)
	}
	wg.Wait()

	report := AnalyzeSharded(allMetrics, 10*time.Second, 2)

	if report.TotalRequests != workers*3 {
		t.Errorf("expected %d total, got %d", workers*3, report.TotalRequests)
	}
	if report.TotalFailures != workers {
		t.Errorf("expected %d failures, got %d", workers, report.TotalFailures)
	}
	if report.TotalSuccess != workers*2 {
		t.Errorf("expected %d successes, got %d", workers*2, report.TotalSuccess)
	}
}

// TestAnalyzeSharded_MultipleRequestNames verifies that per-name sharding
// correctly partitions and then merges disjoint request names.
func TestAnalyzeSharded_MultipleRequestNames(t *testing.T) {
	names := []string{"login", "dashboard", "checkout", "logout"}
	allMetrics := make([][]runner.RequestMetric, 30)
	var wg sync.WaitGroup
	wg.Add(30)
	for i := 0; i < 30; i++ {
		go func(idx int) {
			defer wg.Done()
			slice := make([]runner.RequestMetric, len(names))
			for j, n := range names {
				slice[j] = runner.RequestMetric{
					Name: n, Protocol: "HTTP",
					StatusCode: 200, Duration: time.Duration(idx+1) * time.Millisecond,
				}
			}
			allMetrics[idx] = slice
		}(i)
	}
	wg.Wait()

	report := AnalyzeSharded(allMetrics, 5*time.Second, 8)

	if len(report.PerRequest) != len(names) {
		t.Errorf("expected %d distinct requests, got %d", len(names), len(report.PerRequest))
	}
	for _, stat := range report.PerRequest {
		if stat.TotalRuns != 30 {
			t.Errorf("request %q: expected 30 runs, got %d", stat.Name, stat.TotalRuns)
		}
	}
}

// TestAnalyzeSharded_StatusCodes ensures status-code aggregation is correct
// under concurrent writes from multiple shards.
func TestAnalyzeSharded_StatusCodes(t *testing.T) {
	const workers = 10
	allMetrics := make([][]runner.RequestMetric, workers)
	for i := 0; i < workers; i++ {
		allMetrics[i] = []runner.RequestMetric{
			{Name: "req", Protocol: "HTTP", StatusCode: 200, Duration: 1 * time.Millisecond},
			{Name: "req", Protocol: "HTTP", StatusCode: 404, Duration: 2 * time.Millisecond, Error: errors.New("not found")},
		}
	}

	report := AnalyzeSharded(allMetrics, 5*time.Second, 4)

	if report.StatusCodes[200] != workers {
		t.Errorf("status 200: expected %d, got %d", workers, report.StatusCodes[200])
	}
	if report.StatusCodes[404] != workers {
		t.Errorf("status 404: expected %d, got %d", workers, report.StatusCodes[404])
	}
}
