package runner

import (
	"reqx/internal/collection"
	"reqx/internal/environment"
	"reqx/internal/http_executor"
	"reqx/internal/progress"
	"reqx/internal/scripting"
	"sync"
)

// WorkerConfig holds the shared, read-only inputs every worker needs.
type WorkerConfig struct {
	Coll         *collection.Collection
	BaseEnv      *environment.Environment
	NoCookies    bool
	ClearCookies bool
	Verbosity    int // VerbosityQuiet suppresses per-request logs and shows progress bar
}

// Run distributes totalIterations jobs across the pool and returns all results.
// When Verbosity is VerbosityQuiet, a real-time progress bar is shown instead of logs.
func (wp *WorkerPool) Run(cfg WorkerConfig, totalIterations int) []WorkerResult {
	jobs    := make(chan WorkerJob, totalIterations)
	results := make(chan WorkerResult, totalIterations)

	var bar *progress.Bar
	if cfg.Verbosity == VerbosityQuiet {
		bar = progress.NewBar(totalIterations, wp.numWorkers)
		bar.Start()
	}

	var wg sync.WaitGroup
	for w := 0; w < wp.numWorkers; w++ {
		workerID := w + 1
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			runWorker(cfg, id, jobs, results, bar)
		}(workerID)
	}

	for i := 1; i <= totalIterations; i++ {
		jobs <- WorkerJob{IterationIndex: i}
	}
	close(jobs)

	go func() {
		wg.Wait()
		if bar != nil {
			bar.Stop()
		}
		close(results)
	}()

	all := make([]WorkerResult, 0, totalIterations)
	for r := range results {
		all = append(all, r)
	}
	return all
}

// runWorker consumes jobs and executes one CollectionRunner per job.
func runWorker(cfg WorkerConfig, workerID int, jobs <-chan WorkerJob, results chan<- WorkerResult, bar *progress.Bar) {
	for job := range jobs {
		metrics, err := executeIteration(cfg, workerID)
		if bar != nil {
			bar.IncrementDone()
			if err != nil {
				bar.IncrementErrors()
			}
		}
		results <- WorkerResult{IterationIndex: job.IterationIndex, Metrics: metrics, Err: err}
	}
}

// executeIteration builds isolated state and runs one full collection pass.
func executeIteration(cfg WorkerConfig, workerID int) ([]RequestMetric, error) {
	ctx := NewRuntimeContext()
	if cfg.BaseEnv != nil {
		ctx.SetEnvironment(cfg.BaseEnv.Clone())
	}

	exec := http_executor.NewDefaultExecutor()
	if cfg.NoCookies {
		exec.DisableCookies()
	}

	engine := NewCollectionRunner(exec, nil, nil, scripting.NewGojaRunner())
	engine.SetVerbosity(cfg.Verbosity)
	if cfg.ClearCookies {
		engine.SetClearCookiesPerRequest(true)
	}

	metrics, err := engine.Run(cfg.Coll, ctx)

	// Tag every metric with this worker's ID for export traceability
	for i := range metrics {
		metrics[i].WorkerID = workerID
	}
	return metrics, err
}