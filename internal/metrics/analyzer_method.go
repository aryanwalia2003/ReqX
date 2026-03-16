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

	var globalDurations []time.Duration
	var totalSuccess, totalFailures int

	for _, iterMetrics := range allMetrics {
		for _, m := range iterMetrics {
			if m.Protocol != "HTTP" && m.Protocol != "" {
				continue // skip WS/SOCKET for latency stats
			}
			stat, exists := byName[m.Name]
			if !exists {
				stat = &RequestStat{Name: m.Name}
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
			if m.Duration > 0 {
				stat.Durations = append(stat.Durations, m.Duration)
				globalDurations = append(globalDurations, m.Duration)
			}
		}
	}

	// Compute per-request percentiles
	perRequest := make([]RequestStat, 0, len(order))
	for _, name := range order {
		s := byName[name]
		sorted := sortDurations(s.Durations)
		s.P50 = percentile(sorted, 0.50)
		s.P90 = percentile(sorted, 0.90)
		s.P95 = percentile(sorted, 0.95)
		s.P99 = percentile(sorted, 0.99)
		s.AvgDuration = avg(sorted)
		s.Durations = sorted
		perRequest = append(perRequest, *s)
	}

	// Global percentiles
	allSorted := sortDurations(globalDurations)
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
		AvgLatency:    avg(allSorted),
		P50:           percentile(allSorted, 0.50),
		P90:           percentile(allSorted, 0.90),
		P95:           percentile(allSorted, 0.95),
		P99:           percentile(allSorted, 0.99),
		RPS:           rps,
		TotalDuration: totalDuration,
		PerRequest:    perRequest,
	}
}

func avg(sorted []time.Duration) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	var sum time.Duration
	for _, d := range sorted {
		sum += d
	}
	return sum / time.Duration(len(sorted))
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