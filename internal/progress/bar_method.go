package progress

import (
	"fmt"
	"strings"
	"time"
)

const barWidth = 30

// render writes a single overwriting line to stdout.
func (b *Bar) render() {
	done := b.done.Load()
	errs := b.errors.Load()
	total := b.total

	pct := float64(done) / float64(total)
	if pct > 1 {
		pct = 1
	}

	filled := int(pct * barWidth)
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	elapsed := time.Since(b.startTime)
	var rps float64
	if elapsed.Seconds() > 0 {
		rps = float64(done) / elapsed.Seconds()
	}

	fmt.Printf("\r  [%s] %3.0f%%  Workers: %d  Done: %d/%d  Errors: %d  RPS: %.1f  Elapsed: %s   ",
		bar,
		pct*100,
		b.workers,
		done, total,
		errs,
		rps,
		fmtElapsed(elapsed),
	)
}

func fmtElapsed(d time.Duration) string {
	d = d.Round(time.Second)
	m := d / time.Minute
	s := (d % time.Minute) / time.Second
	return fmt.Sprintf("%02d:%02d", m, s)
}