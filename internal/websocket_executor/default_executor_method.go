package websocket_executor

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"reqx/internal/collection"
	"reqx/internal/errs"

	"github.com/gorilla/websocket"
)

// Execute handles the raw WebSocket connection lifecycle.
func (e *defaultWebSocketExecutor) Execute(rawURL string, headers map[string]string, events []collection.WebSocketEvent, readyChan chan error, stopChan chan struct{}) error {
	if rawURL == "" {
		err := errs.InvalidInput("invalid websocket url: empty")
		if readyChan != nil {
			readyChan <- err
		}
		return err
	}

	// Ensure URL has ws:// or wss:// scheme
	if !strings.HasPrefix(rawURL, "ws://") && !strings.HasPrefix(rawURL, "wss://") {
		err := errs.InvalidInput("WebSocket URL must start with ws:// or wss://")
		if readyChan != nil {
			readyChan <- err
		}
		return err
	}

	reqHeaders := http.Header{}
	for k, v := range headers {
		reqHeaders.Add(k, v)
	}

	if !e.quiet {
		fmt.Printf("Connecting to WebSocket Server: %s\n", rawURL)
	}
	conn, _, err := sharedDialer.Dial(rawURL, reqHeaders)
	if err != nil {
		if readyChan != nil {
			readyChan <- err
		}
		return errs.Wrap(err, errs.KindInternal, "Failed to connect to websocket")
	}
	defer conn.Close()

	// Signal that the connection is ready for async flows
	if readyChan != nil {
		readyChan <- nil
	}
	if !e.quiet {
		fmt.Println("Connected successfully.")
	}

	// Set handlers for better debugging (only in verbose mode)
	if !e.quiet {
		conn.SetCloseHandler(func(code int, text string) error {
			fmt.Printf("\x1b[31m\n[WS_CLOSE_FRAME] Code: %d, Message: %s\x1b[0m\n", code, text)
			return nil
		})
		conn.SetPingHandler(func(appData string) error {
			fmt.Printf("\x1b[36m\n[WS_PING] Received (%s)\x1b[0m\n", appData)
			return conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
		})
		conn.SetPongHandler(func(appData string) error {
			fmt.Printf("\x1b[36m\n[WS_PONG] Received (%s)\x1b[0m\n", appData)
			return nil
		})
	} else {
		// In quiet mode, still respond to pings silently
		conn.SetPingHandler(func(appData string) error {
			return conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
		})
	}

	// Heartbeat (Keep-alive)
	heartbeatStop := make(chan struct{})
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := conn.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(time.Second)); err != nil {
					return
				}
			case <-heartbeatStop:
				return
			}
		}
	}()
	defer close(heartbeatStop)

	var mu sync.Mutex
	expectedMessages := 0
	for _, ev := range events {
		if ev.Type == "listen" {
			expectedMessages++
		}
	}

	done := make(chan struct{})

	// Background Reader
	go func() {
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				// Check if this was an intentional close from our side
				select {
				case <-stopChan:
					// Expected closure, don't log as error
					return
				default:
					if !e.quiet {
						if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
							fmt.Printf("\x1b[31m\n[WS_ERROR] Abnormal Disconnect (1006?): %v\x1b[0m\n", err)
						} else {
							fmt.Printf("\n[READER INFO] WebSocket closed gracefully: %v\n", err)
						}
					}
					return
				}
			}

			mu.Lock()
			if !e.quiet {
				fmt.Printf("\x1b[35m\n[WS_INCOMING] (%d bytes, type %d): %s\x1b[0m\n", len(message), mt, string(message))
			}

			// Only decrement counters if we are in Synchronous Mode
			if stopChan == nil {
				if expectedMessages > 0 {
					expectedMessages--
					if expectedMessages == 0 {
						select {
						case <-done:
						default:
							close(done)
						}
					}
				}
			}
			mu.Unlock()
		}
	}()

	// Emit predefined messages
	for _, ev := range events {
		if ev.Type == "emit" {
			if !e.quiet {
				fmt.Printf("[EMIT] Payload: %s\n", ev.Payload)
			}
			err := conn.WriteMessage(websocket.TextMessage, []byte(ev.Payload))
			if err != nil {
				return errs.Wrap(err, errs.KindInternal, "Failed to write message")
			}
			time.Sleep(200 * time.Millisecond)
		}
	}

	// WAIT LOGIC (Async vs Sync)
	if stopChan != nil {
		// ASYNC MODE: Wait indefinitely until Runner sends stop signal
		<-stopChan
		if !e.quiet {
			fmt.Println("\nClosing Background WebSocket connection...")
		}
		return nil
	}

	// SYNC MODE: Wait for a specific number of messages to arrive
	mu.Lock()
	remaining := expectedMessages
	mu.Unlock()

	if remaining > 0 {
		if !e.quiet {
			fmt.Printf("Waiting up to %v for %d message(s)...\n", e.timeout, remaining)
		}
		select {
		case <-done:
			if !e.quiet {
				fmt.Println("All expected messages received.")
			}
		case <-time.After(e.timeout):
			if !e.quiet {
				mu.Lock()
				missed := expectedMessages
				mu.Unlock()
				fmt.Printf("Timeout reached! Missed %d message(s).\n", missed)
			}
		}
	} else if stopChan == nil {
		time.Sleep(1 * time.Second)
	}

	if !e.quiet {
		fmt.Println("Closing WebSocket connection.")
	}
	return nil
}
