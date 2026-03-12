package scripting

import (
	"postman-cli/internal/collection"
	"postman-cli/internal/environment"
)

// ScriptRunner defines the interface for executing pre-request and test scripts.
type ScriptRunner interface {
	Execute(script *collection.Script, env *environment.Environment, resp *ResponseAPI) error //takes a pointer to the script struct, a pointer to the environment struct, and a response pointer
}

