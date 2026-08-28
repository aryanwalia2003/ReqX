package wailsapp

// BindTargets returns every service exposed to the frontend. Add a new
// service here as it's wired in — each becomes
// window.go.services.<ServiceName>.<methodName>(...) in the generated JS
// bindings once `wails generate module` (or `wails dev`) runs.
func (a *App) BindTargets() []interface{} {
	return []interface{}{
		a.Example,
		a.Request,
	}
}
