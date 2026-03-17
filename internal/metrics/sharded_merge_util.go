package metrics

import "time"

type mergedStats struct {
	order           []string
	byName          map[string]*RequestStat
	globalDurations []time.Duration
	totalSuccess    int
	totalFailures   int
}

func mergeShardResults(results []shardResult, order []string) mergedStats {
	byName := make(map[string]*RequestStat, len(order))
	var globalDurations []time.Duration
	var totalSuccess, totalFailures int

	for i := range results {
		r := results[i]
		totalSuccess += r.totalSuccess
		totalFailures += r.totalFailures
		globalDurations = append(globalDurations, r.globalDurations...)
		mergeShardMaps(byName, r.byName)
	}

	return mergedStats{
		order:           order,
		byName:          byName,
		globalDurations: globalDurations,
		totalSuccess:    totalSuccess,
		totalFailures:   totalFailures,
	}
}

func mergeShardMaps(dst map[string]*RequestStat, src map[string]*RequestStat) {
	for name, stat := range src {
		existing, ok := dst[name]
		if !ok {
			dst[name] = stat
			continue
		}
		existing.TotalRuns += stat.TotalRuns
		existing.Successes += stat.Successes
		existing.Failures += stat.Failures
		existing.Durations = append(existing.Durations, stat.Durations...)
		mergeErrorGroups(&existing.TopErrors, stat.TopErrors)
	}
}

func mergeErrorGroups(dst *[]ErrorGroup, src []ErrorGroup) {
	for _, g := range src {
		for i := range *dst {
			if (*dst)[i].Message == g.Message {
				(*dst)[i].Count += g.Count
				goto next
			}
		}
		*dst = append(*dst, g)
	next:
	}
}

func finalizeReport(m mergedStats, totalDuration time.Duration) Report {
	// Per-request percentiles in original "first-seen" order.
	perRequest := make([]RequestStat, 0, len(m.order))
	for _, name := range m.order {
		s := m.byName[name]
		if s == nil {
			continue
		}
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
	allSorted := sortDurations(m.globalDurations)
	totalReqs := m.totalSuccess + m.totalFailures

	var successRate float64
	if totalReqs > 0 {
		successRate = float64(m.totalSuccess) / float64(totalReqs) * 100
	}

	var rps float64
	if totalDuration > 0 {
		rps = float64(totalReqs) / totalDuration.Seconds()
	}

	return Report{
		TotalRequests: totalReqs,
		TotalSuccess:  m.totalSuccess,
		TotalFailures: m.totalFailures,
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

