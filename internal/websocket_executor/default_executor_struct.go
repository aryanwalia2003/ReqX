package websocket_executor

import "time"

type defaultWebSocketExecutor struct{
	timeout time.Duration
} 

//defaultWebSocketExecutor implements the WebSocketExecutor interface