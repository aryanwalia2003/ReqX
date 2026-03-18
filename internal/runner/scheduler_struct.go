package runner

import (
	"context"
	"sync/atomic"
	"sync"
)

// Scheduler orchestrates duration, RPS, and stage-based load tests.
// It is the Phase 3 replacement for WorkerPool when dynamic control is needed.
type Scheduler struct {
	cfg SchedulerConfig

	// Live state — read by the progress renderer.
	activeWorkers atomic.Int64
	completedIterations atomic.Int64
	failedIterations    atomic.Int64

	// Conductor-controlled desired concurrency.
	// Workers with id > desiredWorkers go idle (no spinning).
	desiredWorkers atomic.Int64

	results chan WorkerResult
	wg      sync.WaitGroup
	cancel  context.CancelFunc

	// Wakes idled workers when desiredWorkers increases.
	// This avoids a central per-iteration job queue while still allowing stage ramps.
	wakeMu sync.Mutex
	wakeCh chan struct{}
}