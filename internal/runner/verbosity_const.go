package runner

// Verbosity controls how much per-request output the CollectionRunner emits.
const (
	VerbosityQuiet  = 0 // no per-request logs; progress bar handles output
	VerbosityNormal = 1 // status line per request (default)
	VerbosityFull   = 2 // full headers, body, and timing breakdown
)