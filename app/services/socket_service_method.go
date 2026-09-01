package services

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"reqx/internal/errs"
)

// SocketMessage is one line of the debugger's live transcript, pushed to the
// frontend as "socket:message" events as frames arrive.
type SocketMessage struct {
	Direction string `json:"direction"` // "in" | "out" | "system"
	Data      string `json:"data"`
	EventName string `json:"eventName,omitempty"` // decoded Socket.IO event name, if any
	Timestamp int64  `json:"timestamp"`
}

func (s *SocketService) emit(direction, data, eventName string) {
	// appCtx is unset in tests and before Wails' OnStartup fires —
	// EventsEmit would log.Fatal on a nil ctx, so just drop the event.
	if appCtx == nil {
		return
	}
	wailsruntime.EventsEmit(appCtx, "socket:message", SocketMessage{
		Direction: direction, Data: data, EventName: eventName,
		Timestamp: time.Now().UnixMilli(),
	})
}

// ConnectSocketInput mirrors the CLI's `reqx ws`/`reqx sio` connect args.
type ConnectSocketInput struct {
	URL      string            `json:"url"`
	Protocol string            `json:"protocol"` // "ws" | "sio"
	Headers  map[string]string `json:"headers,omitempty"`
}

// Connect opens the debugger's one active connection — same raw-dial and
// Socket.IO v4 handshake URL rewrite as `reqx ws`/`reqx sio`
// (cmd/ws_cmd_ctor.go, cmd/sio_cmd_ctor.go), then hands off to a background
// read loop instead of blocking on stdin.
func (s *SocketService) Connect(input ConnectSocketInput) error {
	s.mu.Lock()
	if s.conn != nil {
		s.mu.Unlock()
		return errs.InvalidInput("already connected — disconnect first")
	}
	s.mu.Unlock()

	rawURL := strings.TrimSpace(input.URL)
	protocol := input.Protocol
	if protocol == "" {
		protocol = "ws"
	}

	dialURL, err := socketDialURL(rawURL, protocol)
	if err != nil {
		return err
	}

	reqHeaders := http.Header{}
	for k, v := range input.Headers {
		reqHeaders.Add(k, v)
	}

	conn, _, err := websocket.DefaultDialer.Dial(dialURL, reqHeaders)
	if err != nil {
		return errs.Wrap(err, errs.KindExternal, "failed to connect")
	}

	s.mu.Lock()
	s.conn = conn
	s.protocol = protocol
	s.sioReady = false
	s.mu.Unlock()

	s.emit("system", "Connected", "")
	go s.readLoop(conn, protocol)

	return nil
}

// socketDialURL turns a plain http(s)/ws(s) URL into the URL actually dialed.
// For "sio" it applies the same rewrite as the CLI: http→ws, https→wss,
// default path "/socket.io/", and the Engine.IO v4 websocket query params.
func socketDialURL(rawURL, protocol string) (string, error) {
	if protocol != "sio" {
		if !strings.HasPrefix(rawURL, "ws://") && !strings.HasPrefix(rawURL, "wss://") {
			return "", errs.InvalidInput("URL must start with ws:// or wss://")
		}
		return rawURL, nil
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", errs.InvalidInput("invalid URL: " + err.Error())
	}
	switch u.Scheme {
	case "http":
		u.Scheme = "ws"
	case "https":
		u.Scheme = "wss"
	}
	if u.Path == "" || u.Path == "/" {
		u.Path = "/socket.io/"
	}
	q := u.Query()
	q.Set("EIO", "4")
	q.Set("transport", "websocket")
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// readLoop streams every incoming frame to the frontend until the connection
// closes. For "sio" it also speaks the Engine.IO handshake (0→40, 2→3) so the
// server sees a well-behaved client, same as the CLI's REPL.
func (s *SocketService) readLoop(conn *websocket.Conn, protocol string) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			s.emit("system", "Disconnected: "+err.Error(), "")
			s.mu.Lock()
			if s.conn == conn {
				s.conn = nil
			}
			s.mu.Unlock()
			return
		}

		msgStr := string(message)
		if protocol != "sio" {
			s.emit("in", msgStr, "")
			continue
		}

		switch {
		case strings.HasPrefix(msgStr, "42"):
			var arr []interface{}
			if json.Unmarshal([]byte(msgStr[2:]), &arr) == nil && len(arr) > 0 {
				eventName, _ := arr[0].(string)
				payload := ""
				if len(arr) > 1 {
					if b, err := json.Marshal(arr[1]); err == nil {
						payload = string(b)
					}
				}
				s.emit("in", payload, eventName)
			}
		case strings.HasPrefix(msgStr, "40"):
			s.mu.Lock()
			s.sioReady = true
			s.mu.Unlock()
			s.emit("system", "Socket.IO connected", "")
		case strings.HasPrefix(msgStr, "0"):
			s.writeRaw("40")
		case strings.HasPrefix(msgStr, "2"):
			s.writeRaw("3")
		default:
			s.emit("system", msgStr, "")
		}
	}
}

// writeRaw sends a low-level Engine.IO control frame (handshake replies) —
// never user-authored, so it isn't echoed to the transcript as "out".
func (s *SocketService) writeRaw(text string) {
	s.mu.Lock()
	conn := s.conn
	s.mu.Unlock()
	if conn == nil {
		return
	}
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	_ = conn.WriteMessage(websocket.TextMessage, []byte(text))
}

// Send writes a raw text frame — the plain-WebSocket send.
func (s *SocketService) Send(text string) error {
	s.mu.Lock()
	conn := s.conn
	s.mu.Unlock()
	if conn == nil {
		return errs.InvalidInput("not connected")
	}

	s.writeMu.Lock()
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	s.writeMu.Unlock()
	if err != nil {
		return errs.Wrap(err, errs.KindExternal, "failed to send")
	}
	s.emit("out", text, "")
	return nil
}

// Emit sends a Socket.IO event as `42["name", payload]` — payload is parsed
// as JSON when possible, else sent as a plain string, matching `reqx sio`'s
// REPL `emit` command.
func (s *SocketService) Emit(eventName string, payload string) error {
	s.mu.Lock()
	conn := s.conn
	protocol := s.protocol
	ready := s.sioReady
	s.mu.Unlock()
	if conn == nil {
		return errs.InvalidInput("not connected")
	}
	if protocol != "sio" {
		return errs.InvalidInput("emit is only available for Socket.IO connections")
	}
	if !ready {
		return errs.InvalidInput("cannot emit — waiting for the Socket.IO handshake to complete")
	}

	var payloadVal interface{} = ""
	if trimmed := strings.TrimSpace(payload); trimmed != "" {
		if err := json.Unmarshal([]byte(trimmed), &payloadVal); err != nil {
			payloadVal = payload
		}
	}
	packetBytes, err := json.Marshal([]interface{}{eventName, payloadVal})
	if err != nil {
		return errs.Wrap(err, errs.KindInternal, "failed to encode event")
	}
	text := "42" + string(packetBytes)

	s.writeMu.Lock()
	err = conn.WriteMessage(websocket.TextMessage, []byte(text))
	s.writeMu.Unlock()
	if err != nil {
		return errs.Wrap(err, errs.KindExternal, "failed to emit")
	}

	payloadStr, _ := json.Marshal(payloadVal)
	s.emit("out", string(payloadStr), eventName)
	return nil
}

// Disconnect closes the active connection, if any — a no-op when already
// disconnected so the frontend can call it freely on cleanup.
func (s *SocketService) Disconnect() error {
	s.mu.Lock()
	conn := s.conn
	s.conn = nil
	s.sioReady = false
	s.mu.Unlock()
	if conn == nil {
		return nil
	}
	return conn.Close()
}

// IsConnected reports whether a connection is currently open.
func (s *SocketService) IsConnected() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.conn != nil
}
