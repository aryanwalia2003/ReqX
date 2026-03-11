package runner

import "postman-cli/internal/environment"

// NewRuntimeContext constructs a new RuntimeContext.
func NewRuntimeContext() *RuntimeContext {
	return &RuntimeContext{
		GlobalVariables: make(map[string]interface{}),
		Environment:     environment.NewEnvironment("default"),
	}
}
