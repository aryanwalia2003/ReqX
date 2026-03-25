package socketio_executor

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"reqx/internal/collection"

	"github.com/gorilla/websocket"
)

// TestSocketIOExecutor_DataRace verifies there is no race on expectedListeners
// between the background reader (mutates under mu.Lock) and the main goroutine
// (reads under mu.Lock at the bottom of Execute).
// Run with: go test ./internal/socketio_executor/... -race
func TestSocketIOExecutor_DataRace(t *testing.T) {
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()

		// 1. Engine.IO Open packet
		c.WriteMessage(websocket.TextMessage, []byte(`0{"sid":"123","upgrades":[],"pingInterval":25000,"pingTimeout":5000}`))

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				break
			}
			text := string(msg)
			if strings.HasPrefix(text, "40") {
				c.WriteMessage(websocket.TextMessage, []byte("40")) // Socket.IO connected

				// Dispatch event quickly to race with the main goroutine read
				time.Sleep(10 * time.Millisecond)
				c.WriteMessage(websocket.TextMessage, []byte(`42["testEvent",{"data":"value"}]`))
			}
		}
	}))
	defer s.Close()

	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")
	exec := NewDefaultSocketIOExecutor()
	if e, ok := exec.(interface{ SetQuiet(bool) }); ok {
		e.SetQuiet(true)
	}

	events := []collection.SocketIOEvent{
		{Type: "listen", Name: "testEvent"},
	}
	_ = exec.Execute(wsURL, nil, events, nil, nil)
}

// TestSocketIOExecutor_Timeout validates that Execute returns cleanly when
// the expected event never arrives within the timeout.
func TestSocketIOExecutor_Timeout(t *testing.T) {
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()

		c.WriteMessage(websocket.TextMessage, []byte(`0{"sid":"123","upgrades":[],"pingInterval":25000,"pingTimeout":5000}`))

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				break
			}
			if strings.HasPrefix(string(msg), "40") {
				c.WriteMessage(websocket.TextMessage, []byte("40"))
			}
		}
	}))
	defer s.Close()

	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")
	exec := NewDefaultSocketIOExecutor()
	if e, ok := exec.(interface{ SetQuiet(bool) }); ok {
		e.SetQuiet(true)
	}

	events := []collection.SocketIOEvent{
		{Type: "listen", Name: "neverHappens"},
	}
	_ = exec.Execute(wsURL, nil, events, nil, nil)
}

// TestSocketIOExecutor_ConcurrentWriteMutex verifies that the background reader
// goroutine (writing heartbeat pongs "3") and the main goroutine (writing emits "42[...]")
// never call conn.WriteMessage concurrently.
//
// This directly exercises the writeMu fix. The server rapidly sends Engine.IO pings
// ("2") while the client is in the middle of its emit loop.
func TestSocketIOExecutor_ConcurrentWriteMutex(t *testing.T) {
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()

		c.WriteMessage(websocket.TextMessage, []byte(`0{"sid":"abc","upgrades":[],"pingInterval":25000,"pingTimeout":5000}`))

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				break
			}
			text := string(msg)

			if strings.HasPrefix(text, "40") {
				c.WriteMessage(websocket.TextMessage, []byte("40"))

				// Rapidly fire Engine.IO pings while client is emitting
				go func() {
					for i := 0; i < 15; i++ {
						c.WriteMessage(websocket.TextMessage, []byte("2"))
						time.Sleep(10 * time.Millisecond)
					}
				}()
			}
		}
	}))
	defer s.Close()

	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")
	exec := NewDefaultSocketIOExecutor()
	if e, ok := exec.(interface{ SetQuiet(bool) }); ok {
		e.SetQuiet(true)
	}

	// Emit several events while pings are incoming — this is the exact race scenario
	events := []collection.SocketIOEvent{
		{Type: "emit", Name: "action", Payload: `{"step":1}`},
		{Type: "emit", Name: "action", Payload: `{"step":2}`},
		{Type: "emit", Name: "action", Payload: `{"step":3}`},
		{Type: "emit", Name: "action", Payload: `{"step":4}`},
		{Type: "emit", Name: "action", Payload: `{"step":5}`},
	}
	_ = exec.Execute(wsURL, nil, events, nil, nil)
}

// TestSocketIOExecutor_AsyncMode verifies that Execute correctly handles
// the stopChan signal in async mode without any race on the channel close.
func TestSocketIOExecutor_AsyncMode(t *testing.T) {
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()

		c.WriteMessage(websocket.TextMessage, []byte(`0{"sid":"xyz","upgrades":[],"pingInterval":25000,"pingTimeout":5000}`))

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				break
			}
			if strings.HasPrefix(string(msg), "40") {
				c.WriteMessage(websocket.TextMessage, []byte("40"))
			}
		}
	}))
	defer s.Close()

	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")
	exec := NewDefaultSocketIOExecutor()
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

	close(stopChan)
	time.Sleep(200 * time.Millisecond)
}
