package socketio_executor

import "postman-cli/internal/collection"

// SocketIOExecutor defines the interface for executing Socket.IO request flows.
type SocketIOExecutor interface {
	Execute(url string, headers map[string]string, events []collection.SocketIOEvent) error
}
