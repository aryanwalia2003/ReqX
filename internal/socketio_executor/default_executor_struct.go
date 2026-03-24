package socketio_executor

import (
	"time"
)

// DefaultSocketIOExecutor implements the Socket.IO flow using njones/socketio.
type DefaultSocketIOExecutor struct {
	timeout time.Duration
	quiet   bool // when true, suppresses all per-event console output
}
