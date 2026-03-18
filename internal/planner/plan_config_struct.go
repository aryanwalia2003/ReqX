package planner

// PlanConfig carries the CLI flags that determine how a Collection is transformed
// into an ExecutionPlan. It is populated once in cmd/ and passed to BuildExecutionPlan.
type PlanConfig struct {
	// Filtering: only include requests whose names contain one of these substrings.
	RequestFilters []string

	// Injection: temporarily insert a new request at a 1-based index position.
	InjIndex  string // 1-based position string, e.g. "3"
	InjName   string
	InjMethod string
	InjURL    string
	InjBody   string
	InjHeaders []string
}