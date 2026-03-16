package progress

import (
	"fmt"
	"sync/atomic"
	"time"
)

// Bar renders a real-time progress bar to stdout using carriage-return overwriting.
// It is safe to call Increment from multiple goroutines.
type Bar struct {
	total     int64
	done      atomic.Int64
	errors    atomic.Int64
	workers   int
	startTime time.Time
	stopCh    chan struct{}
}

// NewBar constructs a Bar for a run with totalJobs iterations and numWorkers goroutines.
func NewBar(totalJobs, numWorkers int) *Bar {
	return &Bar{
		total:     int64(totalJobs),
		workers:   numWorkers,
		startTime: time.Now(),
		stopCh:    make(chan struct{}),
	}
}

// Start begins the background render loop.
func (b *Bar) Start() {
	go func() {
		ticker := time.NewTicker(250 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				b.render()
			case <-b.stopCh:
				b.render() // final render
				fmt.Println() // newline after bar
				return
			}
		}
	}()
}

// Stop halts the render loop. Call after all jobs are done.
func (b *Bar) Stop() { close(b.stopCh) }

// IncrementDone records one completed iteration (pass or fail).
func (b *Bar) IncrementDone() { b.done.Add(1) }

// IncrementErrors records one failed iteration.
func (b *Bar) IncrementErrors() { b.errors.Add(1) }