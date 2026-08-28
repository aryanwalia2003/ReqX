package wailsapp

import (
	"context"

	"reqx/app/services"
)

// BindTargets returns every service exposed to the frontend. Add a new
// service here as it's wired in — each becomes
// window.go.services.<ServiceName>.<methodName>(...) in the generated JS
// bindings once `wails generate module` (or `wails dev`) runs.
func (a *App) BindTargets() []interface{} {
	return []interface{}{
		a.Example,
		a.Request,
		a.Collection,
	}
}

// OnStartup runs once wails.Run actually starts the event loop — this is
// the first point a runtime context exists, needed for OS-level runtime
// calls like file dialogs (see services.PickFile). Wired as
// options.App.OnStartup in main.go.
func (a *App) OnStartup(ctx context.Context) {
	services.SetAppContext(ctx)
}
