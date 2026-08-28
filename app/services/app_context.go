package services

import "context"

// appCtx is the Wails runtime context — needed for OS-level runtime calls
// like file dialogs. Package-level and set once via SetAppContext (from
// wailsapp's OnStartup, after wails.Run actually starts the event loop),
// not a field+setter method on any individual service: Wails' Bind only
// walks the exported methods of the structs handed to it, so a method here
// would risk becoming an accidental (and useless — the frontend can't
// supply a context.Context) frontend-callable binding. A free function
// never has that problem.
var appCtx context.Context

// SetAppContext wires the Wails runtime context for every service in this
// package that needs one.
func SetAppContext(ctx context.Context) {
	appCtx = ctx
}
