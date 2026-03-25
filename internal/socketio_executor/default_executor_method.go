package socketio_executor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"

	"reqx/internal/collection"
	"reqx/internal/errs"
)

// Execute runs the Socket.IO flow, emitting and listening to defined events using raw V4 WebSockets.
func (e *DefaultSocketIOExecutor) Execute(rawURL string, headers map[string]string, events []collection.SocketIOEvent, readyChan chan error, stopChan chan struct{}) error {
	if rawURL == "" {
		if readyChan != nil {
			readyChan <- errs.InvalidInput("invalid socket.io url: empty")
		}
		return errs.InvalidInput("invalid socket.io url: empty")
	}

	// 1. Format URL for WebSocket and Engine.IO v4
	u, err := url.Parse(rawURL)
	if err != nil {
		if readyChan != nil {
			readyChan <- errs.Wrap(err, errs.KindInvalidInput, "Invalid URL format")
		}
		return errs.Wrap(err, errs.KindInvalidInput, "Invalid URL format")
	}
	if u.Scheme == "http" {
		u.Scheme = "ws"
	} else if u.Scheme == "https" {
		u.Scheme = "wss"
	}
	if u.Path == "" || u.Path == "/" {
		u.Path = "/socket.io/"
	}

	q := u.Query()
	q.Set("EIO", "4") // FORCE SOCKET.IO v4
	q.Set("transport", "websocket")
	u.RawQuery = q.Encode()

	// 2. Prepare Custom Headers (e.g., Cookies, Authorization)
	reqHeaders := http.Header{}
	for k, v := range headers {
		reqHeaders.Add(k, v)
	}

	if !e.quiet {
		fmt.Printf("Connecting to Socket.IO Server (v4): %s\n", u.String())
	}

	// 3. Connect via raw WebSocket
	conn, _, err := sharedDialer.Dial(u.String(), reqHeaders)
	if err != nil {
		if readyChan != nil {
			readyChan <- errs.Wrap(err, errs.KindInternal, "Failed to connect to websocket")
		}
		return errs.Wrap(err, errs.KindInternal, "Failed to connect to websocket")
	}
	defer conn.Close()

	// 4. State Management for Listeners
	var mu sync.Mutex    // Guards expectedListeners and listenTargets
	var writeMu sync.Mutex // Guards conn.WriteMessage (gorilla forbids concurrent writers)
	expectedListeners := 0
	listenTargets := make(map[string]int)

	for _, ev := range events {
		if ev.Type == "listen" {
			expectedListeners++
			listenTargets[ev.Name]++
			if !e.quiet {
				fmt.Printf("Registered listener for event: '%s'\n", ev.Name)
			}
		}
	}

	done := make(chan struct{})
	connected := make(chan struct{}) // To ensure we wait for '40' before emitting

	// 5. Background Reader (Handles Protocol & Incoming Events)
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return
			}

			msgStr := string(message)

			if strings.HasPrefix(msgStr, "0") {
				// Engine.IO Open -> Send Socket.IO Connect (40)
				writeMu.Lock()
				conn.WriteMessage(websocket.TextMessage, []byte("40"))
				writeMu.Unlock()
			} else if strings.HasPrefix(msgStr, "2") {
				// Engine.IO Ping -> Reply with Pong (3)
				writeMu.Lock()
				conn.WriteMessage(websocket.TextMessage, []byte("3"))
				writeMu.Unlock()
			} else if strings.HasPrefix(msgStr, "40") {
				// Socket.IO Connected
				if !e.quiet {
					fmt.Println("Connected successfully.")
				}
				select {
				case <-connected:
				default:
					close(connected) // Signal that it's safe to emit
					if readyChan != nil {
						readyChan <- nil
					}
				}
			} else if strings.HasPrefix(msgStr, "42") {
				// Incoming Event
				dataStr := msgStr[2:]
				
				// Use gjson to sniff the event name without unmarshalling the entire array
				res := gjson.Get(dataStr, "0")
				if res.Exists() {
					eventName := res.String()
					mu.Lock()
					isListening := false

					for _, ev := range events {
						if ev.Type == "listen" && ev.Name == eventName {
							isListening = true
							break
						}
					}

					if isListening {
						if !e.quiet {
							// Extract payload as raw string from the array's second element
							payload := gjson.Get(dataStr, "1").Raw
							fmt.Printf("\n[RECEIVED] Event: '%s' | Data: %v\n", eventName, payload)
						}

						// Only decrement and track target counts if we are in Synchronous mode
						if stopChan == nil {
							if needed := listenTargets[eventName]; needed > 0 {
								listenTargets[eventName]--
								expectedListeners--
								if expectedListeners == 0 {
									select {
									case <-done:
									default:
										close(done)
									}
								}
							}
						}
					}
					mu.Unlock()
				}
			}
		}
	}()

	// Wait up to 5 seconds for the Socket.IO connection to fully establish
	select {
	case <-connected:
	case <-time.After(5 * time.Second):
		err := errs.Internal("Timeout waiting for Socket.IO connect (40) packet")
		if readyChan != nil {
			readyChan <- err
		}
		return err
	}

	// 6. Emit predefined events
	for _, ev := range events {
		if ev.Type == "emit" {
			if !e.quiet {
				fmt.Printf("[EMIT] Event: '%s' | Payload: %s\n", ev.Name, ev.Payload)
			}

			nameBytes, _ := json.Marshal(ev.Name)
			var finalMessage string

			if ev.Payload == "" {
				finalMessage = "42[" + string(nameBytes) + "]"
			} else if gjson.Valid(ev.Payload) {
				// If the payload is perfectly valid JSON (object, array, number, quoted string, boolean, null),
				// embed it directly.
				finalMessage = "42[" + string(nameBytes) + "," + ev.Payload + "]"
			} else {
				// Otherwise, it was just unquoted text; JSON encode it as a simple string.
				payloadBytes, _ := json.Marshal(ev.Payload)
				finalMessage = "42[" + string(nameBytes) + "," + string(payloadBytes) + "]"
			}

			writeMu.Lock()
			conn.WriteMessage(websocket.TextMessage, []byte(finalMessage))
			writeMu.Unlock()
			time.Sleep(10 * time.Millisecond) // Slight delay between emits to preserve ordering without hanging VUs
		}
	}

	// ========================================================
	// 7. WAIT LOGIC (Async vs Sync)
	// ========================================================

	// ASYNC MODE: Wait indefinitely until Runner sends stop signal
	if stopChan != nil {
		<-stopChan
		if !e.quiet {
			fmt.Println("\nClosing Background Socket.IO connection...")
		}
		return nil
	}

	// SYNC MODE: Wait for specific events to arrive
	mu.Lock()
	remaining := expectedListeners
	mu.Unlock()

	if remaining > 0 {
		if !e.quiet {
			fmt.Printf("Waiting up to %v for expected listener(s)...\n", e.timeout)
		}
		select {
		case <-done:
			if !e.quiet {
				fmt.Println("All expected events received.")
			}
		case <-time.After(e.timeout):
			if !e.quiet {
				mu.Lock()
				missed := expectedListeners
				mu.Unlock()
				fmt.Printf("Timeout reached! Missed %d event(s).\n", missed)
			}
		}
	} else {
		// Just wait a tiny bit to ensure final emits go out before closing the conn
		time.Sleep(1 * time.Second)
	}

	if !e.quiet {
		fmt.Println("Closing Socket.IO connection.")
	}
	return nil
}
