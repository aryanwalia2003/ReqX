package websocket_executor

import "time"

func NewDefaultWebSocketExecutor() WebSocketExecutor{
	return &defaultWebSocketExecutor{
		timeout: 10 * time.Second, //this timeout is for the connection between the client and the server
	}
}

//constructor of defaultWebSocketExecutor