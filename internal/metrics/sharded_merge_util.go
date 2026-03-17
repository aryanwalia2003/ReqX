package metrics

import (
	"time"

	"github.com/HdrHistogram/hdrhistogram-go"
)

type mergedStats struct {
	order           []string
	byName          map[string]*RequestStat
	globalHistogram *hdrhistogram.Histogram
	totalSuccess    int
	totalFailures   int
}

func mergeShardResults(results []shardResult, order []string) mergedStats {
	byName := make(map[string]*RequestStat, len(order))
	globalHistogram := newHistogram()
	var totalSuccess, totalFailures int

	for i := range results {
		r := results[i]
		totalSuccess += r.totalSuccess
		totalFailures += r.totalFailures
		if r.globalHistogram != nil {
			_ = globalHistogram.Merge(r.globalHistogram)
		}
		mergeShardMaps(byName, r.byName)
	}

	return mergedStats{
		order:           order,
		byName:          byName,
		globalHistogram: globalHistogram,
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
		if existing.Histogram == nil {
			existing.Histogram = newHistogram()
		}
		if stat.Histogram != nil {
			_ = existing.Histogram.Merge(stat.Histogram)
		}
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
		s.P50 = durFromQuantileMs(s.Histogram, 50)
		s.P90 = durFromQuantileMs(s.Histogram, 90)
		s.P95 = durFromQuantileMs(s.Histogram, 95)
		s.P99 = durFromQuantileMs(s.Histogram, 99)
		s.AvgDuration = durFromMeanMs(s.Histogram)
		perRequest = append(perRequest, *s)
	}

	// Global percentiles
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
		AvgLatency:    durFromMeanMs(m.globalHistogram),
		P50:           durFromQuantileMs(m.globalHistogram, 50),
		P90:           durFromQuantileMs(m.globalHistogram, 90),
		P95:           durFromQuantileMs(m.globalHistogram, 95),
		P99:           durFromQuantileMs(m.globalHistogram, 99),
		RPS:           rps,
		TotalDuration: totalDuration,
		PerRequest:    perRequest,
	}
}

