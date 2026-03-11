package scripting

// NewGojaRunner constructs a ScriptRunner that uses the dop251/goja VM.
func NewGojaRunner() ScriptRunner {
	return &GojaRunner{}
}
