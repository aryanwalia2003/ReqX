package metrics

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// PrintReport renders a full load test report to stdout.
// This is the single source of truth for all terminal output after a run.
func PrintReport(r Report) {
	sep := strings.Repeat("═", 70)
	thin := strings.Repeat("─", 70)

	header := color.New(color.FgHiCyan, color.Bold)
	header.Printf("\n%s\n", sep)
	header.Println("  LOAD TEST REPORT")
	header.Printf("%s\n", sep)

	// ── Global stats ──────────────────────────────────────────
	fmt.Printf("  Total Requests : %d\n", r.TotalRequests)
	successPct := color.GreenString("%.2f%%", r.SuccessRate)
	if r.SuccessRate < 95 {
		successPct = color.RedString("%.2f%%", r.SuccessRate)
	}
	fmt.Printf("  Success Rate   : %s  (%s passed / %s failed)\n",
		successPct,
		color.GreenString("%d", r.TotalSuccess),
		color.RedString("%d", r.TotalFailures),
	)
	fmt.Printf("  Throughput     : %s req/s\n", color.CyanString("%.2f", r.RPS))
	fmt.Printf("  Total Run Time : %v\n", r.TotalDuration.Round(1000000))

	// ── Latency percentiles ───────────────────────────────────
	header.Printf("\n  %s\n", "LATENCY  (all HTTP requests)")
	header.Printf("  %s\n", thin[:40])
	fmt.Printf("  %-8s %v\n", "Avg", fmtDur(r.AvgLatency))
	fmt.Printf("  %-8s %v\n", "P50", fmtDur(r.P50))
	fmt.Printf("  %-8s %v\n", "P90", fmtDur(r.P90))
	fmt.Printf("  %-8s %v\n", color.YellowString("P95"), color.YellowString(fmtDur(r.P95)))
	fmt.Printf("  %-8s %v\n", color.RedString("P99"), color.RedString(fmtDur(r.P99)))

	// ── Per-request breakdown ─────────────────────────────────
	header.Printf("\n  %s\n", "PER-REQUEST BREAKDOWN")
	fmt.Printf("  %-32s %7s %7s %7s  %8s  %8s  %s\n",
		"Request", "Runs", "Pass", "Fail", "Avg", "P95", "Top Error")
	header.Printf("  %s\n", thin)

	for _, s := range r.PerRequest {
		failCol := fmt.Sprintf("%7d", s.Failures)
		if s.Failures > 0 {
			failCol = color.RedString("%7d", s.Failures)
		}
		topErr := "-"
		if len(s.TopErrors) > 0 {
			e := s.TopErrors[0]
			msg := e.Message
			if len(msg) > 28 {
				msg = msg[:25] + "..."
			}
			topErr = fmt.Sprintf("%s ×%d", msg, e.Count)
		}
		fmt.Printf("  %-32s %7d %7d %s  %8s  %8s  %s\n",
			truncate(s.Name, 32),
			s.TotalRuns,
			s.Successes,
			failCol,
			fmtDur(s.AvgDuration),
			fmtDur(s.P95),
			topErr,
		)
	}
	header.Printf("  %s\n\n", sep)
}

func fmtDur(d interface{ Milliseconds() int64 }) string {
	ms := d.Milliseconds()
	if ms >= 1000 {
		return fmt.Sprintf("%.2fs", float64(ms)/1000)
	}
	return fmt.Sprintf("%dms", ms)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}