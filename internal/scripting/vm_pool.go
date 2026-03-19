package scripting

import (
	"sync"
	"sync/atomic"

	"github.com/dop251/goja"
	"github.com/fatih/color"
)

var vmAllocCount int32

// vmPool recycles *goja.Runtime instances so that the expensive
// goja.New() allocation only happens once per concurrent goroutine
// rather than once per script execution.
var vmPool = sync.Pool{
	New: func() any {
		allocs := atomic.AddInt32(&vmAllocCount, 1)
		color.Cyan("[POOL] Allocating fresh JavaScript VM... (Total VMs: %d)\n", allocs)
		
		vm := goja.New()
		vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
		return vm
	},
}

// acquireVM pulls a ready-to-use VM from the pool.
func acquireVM() *goja.Runtime {
	return vmPool.Get().(*goja.Runtime)
}

// releaseVM clears injected globals and returns the VM to the pool.
func releaseVM(vm *goja.Runtime) {
	// Wipe script-injected bindings to prevent state bleed.
	vm.Set("pm", nil)
	vm.Set("console", nil)
	vmPool.Put(vm)
}
