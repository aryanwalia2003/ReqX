package services

import "testing"

// A nil db (history.Open failed at startup) must degrade to empty results,
// never a panic or an error — same graceful-skip the CLI does.
func TestHistoryService_NilDB(t *testing.T) {
	s := NewHistoryService(nil)

	runs, err := s.ListRuns(10)
	if err != nil || len(runs) != 0 {
		t.Errorf("ListRuns() = %+v, %v; want empty slice, nil error", runs, err)
	}

	stats, err := s.GetRunStats("any-id")
	if err != nil || len(stats) != 0 {
		t.Errorf("GetRunStats() = %+v, %v; want empty slice, nil error", stats, err)
	}

	nodes, err := s.GetDAGNodes("any-id")
	if err != nil || len(nodes) != 0 {
		t.Errorf("GetDAGNodes() = %+v, %v; want empty slice, nil error", nodes, err)
	}
}
