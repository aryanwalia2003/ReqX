package ui

import "reqx/internal/history"

// Server serves the embedded dashboard and history API.
type Server struct {
	db   *history.DB
	port int
}
