package metrics

import (
	"time"

	"reqx/internal/runner"
)

type shardResult struct {
	byName          map[string]*RequestStat
	globalDurations []time.Duration
	totalSuccess    int
	totalFailures   int
}

func consumeShard(ch <-chan runner.RequestMetric) shardResult {
	byName := make(map[string]*RequestStat, 64)
	var globalDurations []time.Duration
	var totalSuccess, totalFailures int

	for m := range ch {
		stat, ok := byName[m.Name]
		if !ok {
			stat = &RequestStat{Name: m.Name}
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
			stat.Durations = append(stat.Durations, m.Duration)
			globalDurations = append(globalDurations, m.Duration)
		}
	}

	return shardResult{
		byName:          byName,
		globalDurations: globalDurations,
		totalSuccess:    totalSuccess,
		totalFailures:   totalFailures,
	}
}

