package services

import (
	"sync"

	"github.com/gorilla/websocket"
)

// SocketService drives one interactive WebSocket/Socket.IO debugger
// connection at a time — the desktop-app equivalent of the CLI's `reqx ws`
// and `reqx sio` REPLs, except every frame is pushed to the frontend via a
// Wails event instead of stdout.
type SocketService struct {
	mu       sync.Mutex
	writeMu  sync.Mutex // guards conn.WriteMessage (gorilla forbids concurrent writers)
	conn     *websocket.Conn
	protocol string // "ws" | "sio"
	sioReady bool   // true once the Socket.IO (not just TCP) handshake completes
}
