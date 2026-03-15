package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"reqx/internal/errs"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

func NewWSCmd() *cobra.Command {
	var headers []string
	
	c := &cobra.Command{
		Use:   "ws [url]",
		Short: "Start an interactive raw WebSocket session",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rawURL := args[0]
			
			if !strings.HasPrefix(rawURL, "ws://") && !strings.HasPrefix(rawURL, "wss://") {
				return errs.InvalidInput("URL must start with ws:// or wss://")
			}

			reqHeaders := http.Header{}
			for _, h := range headers {
				parts := strings.SplitN(h, ":", 2)
				if len(parts) == 2 {
					reqHeaders.Add(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
				}
			}

			conn, _, err := websocket.DefaultDialer.Dial(rawURL, reqHeaders)
			if err != nil {
				return errs.Wrap(err, errs.KindInternal, "Failed to connect")
			}
			defer conn.Close()

			color.Green("✅ Connected! Type a message and press Enter to send. Press Ctrl+C to quit.")
			
			// Background listener
			go func() {
				// Set handlers for better debugging
				conn.SetCloseHandler(func(code int, text string) error {
					color.Red("\n❌ [WS_CLOSE_FRAME] Code: %d, Message: %s", code, text)
					return nil
				})

				conn.SetPingHandler(func(appData string) error {
					fmt.Printf("\x1b[36m\n[WS_PING] Received (%s)\x1b[0m\n> ", appData)
					return conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
				})

				conn.SetPongHandler(func(appData string) error {
					fmt.Printf("\x1b[36m\n[WS_PONG] Received (%s)\x1b[0m\n> ", appData)
					return nil
				})

				for {
					mt, message, err := conn.ReadMessage()
					if err != nil {
						if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
							color.Red("\n❌ Abnormal Disconnect: %v", err)
							fmt.Printf("\x1b[33m[DEBUG] Error Type: %T, Value: %+v\x1b[0m\n", err, err)
						} else {
							color.Yellow("\nℹ️  Connection closed: %v", err)
						}
						os.Exit(0)
					}
					color.Cyan("\n⬇️  [RECEIVED] (type %d): %s\n> ", mt, string(message))
				}
			}()

			// Heartbeat goroutine
			go func() {
				ticker := time.NewTicker(30 * time.Second)
				defer ticker.Stop()
				for {
					<-ticker.C
					if err := conn.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(time.Second)); err != nil {
						return
					}
				}
			}()

			// Main thread for sending
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Print("> ")
			for scanner.Scan() {
				text := scanner.Text()
				err := conn.WriteMessage(websocket.TextMessage, []byte(text))
				if err != nil {
					return errs.Wrap(err, errs.KindInternal, "Failed to send message")
				}
				color.Yellow("⬆️  [SENT]: %s", text)
				fmt.Print("> ")
			}
			
			return nil
		},
	}
	
	c.Flags().StringSliceVarP(&headers, "header", "H", []string{}, "Custom headers")
	return c
}