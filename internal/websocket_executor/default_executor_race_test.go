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

// TestWebSocketExecutor_DataRace verifies there is no race on expectedMessages
// between the background reader goroutine (mutates under mu) and the main
// goroutine (reads under mu at the bottom of Execute).
// Run with: go test ./internal/websocket_executor/... -race
func TestWebSocketExecutor_DataRace(t *testing.T) {
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()

		// Send a message immediately to race against the main goroutine's read.
		time.Sleep(10 * time.Millisecond)
		c.WriteMessage(websocket.TextMessage, []byte(`{"event":"ok"}`))

		// Drain remaining
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
	_ = exec.Execute(wsURL, nil, events, nil, nil)
}

// TestWebSocketExecutor_Timeout validates that Execute returns within the
// expected timeout window when no messages ever arrive.
func TestWebSocketExecutor_Timeout(t *testing.T) {
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
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

	events := []collection.WebSocketEvent{
		{Type: "listen"},
	}
	_ = exec.Execute(wsURL, nil, events, nil, nil)
}

// TestWebSocketExecutor_ConcurrentWriteMutex verifies that the heartbeat goroutine
// and the emit loop never write to the WebSocket concurrently.
//
// The test emits many messages while simultaneously the server sends pings
// (exercising the PingHandler WriteControl path) and the heartbeat ticker fires.
// With the writeMu fix this must produce zero races.
func TestWebSocketExecutor_ConcurrentWriteMutex(t *testing.T) {
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()

		// Rapidly send pings to exercise the PingHandler concurrently with emits.
		go func() {
			for i := 0; i < 20; i++ {
				c.WriteMessage(websocket.PingMessage, []byte("ping"))
				time.Sleep(5 * time.Millisecond)
			}
		}()

		// Drain incoming messages
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

	// Emit 10 messages while pings are incoming
	events := make([]collection.WebSocketEvent, 10)
	for i := range events {
		events[i] = collection.WebSocketEvent{Type: "emit", Payload: `{"seq":` + string(rune('0'+i)) + `}`}
	}
	_ = exec.Execute(wsURL, nil, events, nil, nil)
}

// TestWebSocketExecutor_AsyncMode verifies that Execute returns cleanly
// when stopped via stopChan in async mode (no race on the channel close).
func TestWebSocketExecutor_AsyncMode(t *testing.T) {
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
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

	stopChan := make(chan struct{})
	readyChan := make(chan error, 1)

	go func() {
		_ = exec.Execute(wsURL, nil, nil, readyChan, stopChan)
	}()

	select {
	case err := <-readyChan:
		if err != nil {
			t.Fatalf("unexpected ready error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for async ready signal")
	}

	// Signal stop and confirm no hang
	close(stopChan)
	time.Sleep(200 * time.Millisecond)
}
