package services

// NewSocketService constructs an idle debugger service — no connection until
// Connect is called from the frontend.
func NewSocketService() *SocketService {
	return &SocketService{}
}
