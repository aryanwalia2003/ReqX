package dag

import (
	"fmt"
	"strconv"
	"strings"
)

// EvalContext carries the observable result of a completed node so that
// condition expressions can be evaluated against it.
type EvalContext struct {
	StatusCode int
	DurationMs int64
	Failed     bool // true if HTTP >= 400 or execution error occurred
}

// EvalCondition evaluates a condition string against ctx and returns true when
// the dependent node should be allowed to run.
//
// Supported expression forms:
//
//	"status == 200"
//	"status != 200"
//	"status >= 200"
//	"status <= 299"
//	"status > 199"
//	"status < 400"
//	"duration_ms < 500"
//	"failed == false"
//	"failed == true"
//
// An empty condition string always returns true (no guard).
func EvalCondition(condition string, ctx EvalContext) (bool, error) {
	condition = strings.TrimSpace(condition)
	if condition == "" {
		return true, nil
	}

	parts := strings.Fields(condition)
	if len(parts) != 3 {
		return false, fmt.Errorf("dag: condition %q must be of the form '<field> <op> <value>'", condition)
	}

	field, op, rawVal := parts[0], parts[1], parts[2]

	var lhs int64
	switch field {
	case "status":
		lhs = int64(ctx.StatusCode)
	case "duration_ms":
		lhs = ctx.DurationMs
	case "failed":
		boolVal := rawVal == "true"
		switch op {
		case "==":
			return ctx.Failed == boolVal, nil
		case "!=":
			return ctx.Failed != boolVal, nil
		default:
			return false, fmt.Errorf("dag: operator %q is not supported for 'failed' field", op)
		}
	default:
		return false, fmt.Errorf("dag: unknown condition field %q — supported: status, duration_ms, failed", field)
	}

	rhs, err := strconv.ParseInt(rawVal, 10, 64)
	if err != nil {
		return false, fmt.Errorf("dag: condition value %q is not a valid integer", rawVal)
	}

	switch op {
	case "==":
		return lhs == rhs, nil
	case "!=":
		return lhs != rhs, nil
	case ">":
		return lhs > rhs, nil
	case ">=":
		return lhs >= rhs, nil
	case "<":
		return lhs < rhs, nil
	case "<=":
		return lhs <= rhs, nil
	default:
		return false, fmt.Errorf("dag: unsupported operator %q — supported: ==, !=, >, >=, <, <=", op)
	}
}