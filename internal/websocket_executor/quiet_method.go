package websocket_executor

// SetQuiet silences all per-message console output when q is true.
// Call this before Execute to suppress [WS_INCOMING] / "Connected" spam
// during load tests so that stdout is not the CPU bottleneck.
func (e *defaultWebSocketExecutor) SetQuiet(q bool) {
	e.quiet = q
}
