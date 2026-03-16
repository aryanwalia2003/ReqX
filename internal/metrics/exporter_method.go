package metrics

import (
	"encoding/json"
	"os"
	"reqx/internal/errs"
	"reqx/internal/runner"
)

// exportRecord is the JSON shape written per request.
type exportRecord struct {
	Iteration  int    `json:"iteration"`
	WorkerID   int    `json:"worker_id"`
	Name       string `json:"name"`
	Protocol   string `json:"protocol"`
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	DurationMs int64  `json:"duration_ms"`
	Error      string `json:"error,omitempty"`
}

// ExportJSON writes all raw metrics to path as newline-delimited JSON.
// Each line is one complete JSON object, making it streamable and jq-friendly.
func ExportJSON(allMetrics [][]runner.RequestMetric, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return errs.Wrap(err, errs.KindInternal, "could not create export file")
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	for iterIdx, iterMetrics := range allMetrics {
		for _, m := range iterMetrics {
			rec := exportRecord{
				Iteration:  iterIdx + 1,
				WorkerID:   m.WorkerID,
				Name:       m.Name,
				Protocol:   m.Protocol,
				StatusCode: m.StatusCode,
				Status:     m.StatusString,
				DurationMs: m.Duration.Milliseconds(),
			}
			if m.Error != nil {
				rec.Error = m.Error.Error()
			}
			if err := enc.Encode(rec); err != nil {
				return errs.Wrap(err, errs.KindInternal, "failed to encode metric record")
			}
		}
	}
	return nil
}