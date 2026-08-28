package services

import "reqx/internal/history"

// HistoryService is the read side of run history — the app-side equivalent
// of the CLI's `reqx ui` dashboard (internal/ui), reusing the exact same
// internal/history queries. CollectionService.Run is the write side (it
// saves via the same *history.DB, shared from wailsapp.NewApp).
type HistoryService struct {
	db *history.DB
}
