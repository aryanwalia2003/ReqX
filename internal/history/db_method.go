package history

import (
	"fmt"
	"reqx/internal/metrics"
	"time"
)

// SaveRun persists an aggregated test report in a single transaction.
// Called once after the test completes — no concurrent writes.
func (d *DB) SaveRun(collection string, r metrics.Report) error {
	id := fmt.Sprintf("%d", time.Now().UnixNano())

	tx, err := d.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	errorPct := 0.0
	if r.TotalRequests > 0 {
		errorPct = float64(r.TotalFailures) / float64(r.TotalRequests) * 100
	}

	_, err = tx.Exec(
		`INSERT INTO test_runs(id,collection,total_reqs,rps,p95_ms,error_pct) VALUES(?,?,?,?,?,?)`,
		id, collection, r.TotalRequests, r.RPS, r.P95.Milliseconds(), errorPct,
	)
	if err != nil {
		return err
	}

	for _, s := range r.PerRequest {
		_, err = tx.Exec(
			`INSERT INTO request_stats(run_id,name,successes,failures,p95_ms,avg_ms) VALUES(?,?,?,?,?,?)`,
			id, s.Name, s.Successes, s.Failures, s.P95.Milliseconds(), s.AvgDuration.Milliseconds(),
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
