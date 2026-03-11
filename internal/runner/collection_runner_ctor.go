package runner

import (
	"postman-cli/internal/http_executor"
	"postman-cli/internal/scripting"
)

// NewCollectionRunner constructs the orchestrator engine.
func NewCollectionRunner(exec http_executor.RequestExecutor, script scripting.ScriptRunner) *CollectionRunner {
	if exec == nil {
		exec = http_executor.NewDefaultExecutor()
	}
	if script == nil {
		script = scripting.NewGojaRunner()
	}
	
	return &CollectionRunner{
		executor:     exec,
		scriptRunner: script,
	}
}
