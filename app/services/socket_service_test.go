package services

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var testUpgrader = websocket.Upgrader{}

func wsURLOf(srv *httptest.Server) string {
	return "ws" + strings.TrimPrefix(srv.URL, "http")
}

func TestSocketService_ConnectSendDisconnect_WS(t *testing.T) {
	received := make(chan string, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		_, msg, err := conn.ReadMessage()
		if err == nil {
			received <- string(msg)
		}
	}))
	defer srv.Close()

	s := NewSocketService()
	if err := s.Connect(ConnectSocketInput{URL: wsURLOf(srv), Protocol: "ws"}); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	if !s.IsConnected() {
		t.Fatal("expected IsConnected() true after Connect")
	}

	if err := s.Send("hello"); err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	select {
	case got := <-received:
		if got != "hello" {
			t.Errorf("server received %q, want %q", got, "hello")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for server to receive message")
	}

	if err := s.Disconnect(); err != nil {
		t.Fatalf("Disconnect() error = %v", err)
	}
	if s.IsConnected() {
		t.Fatal("expected IsConnected() false after Disconnect")
	}
	if err := s.Send("late"); err == nil {
		t.Fatal("expected Send() after Disconnect to fail")
	}
	// Disconnect must be safe to call again (frontend cleanup calls it freely).
	if err := s.Disconnect(); err != nil {
		t.Fatalf("second Disconnect() error = %v", err)
	}
}

func TestSocketService_Connect_RejectsWhenAlreadyConnected(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		time.Sleep(2 * time.Second)
	}))
	defer srv.Close()

	s := NewSocketService()
	if err := s.Connect(ConnectSocketInput{URL: wsURLOf(srv), Protocol: "ws"}); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer s.Disconnect()

	if err := s.Connect(ConnectSocketInput{URL: wsURLOf(srv), Protocol: "ws"}); err == nil {
		t.Fatal("expected second Connect() to fail while already connected")
	}
}

func TestSocketService_Connect_RejectsBadWSURL(t *testing.T) {
	s := NewSocketService()
	if err := s.Connect(ConnectSocketInput{URL: "http://example.com", Protocol: "ws"}); err == nil {
		t.Fatal("expected error for a non-ws:// URL with protocol \"ws\"")
	}
}

func TestSocketService_Emit_RejectsForPlainWS(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		time.Sleep(2 * time.Second)
	}))
	defer srv.Close()

	s := NewSocketService()
	if err := s.Connect(ConnectSocketInput{URL: wsURLOf(srv), Protocol: "ws"}); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer s.Disconnect()

	if err := s.Emit("evt", "{}"); err == nil {
		t.Fatal("expected Emit() to fail on a plain WebSocket connection")
	}
}

func TestSocketService_Emit_RejectsBeforeSocketIOHandshake(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		// Never sends the Engine.IO "0" open frame, so the client's
		// handshake never completes.
		time.Sleep(2 * time.Second)
	}))
	defer srv.Close()

	s := NewSocketService()
	if err := s.Connect(ConnectSocketInput{URL: srv.URL, Protocol: "sio"}); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer s.Disconnect()

	if err := s.Emit("evt", "{}"); err == nil {
		t.Fatal("expected Emit() to fail before the Socket.IO handshake completes")
	}
}

// TestSocketService_Emit_SocketIOHandshakeAndFrame drives a minimal
// Engine.IO v4 handshake (server sends "0", client auto-replies "40") then
// emits an event and asserts the server sees the exact `42[...]` frame.
func TestSocketService_Emit_SocketIOHandshakeAndFrame(t *testing.T) {
	frameCh := make(chan string, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		if err := conn.WriteMessage(websocket.TextMessage, []byte("0{}")); err != nil {
			return
		}
		// Client should auto-reply "40" to the Engine.IO open frame.
		if _, msg, err := conn.ReadMessage(); err != nil || string(msg) != "40" {
			return
		}
		// Server confirms the Socket.IO-level connect — this is what flips
		// the client's sioReady gate.
		if err := conn.WriteMessage(websocket.TextMessage, []byte("40{}")); err != nil {
			return
		}

		_, msg, err := conn.ReadMessage()
		if err == nil {
			frameCh <- string(msg)
		}
	}))
	defer srv.Close()

	s := NewSocketService()
	if err := s.Connect(ConnectSocketInput{URL: srv.URL, Protocol: "sio"}); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer s.Disconnect()

	// Wait for the background readLoop to complete the handshake.
	deadline := time.After(2 * time.Second)
	for {
		s.mu.Lock()
		ready := s.sioReady
		s.mu.Unlock()
		if ready {
			break
		}
		select {
		case <-deadline:
			t.Fatal("timed out waiting for Socket.IO handshake")
		case <-time.After(10 * time.Millisecond):
		}
	}

	if err := s.Emit("greet", `{"name":"dev"}`); err != nil {
		t.Fatalf("Emit() error = %v", err)
	}

	select {
	case got := <-frameCh:
		want := `42["greet",{"name":"dev"}]`
		if got != want {
			t.Errorf("server received frame %q, want %q", got, want)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for server to receive the emit frame")
	}
}

func TestSocketDialURL_SIORewritesScheme(t *testing.T) {
	got, err := socketDialURL("http://localhost:3000", "sio")
	if err != nil {
		t.Fatalf("socketDialURL() error = %v", err)
	}
	if !strings.HasPrefix(got, "ws://localhost:3000/socket.io/?") {
		t.Errorf("socketDialURL() = %q, want ws://localhost:3000/socket.io/?...", got)
	}
	if !strings.Contains(got, "EIO=4") || !strings.Contains(got, "transport=websocket") {
		t.Errorf("socketDialURL() = %q, missing EIO/transport query params", got)
	}
}
