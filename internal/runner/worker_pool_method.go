package runner

import (
	"reqx/internal/environment"
	"reqx/internal/http_executor"
	"reqx/internal/personas"
	"reqx/internal/planner"
	"reqx/internal/progress"
	"reqx/internal/scripting"
	"sync"
)

// WorkerConfig holds the shared, read-only inputs every WorkerPool goroutine needs.
type WorkerConfig struct {
	Plan         *planner.ExecutionPlan   // immutable; replaces Coll
	BaseEnv      *environment.Environment // cloned per worker
	NoCookies    bool
	ClearCookies bool
	Verbosity    int
	Personas     []personas.Persona
}

// Run distributes totalIterations jobs across the pool and returns all results.
func (wp *WorkerPool) Run(cfg WorkerConfig, totalIterations int) []WorkerResult {
	jobs := make(chan WorkerJob, totalIterations)
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
			poolWorkerLoop(cfg, id, jobs, results, bar)
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

// poolWorkerLoop consumes jobs and executes one isolated plan pass per job.
func poolWorkerLoop(cfg WorkerConfig, workerID int, jobs <-chan WorkerJob, results chan<- WorkerResult, bar *progress.Bar) {
	for job := range jobs {
		metrics, err := poolExecuteOne(cfg, workerID)
		if bar != nil {
			bar.IncrementDone()
			if err != nil {
				bar.IncrementErrors()
			}
		}
		results <- WorkerResult{IterationIndex: job.IterationIndex, Metrics: metrics, Err: err}
	}
}

func poolExecuteOne(cfg WorkerConfig, workerID int) ([]RequestMetric, error) {
	rtCtx := NewRuntimeContext()
	if cfg.BaseEnv != nil {
		rtCtx.SetEnvironment(cfg.BaseEnv.Clone())
	}
	if len(cfg.Personas) > 0 {
		p := cfg.Personas[(workerID-1)%len(cfg.Personas)]
		applyPersona(rtCtx, p)
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

	metrics, err := engine.Run(cfg.Plan, rtCtx)
	for i := range metrics {
		metrics[i].WorkerID = workerID
	}
	return metrics, err
}
