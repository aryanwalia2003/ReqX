package runner

import (
	"fmt"
	"sync"

	"github.com/dop251/goja"
	"github.com/fatih/color"

	"reqx/internal/dag"
	"reqx/internal/environment"
	"reqx/internal/planner"
)

// RunDAG executes the requests in plan according to the dependency graph.
//
// Parallelism model:
//   - Nodes in the same topological level run concurrently.
//   - A sync.WaitGroup barrier separates each level.
//   - Every goroutine receives a snapshot clone of the shared environment so
//     concurrent pm.env.set calls never race on the same map.
//   - After each level's barrier, writes from all node clones are merged back
//     into ctx.Environment (last-write-wins on key conflicts).
//
// Skip propagation:
//   - A node is skipped when any direct dependency was itself skipped.
//   - A node is skipped when its condition evaluates to false against the
//     representative EvalContext for that level (see worstEvalContext).
func (cr *CollectionRunner) RunDAG(plan *planner.ExecutionPlan, ctx *RuntimeContext) ([]RequestMetric, error) {
	levels, err := dag.TopoSort(plan.DAG)
	if err != nil {
		return nil, err
	}

	n := len(plan.Requests)
	allMetrics := make([]RequestMetric, n)
	evalCtxs := make([]dag.EvalContext, n)
	skipped := make([]bool, n)

	// Top-level RunDAG owns the async connection lifecycle ONLY when
	// PersistConnections is false (single-run or WorkerPool mode).
	// In scheduler VU mode, the worker goroutine in scheduler_worker_method.go
	// owns the stop/wait so sockets survive across iterations.
	defer func() {
		if ctx.PersistConnections {
			return // scheduler worker will close & wait
		}
		if cr.verbosity >= VerbosityNormal {
			color.Cyan("\n[DAG] All levels done. Waiting for background connections...\n")
		}
		ctx.AsyncStopOnce.Do(func() { close(ctx.AsyncStop) })
		ctx.AsyncWG.Wait()
		if cr.verbosity >= VerbosityNormal {
			color.Green("[DAG] All background connections closed cleanly.\n")
		}
	}()

	for levelIdx, level := range levels {
		if cr.verbosity >= VerbosityNormal {
			color.Cyan("\n[DAG] Level %d — %d node(s) in parallel\n", levelIdx, len(level))
		}

		type nodeResult struct {
			reqIdx  int
			metric  RequestMetric
			eval    dag.EvalContext
			skip    bool
			envDiff map[string]string // variables written by pm.env.set during this node
		}

		results := make([]nodeResult, len(level))
		var wg sync.WaitGroup

		for slot, nodeIdx := range level {
			wg.Add(1)
			go func(slot, reqIdx int) {
				defer wg.Done()

				req := plan.Requests[reqIdx]

				// ── 1. Skip propagation ──────────────────────────────────────
				for _, depIdx := range plan.DAG.Edges[reqIdx] {
					if skipped[depIdx] {
						if cr.verbosity >= VerbosityNormal {
							color.Yellow("[DAG] Skip %q — dep %q was skipped\n",
								req.Name, plan.Requests[depIdx].Name)
						}
						results[slot] = nodeResult{
							reqIdx: reqIdx,
							metric: RequestMetric{Name: req.Name, Protocol: "HTTP", StatusString: "SKIPPED"},
							skip:   true,
						}
						return
					}
				}

				// ── 2. Condition check ───────────────────────────────────────
				// Evaluate the condition against the worst result among all
				// declared dependencies. "Worst" means: prefer Failed=true,
				// then highest status code. This makes "failed == false" behave
				// correctly for fan-in nodes (all deps must have succeeded).
				if req.Condition != "" && len(plan.DAG.Edges[reqIdx]) > 0 {
					repr := worstEvalContext(evalCtxs, plan.DAG.Edges[reqIdx])
					ok, condErr := dag.EvalCondition(req.Condition, repr)
					if condErr != nil {
						color.Yellow("[DAG] Condition error on %q: %v — skipping\n", req.Name, condErr)
						results[slot] = nodeResult{
							reqIdx: reqIdx,
							metric: RequestMetric{Name: req.Name, Protocol: "HTTP", StatusString: "SKIPPED", ErrorMsg: condErr.Error()},
							skip:   true,
						}
						return
					}
					if !ok {
						if cr.verbosity >= VerbosityNormal {
							color.Yellow("[DAG] Skip %q — condition %q not met\n", req.Name, req.Condition)
						}
						results[slot] = nodeResult{
							reqIdx: reqIdx,
							metric: RequestMetric{Name: req.Name, Protocol: "HTTP", StatusString: "SKIPPED"},
							skip:   true,
						}
						return
					}
				}

				// ── 3. Isolated environment clone ────────────────────────────
				// Each node goroutine gets its own copy of the current env so
				// concurrent pm.env.set calls are safe. We snapshot the parent
				// env once (before the goroutine was scheduled, under the
				// implicit happens-before of WaitGroup), then give the node a
				// RuntimeContext pointing at this private copy.
				nodeCtx := ctx.CloneForNode()

				// ── 4. Execute ───────────────────────────────────────────────
				subPlan := &planner.ExecutionPlan{
					Requests:        plan.Requests[reqIdx : reqIdx+1],
					CollectionAuth:  plan.CollectionAuth,
					CompiledScripts: remapCompiledScripts(plan, reqIdx),
					DAG:             nil, // sub-plan always runs linearly
				}

				nodeMetrics, runErr := cr.Run(subPlan, nodeCtx)
				if runErr != nil {
					fmt.Printf("[DAG] Node %q error: %v\n", req.Name, runErr)
				}

				var m RequestMetric
				if len(nodeMetrics) > 0 {
					m = nodeMetrics[0]
				} else {
					m = RequestMetric{Name: req.Name, Protocol: "HTTP", Error: runErr}
				}

				// Compute the diff: keys that the node wrote to its clone.
				envDiff := diffEnv(ctx.Environment.Variables, nodeCtx.Environment.Variables)

				results[slot] = nodeResult{
					reqIdx:  reqIdx,
					metric:  m,
					envDiff: envDiff,
					eval: dag.EvalContext{
						StatusCode: m.StatusCode,
						DurationMs: m.Duration.Milliseconds(),
						Failed:     m.Error != nil || (m.StatusCode >= 400 && m.StatusCode != 0),
					},
				}
			}(slot, nodeIdx)
		}

		wg.Wait()

		// ── 5. Merge env diffs back into shared context ──────────────────────
		// This runs sequentially after the barrier so there is no concurrent
		// access to ctx.Environment at this point. Last-write-wins on conflicts.
		for _, r := range results {
			for k, v := range r.envDiff {
				ctx.Environment.Variables[k] = v
			}
			allMetrics[r.reqIdx] = r.metric
			evalCtxs[r.reqIdx] = r.eval
			skipped[r.reqIdx] = r.skip
		}
	}

	return allMetrics, nil
}

// worstEvalContext picks the most "failed" EvalContext from a set of dependency
// indices. Rules (in priority order):
//  1. If any dep has Failed=true, return that one.
//  2. Otherwise return the dep with the highest StatusCode.
//
// This gives conditions like "failed == false" the correct semantics for fan-in
// nodes: all deps must have succeeded.
func worstEvalContext(evalCtxs []dag.EvalContext, depIndices []int) dag.EvalContext {
	worst := evalCtxs[depIndices[0]]
	for _, idx := range depIndices[1:] {
		c := evalCtxs[idx]
		if c.Failed && !worst.Failed {
			worst = c
			continue
		}
		if !worst.Failed && c.StatusCode > worst.StatusCode {
			worst = c
		}
	}
	return worst
}

// diffEnv returns the keys whose values differ between parent and child env maps.
// These are the variables that the node goroutine wrote during its execution.
func diffEnv(parent, child map[string]string) map[string]string {
	diff := make(map[string]string)
	for k, v := range child {
		if parent[k] != v {
			diff[k] = v
		}
	}
	return diff
}

// remapCompiledScripts returns a CompiledScripts map with reqIdx remapped to 0.
// Run() looks up compiled scripts by position in the sub-plan's Requests slice,
// which is always 0 for single-node sub-plans.
func remapCompiledScripts(plan *planner.ExecutionPlan, reqIdx int) map[planner.ScriptKey]*goja.Program {
	if plan.CompiledScripts == nil {
		return nil
	}
	out := make(map[planner.ScriptKey]*goja.Program)
	for k, v := range plan.CompiledScripts {
		if k.RequestIndex == reqIdx {
			out[planner.ScriptKey{RequestIndex: 0, ScriptType: k.ScriptType}] = v
		}
	}
	return out
}

// cloneEnvVars returns a shallow copy of vars. Callers get their own map so
// concurrent writes never touch the parent's map.
func cloneEnvVars(vars map[string]string) map[string]string {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		out[k] = v
	}
	return out
}

// newEnvSnapshot returns an Environment whose Variables map is a copy of src.
func newEnvSnapshot(src *environment.Environment) *environment.Environment {
	if src == nil {
		return environment.NewEnvironment("snapshot")
	}
	return &environment.Environment{
		Name:      src.Name,
		Variables: cloneEnvVars(src.Variables),
	}
}