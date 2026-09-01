// Types local to the "history" feature. JSON keys are snake_case to match
// internal/history.RunRow/StatRow (Go) exactly — those structs are the API
// response shape as-is, no camelCase DTO layer in between.

/** Mirrors internal/history.RunRow (Go). */
export interface RunRow {
  id: string
  ts: string
  collection: string
  total_reqs: number
  rps: number
  p95_ms: number
  error_pct: number
}

/** Mirrors internal/history.StatRow (Go). */
export interface StatRow {
  name: string
  successes: number
  failures: number
  p95_ms: number
  avg_ms: number
}

/** Mirrors internal/history.DagNodeRow (Go) — one node in a collection's
 * dependency graph (only present when it has depends_on edges). Used both
 * for a past run (HistoryService.GetRunStats' sibling GetDAGNodes) and a
 * just-finished one (CollectionService.Run's RunCollectionOutput.dagNodes). */
export interface DagNodeRow {
  name: string
  status: string
  duration_ms: number
  level_idx: number
  depends_on: string[]
}
