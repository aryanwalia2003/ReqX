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

// TestSocketIOExecutor_DataRace triggers the known unprotected read
// on expectedListeners inside Execute() against the background reader mutex.
// Run with "go test ./internal/socketio_executor/... -race".
func TestSocketIOExecutor_DataRace(t *testing.T) {
	upgrader := websocket.Upgrader{}
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
			if strings.HasPrefix(text, "40") { // Client connected
				c.WriteMessage(websocket.TextMessage, []byte("40")) // Server connected

				// Dispatch an event quickly to decrement expectedListeners in the background
				// while the main goroutine evaluates `expectedListeners > 0`
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

func TestSocketIOExecutor_Timeout(t *testing.T) {
	upgrader := websocket.Upgrader{}
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
