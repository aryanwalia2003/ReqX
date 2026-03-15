package collection

type WebSocketEvent struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
	Count   int    `json:"count"`
}

//websocket event looks like this:
// {
// 	"type": "emit",
// 	"payload": "hello",
// 	"count": 1
// }