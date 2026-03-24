package scripting

import (
	"sync"

	"github.com/dop251/goja"
)

// NewGojaRunner constructs a ScriptRunner that uses the dop251/goja VM backed by a sync.Pool.
func NewGojaRunner() ScriptRunner {
	return &GojaRunner{
		pool: &sync.Pool{
			New: func() interface{} {
				vm := goja.New()
				vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
				return vm
			},
		},
	}
}
