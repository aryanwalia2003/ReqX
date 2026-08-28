// Package wailsapp wires the Wails desktop shell: window options, the
// embedded frontend build, and the service structs bound to the frontend.
// It never contains business logic itself — that lives in internal/* and
// gets wrapped by app/services.
package wailsapp

import "reqx/app/services"

// App holds every service struct bound to the frontend. Add a field here
// (and wire it in NewApp) for each new service — Wails binds every exported
// field's exported methods automatically.
type App struct {
	Example    *services.ExampleService
	Request    *services.RequestService
	Collection *services.CollectionService
}
