package services

import "reqx/internal/runner"

// CollectionService opens and runs a whole collection from the desktop app —
// the app-side equivalent of the CLI's `reqx run` command. Like
// RequestService, it holds one CollectionRunner for the app's lifetime so
// runs share the same pooled transport and cookie jar.
type CollectionService struct {
	runner *runner.CollectionRunner
}
