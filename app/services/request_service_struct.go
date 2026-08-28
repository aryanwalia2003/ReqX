package services

import "reqx/internal/runner"

// RequestService sends a single ad-hoc HTTP request from the desktop app —
// the app-side equivalent of the CLI's `reqx req` command. It holds one
// CollectionRunner (and therefore one pooled http.Transport + cookie jar)
// for the lifetime of the app, so repeated sends reuse connections and
// accumulate cookies exactly like the CLI's engine does.
type RequestService struct {
	runner *runner.CollectionRunner
}
