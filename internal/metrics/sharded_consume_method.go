package metrics

import (
	"reqx/internal/runner"

	"github.com/HdrHistogram/hdrhistogram-go"
)

type shardResult struct {
	byName          map[string]*RequestStat
	globalHistogram *hdrhistogram.Histogram
	totalSuccess    int
	totalFailures   int
}

func consumeShard(ch <-chan runner.RequestMetric) shardResult {
	byName := make(map[string]*RequestStat, 64)
	globalHistogram := newHistogram()
	var totalSuccess, totalFailures int

	for m := range ch {
		stat, ok := byName[m.Name]
		if !ok {
			stat = &RequestStat{Name: m.Name, Histogram: newHistogram()}
			byName[m.Name] = stat
		}

		stat.TotalRuns++
		failed := m.Error != nil || (m.StatusCode != 0 && m.StatusCode >= 400)
		if failed {
			stat.Failures++
			totalFailures++
			addError(&stat.TopErrors, errorMessage(m))
		} else {
			stat.Successes++
			totalSuccess++
		}

		if m.Duration > 0 {
			recordDurationMs(stat.Histogram, m.Duration)
			recordDurationMs(globalHistogram, m.Duration)
		}
	}

	return shardResult{
		byName:          byName,
		globalHistogram: globalHistogram,
		totalSuccess:    totalSuccess,
		totalFailures:   totalFailures,
	}
}

