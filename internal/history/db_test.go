package history

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/glebarez/go-sqlite"

	"reqx/internal/metrics"
)

// openTestDB opens a DB against a temp file instead of Open()'s hardcoded
// ~/.reqx/history.db, so tests never touch a real user's history.
func openTestDB(t *testing.T) *DB {
	t.Helper()
	conn, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "history.db"))
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	if _, err := conn.Exec(schema); err != nil {
		t.Fatalf("schema exec error = %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return &DB{conn: conn}
}

func TestDB_SaveRunAndListRuns(t *testing.T) {
	db := openTestDB(t)

	report := metrics.Report{
		TotalRequests: 2,
		TotalSuccess:  1,
		TotalFailures: 1,
		RPS:           10,
		P95:           50 * time.Millisecond,
		PerRequest: []metrics.RequestStat{
			{Name: "Get user", Successes: 1, Failures: 0, P95: 20 * time.Millisecond, AvgDuration: 15 * time.Millisecond},
		},
	}

	if err := db.SaveRun("demo.json", report); err != nil {
		t.Fatalf("SaveRun() error = %v", err)
	}

	runs, err := db.ListRuns(10)
	if err != nil {
		t.Fatalf("ListRuns() error = %v", err)
	}
	if len(runs) != 1 {
		t.Fatalf("len(runs) = %d, want 1", len(runs))
	}
	got := runs[0]
	if got.Collection != "demo.json" || got.TotalReqs != 2 || got.P95Ms != 50 {
		t.Errorf("ListRuns()[0] = %+v, unexpected values", got)
	}
	if got.ErrorPct != 50 {
		t.Errorf("ErrorPct = %v, want 50 (1 failure of 2 requests)", got.ErrorPct)
	}

	stats, err := db.GetRunStats(got.ID)
	if err != nil {
		t.Fatalf("GetRunStats() error = %v", err)
	}
	if len(stats) != 1 || stats[0].Name != "Get user" || stats[0].P95Ms != 20 {
		t.Errorf("GetRunStats() = %+v, want 1 stat for 'Get user' with P95Ms=20", stats)
	}
}

func TestDB_ListRuns_RespectsLimit(t *testing.T) {
	db := openTestDB(t)

	for i := 0; i < 5; i++ {
		if err := db.SaveRun("c.json", metrics.Report{}); err != nil {
			t.Fatalf("SaveRun() #%d error = %v", i, err)
		}
	}

	runs, err := db.ListRuns(3)
	if err != nil {
		t.Fatalf("ListRuns() error = %v", err)
	}
	if len(runs) != 3 {
		t.Errorf("len(runs) = %d, want 3 (limit)", len(runs))
	}
}

func TestDB_GetRunStats_UnknownRun(t *testing.T) {
	db := openTestDB(t)

	stats, err := db.GetRunStats("does-not-exist")
	if err != nil {
		t.Fatalf("GetRunStats() error = %v", err)
	}
	if len(stats) != 0 {
		t.Errorf("GetRunStats() for unknown run = %+v, want empty", stats)
	}
}
