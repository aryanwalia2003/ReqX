package runner

// NewWorkerPool constructs a WorkerPool with the given concurrency level.
func NewWorkerPool(numWorkers int) *WorkerPool {
	if numWorkers < 1 {
		numWorkers = 1
	}
	return &WorkerPool{numWorkers: numWorkers}
}