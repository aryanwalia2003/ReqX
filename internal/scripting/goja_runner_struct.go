package scripting

import (
	"sync"
)

// GojaRunner is the concrete implementation of ScriptRunner using the goja JS engine.
// It uses a sync.Pool of *goja.Runtime to allow parallel execution (DAG) without
// concurrent map access panics or heavy allocation overhead.
type GojaRunner struct {
	pool *sync.Pool
}
