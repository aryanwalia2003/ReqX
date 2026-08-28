package services

import (
	"reqx/internal/history"
	"reqx/internal/runner"
)

// CollectionService opens and runs a whole collection from the desktop app —
// the app-side equivalent of the CLI's `reqx run` command. Like
// RequestService, it holds one CollectionRunner for the app's lifetime so
// runs share the same pooled transport and cookie jar.
type CollectionService struct {
	runner *runner.CollectionRunner

	// history is nil when history.Open() failed at startup (see
	// wailsapp.NewApp) — Run() then just skips saving, same graceful
	// degradation as the CLI (cmd/run_cmd_ctor.go's printAndExport).
	history *history.DB
}
