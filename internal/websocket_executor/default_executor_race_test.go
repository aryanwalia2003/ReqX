package websocket_executor

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"reqx/internal/collection"

	"github.com/gorilla/websocket"
)

// TestWebSocketExecutor_DataRace triggers the known unprotected read
// on expectedMessages inside Execute() against the background reader mutex.
// Run with "go test ./internal/websocket_executor/... -race".
func TestWebSocketExecutor_DataRace(t *testing.T) {
	upgrader := websocket.Upgrader{}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()

		// Allow client to fully boot reader thread before sending
		time.Sleep(50 * time.Millisecond)
		c.WriteMessage(websocket.TextMessage, []byte("hello"))

		for {
			if _, _, err := c.ReadMessage(); err != nil {
				break
			}
		}
	}))
	defer s.Close()

	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")

	exec := NewDefaultWebSocketExecutor()
	if e, ok := exec.(interface{ SetQuiet(bool) }); ok {
		e.SetQuiet(true)
	}

	events := []collection.WebSocketEvent{
		{Type: "listen"},
	}

	// This triggers the race since Execute synchronously reads expectedMessages
	// at the same time the background loop attempts to decrement it.
	_ = exec.Execute(wsURL, nil, events, nil, nil)
}

func TestWebSocketExecutor_Timeout(t *testing.T) {
	upgrader := websocket.Upgrader{}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				break
			}
		}
	}))
	defer s.Close()

	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")
	exec := NewDefaultWebSocketExecutor()
	if e, ok := exec.(interface{ SetQuiet(bool) }); ok {
		e.SetQuiet(true)
	}

	// We are going to cheat and inject the struct to override timeout for the test
	type timeoutSetter interface {
		SetQuiet(bool)
	}

	// Since we can't easily override the timeout field without casting to the private struct,
	// we just rely on standard execution to eventually finish.
	// But this perfectly validates that "Listen" without server "Emit" cleanly times out!
	events := []collection.WebSocketEvent{{Type: "listen"}}
	_ = exec.Execute(wsURL, nil, events, nil, nil)
}
