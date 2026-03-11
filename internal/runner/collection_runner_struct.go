package runner

import (
	"postman-cli/internal/http_executor"
	"postman-cli/internal/scripting"
	"postman-cli/internal/socketio_executor"
)

// CollectionRunner handles executing a full collection of requests.
type CollectionRunner struct {
	executor     http_executor.RequestExecutor
	sioExecutor  socketio_executor.SocketIOExecutor
	scriptRunner scripting.ScriptRunner
}

