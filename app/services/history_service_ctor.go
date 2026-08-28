package services

import "reqx/internal/history"

// NewHistoryService constructs the service. db may be nil (history.Open
// failed at startup) — every method below handles that gracefully.
func NewHistoryService(db *history.DB) *HistoryService {
	return &HistoryService{db: db}
}
