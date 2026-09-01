package wailsapp

import (
	"log"

	"reqx/app/services"
	"reqx/internal/history"
)

// NewApp constructs every service and wires them into the App struct that
// gets bound to the frontend. One history.DB is shared between
// CollectionService (writes a row per run) and HistoryService (reads them
// for the dashboard) — same connection, so a run is visible immediately.
//
// A history.Open failure degrades gracefully (nil db, logged) rather than
// aborting startup — same as the CLI, which silently skips history
// entirely when it can't open the file (cmd/run_cmd_ctor.go).
func NewApp() *App {
	db, err := history.Open() // nil on error — Open never returns a partial DB
	if err != nil {
		log.Printf("[WARN] history.Open failed, run history will not be saved: %v\n", err)
	}

	return &App{
		Example:    services.NewExampleService(),
		Request:    services.NewRequestService(),
		Collection: services.NewCollectionService(db),
		History:    services.NewHistoryService(db),
		Socket:     services.NewSocketService(),
	}
}
