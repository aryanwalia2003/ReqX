package services

import "reqx/internal/runner"

// NewRequestService constructs the service with its own CollectionRunner
// (nil deps make it fall back to the default HTTP executor — see
// runner.NewCollectionRunner).
func NewRequestService() *RequestService {
	cr := runner.NewCollectionRunner(nil, nil, nil, nil)
	cr.SetVerbosity(runner.VerbosityQuiet) // no CLI console output from inside the desktop app
	return &RequestService{runner: cr}
}
