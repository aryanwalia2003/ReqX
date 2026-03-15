package runner

import (
	"reqx/internal/collection"
	"reqx/internal/environment"
	"reqx/internal/http_executor"
	"reqx/internal/scripting"
	"sync"
)

// WorkerConfig holds the shared, read-only inputs every worker needs.
type WorkerConfig struct {
	Coll        *collection.Collection
	BaseEnv     *environment.Environment // cloned per worker — never mutated
	NoCookies   bool
	ClearCookies bool
	Verbose     bool
}

// Run distributes totalIterations jobs across the pool and returns all results in order.
func (wp *WorkerPool) Run(cfg WorkerConfig, totalIterations int) []WorkerResult {
	jobs    := make(chan WorkerJob, totalIterations)
	results := make(chan WorkerResult, totalIterations)

	var wg sync.WaitGroup
	for w := 0; w < wp.numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runWorker(cfg, jobs, results)
		}()
	}

	for i := 1; i <= totalIterations; i++ {
		jobs <- WorkerJob{IterationIndex: i}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	all := make([]WorkerResult, 0, totalIterations)
	for r := range results {
		all = append(all, r)
	}
	return all
}

// runWorker consumes jobs and executes a full CollectionRunner per job.
func runWorker(cfg WorkerConfig, jobs <-chan WorkerJob, results chan<- WorkerResult) {
	for job := range jobs {
		metrics, err := executeIteration(cfg)
		results <- WorkerResult{
			IterationIndex: job.IterationIndex,
			Metrics:        metrics,
			Err:            err,
		}
	}
}

// executeIteration builds isolated state and runs one full collection pass.
func executeIteration(cfg WorkerConfig) ([]RequestMetric, error) {
	ctx := NewRuntimeContext()
	if cfg.BaseEnv != nil {
		ctx.SetEnvironment(cfg.BaseEnv.Clone())
	}

	exec := http_executor.NewDefaultExecutor()
	if cfg.NoCookies {
		exec.DisableCookies()
	}

	engine := NewCollectionRunner(exec, nil, nil, scripting.NewGojaRunner())
	if cfg.ClearCookies {
		engine.SetClearCookiesPerRequest(true)
	}
	if cfg.Verbose {
		engine.SetVerbose(true)
	}

	return engine.Run(cfg.Coll, ctx)
}