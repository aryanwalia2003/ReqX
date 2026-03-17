package metrics

import (
	"reqx/internal/runner"
	"time"
)

// Analyze aggregates all metrics from all iterations into a single Report.
// allMetrics[i] is the slice of RequestMetrics produced by iteration i.
func Analyze(allMetrics [][]runner.RequestMetric, totalDuration time.Duration) Report {
	// Index per-request stats by name, preserving first-seen order.
	order := []string{}
	byName := map[string]*RequestStat{}

	globalHistogram := newHistogram()
	var totalSuccess, totalFailures int

	for _, iterMetrics := range allMetrics {
		for _, m := range iterMetrics {
			if m.Protocol != "HTTP" && m.Protocol != "" {
				continue // skip WS/SOCKET for latency stats
			}
			stat, exists := byName[m.Name]
			if !exists {
				stat = &RequestStat{Name: m.Name, Histogram: newHistogram()}
				byName[m.Name] = stat
				order = append(order, m.Name)
			}
			stat.TotalRuns++
			failed := m.Error != nil || (m.StatusCode != 0 && m.StatusCode >= 400)
			if failed {
				stat.Failures++
				totalFailures++
				msg := errorMessage(m)
				addError(&stat.TopErrors, msg)
			} else {
				stat.Successes++
				totalSuccess++
			}
			recordDurationMs(stat.Histogram, m.Duration)
			recordDurationMs(globalHistogram, m.Duration)
		}
	}

	// Compute per-request percentiles
	perRequest := make([]RequestStat, 0, len(order))
	for _, name := range order {
		s := byName[name]
		s.P50 = durFromQuantileMs(s.Histogram, 50)
		s.P90 = durFromQuantileMs(s.Histogram, 90)
		s.P95 = durFromQuantileMs(s.Histogram, 95)
		s.P99 = durFromQuantileMs(s.Histogram, 99)
		s.AvgDuration = durFromMeanMs(s.Histogram)
		perRequest = append(perRequest, *s)
	}

	// Global percentiles
	totalReqs := totalSuccess + totalFailures
	var successRate float64
	if totalReqs > 0 {
		successRate = float64(totalSuccess) / float64(totalReqs) * 100
	}
	var rps float64
	if totalDuration > 0 {
		rps = float64(totalReqs) / totalDuration.Seconds()
	}

	return Report{
		TotalRequests: totalReqs,
		TotalSuccess:  totalSuccess,
		TotalFailures: totalFailures,
		SuccessRate:   successRate,
		AvgLatency:    durFromMeanMs(globalHistogram),
		P50:           durFromQuantileMs(globalHistogram, 50),
		P90:           durFromQuantileMs(globalHistogram, 90),
		P95:           durFromQuantileMs(globalHistogram, 95),
		P99:           durFromQuantileMs(globalHistogram, 99),
		RPS:           rps,
		TotalDuration: totalDuration,
		PerRequest:    perRequest,
	}
}

func errorMessage(m runner.RequestMetric) string {
	if m.Error != nil {
		return m.Error.Error()
	}
	return m.StatusString
}

func addError(groups *[]ErrorGroup, msg string) {
	for i := range *groups {
		if (*groups)[i].Message == msg {
			(*groups)[i].Count++
			return
		}
	}
	*groups = append(*groups, ErrorGroup{Message: msg, Count: 1})
}