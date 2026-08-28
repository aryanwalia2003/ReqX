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
