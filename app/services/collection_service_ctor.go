package services

import (
	"reqx/internal/history"
	"reqx/internal/runner"
)

// NewCollectionService constructs the service with its own CollectionRunner
// (nil deps fall back to the default HTTP executor). db may be nil — see
// the history field's comment.
func NewCollectionService(db *history.DB) *CollectionService {
	cr := runner.NewCollectionRunner(nil, nil, nil, nil)
	cr.SetVerbosity(runner.VerbosityQuiet)
	return &CollectionService{runner: cr, history: db}
}
