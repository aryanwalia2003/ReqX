package runner

// WorkerJob represents a single iteration assigned to a worker.
type WorkerJob struct {
	IterationIndex int
}

// WorkerResult carries the metrics produced by one completed iteration.
type WorkerResult struct {
	IterationIndex int
	Metrics        []RequestMetric
	Err            error
}

// WorkerPool orchestrates parallel collection runs across N workers.
type WorkerPool struct {
	numWorkers int
}