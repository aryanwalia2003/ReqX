package runner

import (
	"postman-cli/internal/http_executor"
	"postman-cli/internal/scripting"
)

// CollectionRunner handles executing a full collection of requests.
type CollectionRunner struct {
	executor     http_executor.RequestExecutor
	scriptRunner scripting.ScriptRunner
}
