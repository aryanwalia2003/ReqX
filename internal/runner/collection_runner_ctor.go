package runner

import (
	"reqx/internal/http_executor"
	"reqx/internal/scripting"
	"reqx/internal/socketio_executor"
	"reqx/internal/websocket_executor"
	"sync"
)

// NewCollectionRunner constructs the orchestration engine.
// exec must be a *http_executor.DefaultExecutor so cookie control is available.
func NewCollectionRunner(exec *http_executor.DefaultExecutor, sio socketio_executor.SocketIOExecutor, we websocket_executor.WebSocketExecutor, script scripting.ScriptRunner) *CollectionRunner {
	if exec == nil {
		exec = http_executor.NewDefaultExecutor()
	}
	if sio == nil {
		sio = socketio_executor.NewDefaultSocketIOExecutor()
	}
	if we == nil {
		we = websocket_executor.NewDefaultWebSocketExecutor()
	}
	if script == nil {
		script = scripting.NewGojaRunner()
	}

	return &CollectionRunner{
		executor:     exec,
		sioExecutor:  sio,
		weExecutor:   we,
		scriptRunner: script,
		verbosity:    VerbosityNormal,
		wg:          &sync.WaitGroup{},
	}
}

// SetVerbosity controls how much per-request output the runner emits.
func (cr *CollectionRunner) SetVerbosity(v int) {
	cr.verbosity = v
}

// SetVerbose is kept for backward compatibility; it maps to VerbosityFull.
func (cr *CollectionRunner) SetVerbose(v bool) {
	if v {
		cr.verbosity = VerbosityFull
	} else if cr.verbosity < VerbosityNormal {
		cr.verbosity = VerbosityNormal
	}
}

