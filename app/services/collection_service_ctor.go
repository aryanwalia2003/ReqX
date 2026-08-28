package services

import "reqx/internal/runner"

// NewCollectionService constructs the service with its own CollectionRunner
// (nil deps fall back to the default HTTP executor).
func NewCollectionService() *CollectionService {
	cr := runner.NewCollectionRunner(nil, nil, nil, nil)
	cr.SetVerbosity(runner.VerbosityQuiet)
	return &CollectionService{runner: cr}
}
