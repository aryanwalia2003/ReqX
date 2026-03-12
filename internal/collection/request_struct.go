package collection

// Request represents a single API or Socket call.
type Request struct {
	Name     string            `json:"name"`
	Method   string            `json:"method"`
	URL      string            `json:"url"`
	Protocol string            `json:"protocol,omitempty"` // "HTTP", "WS", "SOCKETIO"
	Headers  map[string]string `json:"headers,omitempty"`
	Body     string            `json:"body,omitempty"`
	Auth     *Auth             `json:"auth,omitempty"`
	Events   []SocketIOEvent   `json:"events,omitempty"`
	Scripts  []Script          `json:"scripts,omitempty"`
}

//this struct will look like 
// {
// 	"name": "request_name",
// 	"method": "GET",
// 	"url": "http://localhost:8080",
// 	"headers": {
// 		"Content-Type": "application/json"
// 	},
// 	"body": "{\"key\":\"value\"}",
// 	"scripts": [
// 		{
// 			"name": "script_name",
// 			"type": "pre",
// 			"content": "console.log('hello')"
// 		}
// 	]
// }
