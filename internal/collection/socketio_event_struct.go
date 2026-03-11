package collection

// SocketIOEvent represents an event to either emit or listen for via Socket.IO.
type SocketIOEvent struct {
	Type    string `json:"type"`    // "emit" or "listen"
	Name    string `json:"name"`    // The name of the event
	Payload string `json:"payload"` // The data to send (if emitting)
}
