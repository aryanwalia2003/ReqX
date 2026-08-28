package services

import (
	"reqx/internal/errs"
	"reqx/internal/history"
)

// ListRuns returns the most recent runs, newest first — same query as the
// CLI's `reqx ui` dashboard (internal/ui's apiHistory). limit<=0 defaults to
// 50. A nil db (history.Open failed at startup) degrades to an empty list,
// same as the CLI silently skipping history when it can't open the file.
func (s *HistoryService) ListRuns(limit int) ([]history.RunRow, error) {
	if s.db == nil {
		return []history.RunRow{}, nil
	}
	if limit <= 0 {
		limit = 50
	}
	runs, err := s.db.ListRuns(limit)
	if err != nil {
		return nil, errs.Wrap(err, errs.KindDatabase, "could not list run history")
	}
	if runs == nil {
		runs = []history.RunRow{}
	}
	return runs, nil
}

// GetRunStats returns the per-request breakdown for one run.
func (s *HistoryService) GetRunStats(runID string) ([]history.StatRow, error) {
	if s.db == nil {
		return []history.StatRow{}, nil
	}
	stats, err := s.db.GetRunStats(runID)
	if err != nil {
		return nil, errs.Wrap(err, errs.KindDatabase, "could not load run detail")
	}
	if stats == nil {
		stats = []history.StatRow{}
	}
	return stats, nil
}

// GetDAGNodes returns the scenario-graph nodes for one run (empty for a
// linear, non-DAG run).
func (s *HistoryService) GetDAGNodes(runID string) ([]history.DagNodeRow, error) {
	if s.db == nil {
		return []history.DagNodeRow{}, nil
	}
	nodes, err := s.db.GetDAGNodes(runID)
	if err != nil {
		return nil, errs.Wrap(err, errs.KindDatabase, "could not load run graph")
	}
	return nodes, nil
}
