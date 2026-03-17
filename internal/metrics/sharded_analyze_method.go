package metrics

import (
	"runtime"
	"sync"
	"time"

	"reqx/internal/runner"
)

// AnalyzeSharded aggregates metrics using N sharded goroutines.
// If shards <= 0, an automatic shard count is chosen.
//   - Dispatcher: single goroutine loop in this function (hash(name)%N)
//   - Shard channels: chans[]
//   - Per-shard aggregators: consumeShard() goroutines
//   - Final merge: mergeShardResults() + mergeShardMaps()
func AnalyzeSharded(allMetrics [][]runner.RequestMetric, totalDuration time.Duration, shards int) Report {
	shards = normalizeShardCount(shards)
	if shards == 1 {
		return Analyze(allMetrics, totalDuration)
	}

	chans := make([]chan runner.RequestMetric, shards)
	results := make([]shardResult, shards)

	var wg sync.WaitGroup
	wg.Add(shards)
	for i := 0; i < shards; i++ {
		chans[i] = make(chan runner.RequestMetric, 4096)
		go func(idx int) {
			defer wg.Done()
			results[idx] = consumeShard(chans[idx])
		}(i)
	}

	// Dispatcher preserves "first seen" order deterministically.
	order := make([]string, 0, 64)
	seen := make(map[string]bool, 64)

	for _, iterMetrics := range allMetrics {
		for _, m := range iterMetrics {
			if m.Protocol != "HTTP" && m.Protocol != "" {
				continue
			}
			if !seen[m.Name] {
				seen[m.Name] = true
				order = append(order, m.Name)
			}
			chans[shardFor(m.Name, shards)] <- m
		}
	}

	for i := 0; i < shards; i++ {
		close(chans[i])
	}
	wg.Wait()

	merged := mergeShardResults(results, order)
	return finalizeReport(merged, totalDuration)
}

func normalizeShardCount(shards int) int {
	if shards <= 0 {
		shards = runtime.GOMAXPROCS(0)
		if shards < 2 {
			shards = 2
		}
	}
	if shards > 64 {
		shards = 64
	}
	return shards
}

