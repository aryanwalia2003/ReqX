package ui

import "reqx/internal/history"

// NewServer creates a UI server bound to the given port.
func NewServer(db *history.DB, port int) *Server {
	return &Server{db: db, port: port}
}
