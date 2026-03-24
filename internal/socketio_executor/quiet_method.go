package socketio_executor

// SetQuiet silences all per-event console output when q is true.
// Call this before Execute to suppress the [RECEIVED] / "Connected" spam
// during load tests so that stdout is not the CPU bottleneck.
func (e *DefaultSocketIOExecutor) SetQuiet(q bool) {
	e.quiet = q
}
