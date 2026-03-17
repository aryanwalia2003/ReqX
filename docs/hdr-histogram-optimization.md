# 📈 Optimization: High Dynamic Range (HDR) Histograms

This document outlines the architectural upgrade from raw latency storage to **HDR Histograms**. This change is critical for scaling ReqX to millions of requests without consuming gigabytes of RAM or slowing down the final report generation.

---

## 🛑 The Problem: Raw Latency Storage

Currently, ReqX stores every single request duration in a `[]time.Duration` slice.

### 1. Memory Consumption
If you run a test with **10 million requests**:
- Each `time.Duration` (int64) takes **8 bytes**.
- **10,000,000 * 8 bytes ≈ 80 MB** per request type just for the numbers.
- For a collection with many requests, this quickly reaches hundreds of megabytes or gigabytes.

### 2. CPU Bottleneck (Sorting)
To find the P95 or P99, we currently use `sort.Slice`.
- Sorting 10 million numbers is computationally expensive ($O(N \log N)$).
- The final **Analyze** step becomes a bottleneck, making the CLI hang for several seconds after the test finishes.

---

## 🚀 The Solution: HDR Histograms

An **HDR (High Dynamic Range) Histogram** is a specialized data structure designed to record a large number of values with high precision using a **fixed, small amount of memory**.

### How it works:
- It doesn't store every individual value.
- Instead, it groups values into "buckets" in a logarithmic way.
- It can record values from **1 nanosecond to 1 hour** with configurable precision (e.g., 3 significant digits).
- **Querying percentiles (P50, P99, P99.9) is $O(1)$**—it's nearly instantaneous regardless of whether you've recorded 1,000 or 1,000,000,000 values.

**Analogy:**
- **Raw Latencies:** Writing down every single person's exact height on a piece of paper.
- **HDR Histogram:** Creating a chart with buckets like "150-151cm," "151-152cm," etc., and just incrementing a counter. You lose the individual data point but keep the distribution perfectly.

---

## 🗺️ Implementation Roadmap

### Phase 1: Dependency Management
Add the battle-tested Go library to the project:
```bash
go get github.com/HdrHistogram/hdrhistogram-go
```

### Phase 2: Refactor Metrics Storage
Update `RequestStat` in `internal/metrics/report_struct.go`:
```go
type RequestStat struct {
    Name        string
    TotalRuns   int
    Successes   int
    Failures    int
    // Durations []time.Duration // REMOVED
    Histogram   *hdrhistogram.Histogram // ADDED
    TopErrors   []ErrorGroup
}
```

### Phase 3: Update Sharded Collection
When a shard receives a metric, instead of appending to a slice, it records it in the histogram:
```go
// Inside consumeShard loop
stat.Histogram.RecordValue(m.Duration.Milliseconds())
```

### Phase 4: Instant Analysis
Calculating percentiles becomes a simple query:
```go
s.P95 = time.Duration(s.Histogram.ValueAtQuantile(95)) * time.Millisecond
s.AvgDuration = time.Duration(s.Histogram.Mean()) * time.Millisecond
```
**No more sorting required.**

---

## 🔒 Thread Safety
Our current **Sharded Design** handles thread safety elegantly. Since each shard has its own `consumeShard` goroutine and private maps, we don't need a concurrent histogram. We only merge the histograms at the very end in a single-threaded merge step.

---
*Benefit: ReqX can now handle professional-grade load tests (millions of requests) on standard hardware with minimal overhead.*
