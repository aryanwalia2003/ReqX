package runner

import (
	"reqx/internal/environment"
	"sync"
)

// NewRuntimeContext constructs a new RuntimeContext.
func NewRuntimeContext() *RuntimeContext {
	return &RuntimeContext{
		GlobalVariables: make(map[string]interface{}),
		Environment:     environment.NewEnvironment("default"),
		AsyncWG:         new(sync.WaitGroup),
		AsyncStop:       make(chan struct{}),
		AsyncStopOnce:   new(sync.Once),
		ownsAsyncStop:   true,
	}
}

// CloneForNode creates a lightweight child context for a DAG node goroutine.
// The clone:
//   - gets its own Environment snapshot so concurrent pm.env.set is race-free
//   - shares AsyncWG / AsyncStop / AsyncStopOnce so background tasks are
//     tracked by the single top-level RunDAG defer
//   - has ownsAsyncStop = false so runLinear does NOT close the stop channel
func (rc *RuntimeContext) CloneForNode() *RuntimeContext {
	return &RuntimeContext{
		GlobalVariables: rc.GlobalVariables,
		Environment:     newEnvSnapshot(rc.Environment),
		AsyncWG:         rc.AsyncWG,
		AsyncStop:       rc.AsyncStop,
		AsyncStopOnce:   rc.AsyncStopOnce,
		ownsAsyncStop:   false,
	}
}
