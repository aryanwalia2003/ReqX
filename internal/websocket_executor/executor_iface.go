package websocket_executor

import "reqx/internal/collection"

type WebSocketExecutor interface{
	Execute(url string, headers map[string]string, events []collection.WebSocketEvent, readyChan chan error, stopChan chan struct{})error
}

//websocket executor interface defines the contract for executing websocket requests, it says websocket request will have a url, headers, events, ready channel and stop channel

//events are of the type collection.WebSocketEvent, meaning they will look like this:
// {
// 	"type": "emit",
// 	"payload": "hello",
// 	"count": 1
// }