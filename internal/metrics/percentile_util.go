package metrics

import (
	"sort"
	"time"
)

// percentile returns the p-th percentile (0.0–1.0) of a sorted duration slice.
// The slice must be sorted ascending before calling. Returns 0 for empty input.
func percentile(sorted []time.Duration, p float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(float64(len(sorted)-1) * p)
	return sorted[idx]
}

// sortDurations returns a new sorted copy — never mutates the original.
func sortDurations(in []time.Duration) []time.Duration {
	cp := make([]time.Duration, len(in))
	copy(cp, in)
	sort.Slice(cp, func(i, j int) bool { return cp[i] < cp[j] })
	return cp
}