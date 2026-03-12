package scripting

import (
	"fmt"
	"reflect"
) //reflect is used for deep comparison of values like objects and arrays for eg: pm.expect(pm.response.json()).toEql({key: "value"}) and pm.expect(pm.response.json()).toBe({key: "value"}) and pm.expect(pm.response.json()).toExist()

type defaultExpectBuilder struct {
	value       interface{}
	testResults *TestResults
}

// markLastTestFailed updates the most recently added test block to FALSE.
// Since Goja executes synchronously, `pm.expect` within `pm.test` can modify the last pushed state safely.
func (b *defaultExpectBuilder) markLastTestFailed(reason string) {
	if b.testResults == nil || len(*b.testResults) == 0 {
		return
	}
	idx := len(*b.testResults) - 1
	(*b.testResults)[idx].Passed = false
	(*b.testResults)[idx].Error = reason
}

func (b *defaultExpectBuilder) ToEql(expected interface{}) {
	if !reflect.DeepEqual(b.value, expected) {
		b.markLastTestFailed(fmt.Sprintf("expected %v to equal %v", b.value, expected))
		// Instead of throwing a JS exception that crashes the whole script,
		// we just fail the test block cleanly. 
	}
}

func (b *defaultExpectBuilder) ToBe(expected interface{}) {
	// For JS, simple equality often uses ToBe
	if b.value != expected {
		b.markLastTestFailed(fmt.Sprintf("expected %v to be %v", b.value, expected))
	}
}

func (b *defaultExpectBuilder) ToExist() {
	if b.value == nil {
		b.markLastTestFailed("expected value to exist, got null/undefined")
	}
}
