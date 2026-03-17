package metrics

import (
	"time"

	"github.com/HdrHistogram/hdrhistogram-go"
)

const (
	hdrMinMs     int64 = 1
	hdrMaxMs     int64 = 60 * 60 * 1000 // 1 hour
	hdrSigFigs          = 3
)

func newHistogram() *hdrhistogram.Histogram {
	return hdrhistogram.New(hdrMinMs, hdrMaxMs, hdrSigFigs)
}

func recordDurationMs(h *hdrhistogram.Histogram, d time.Duration) {
	if h == nil || d <= 0 {
		return
	}
	ms := d.Milliseconds()
	if ms < hdrMinMs {
		ms = hdrMinMs
	} else if ms > hdrMaxMs {
		ms = hdrMaxMs
	}
	_ = h.RecordValue(ms)
}

func durFromQuantileMs(h *hdrhistogram.Histogram, q float64) time.Duration {
	if h == nil || h.TotalCount() == 0 {
		return 0
	}
	return time.Duration(h.ValueAtQuantile(q)) * time.Millisecond
}

func durFromMeanMs(h *hdrhistogram.Histogram) time.Duration {
	if h == nil || h.TotalCount() == 0 {
		return 0
	}
	return time.Duration(h.Mean()) * time.Millisecond
}

