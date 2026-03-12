package runner

import (
	"postman-cli/internal/http_executor"
	"postman-cli/internal/scripting"
	"postman-cli/internal/socketio_executor"
)

// NewCollectionRunner constructs the orchestration engine.
// exec must be a *http_executor.DefaultExecutor so cookie control is available.
func NewCollectionRunner(exec *http_executor.DefaultExecutor, sio socketio_executor.SocketIOExecutor, script scripting.ScriptRunner) *CollectionRunner {
	if exec == nil {
		exec = http_executor.NewDefaultExecutor()
	}
	if sio == nil {
		sio = socketio_executor.NewDefaultSocketIOExecutor()
	}
	if script == nil {
		script = scripting.NewGojaRunner()
	}

	return &CollectionRunner{
		executor:     exec,
		sioExecutor:  sio,
		scriptRunner: script,
	}
}

