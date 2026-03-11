package socketio_executor

import (
	"time"
)

// DefaultSocketIOExecutor implements the Socket.IO flow using njones/socketio.
type DefaultSocketIOExecutor struct {
	timeout time.Duration
}
