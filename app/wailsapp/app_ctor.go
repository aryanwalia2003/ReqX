package wailsapp

import "reqx/app/services"

// NewApp constructs every service and wires them into the App struct that
// gets bound to the frontend.
func NewApp() *App {
	return &App{
		Example: services.NewExampleService(),
	}
}
