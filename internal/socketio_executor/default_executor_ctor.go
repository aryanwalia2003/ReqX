package socketio_executor

import "time"

// NewDefaultSocketIOExecutor creates a standard SocketIOExecutor.
func NewDefaultSocketIOExecutor() SocketIOExecutor {
	return &DefaultSocketIOExecutor{
		timeout: 10 * time.Second,
	}
}
