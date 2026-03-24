package websocket_executor

import "time"

type defaultWebSocketExecutor struct {
	timeout time.Duration
	quiet   bool // when true, suppresses all per-message console output
}

//defaultWebSocketExecutor implements the WebSocketExecutor interface
