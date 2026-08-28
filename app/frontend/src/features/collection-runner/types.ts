// Types local to the "collection-runner" feature.

import type { AuthConfig } from '@/features/send-request'

/** Not every field internal/collection.Request (Go) carries — it also has
 * scripts, DAG deps, socket events; those pass through unedited when the
 * loaded Collection is round-tripped back into Run. headers/body/auth are
 * here (not just name/method/url) so a request can be opened into
 * SendRequestPanel for editing (see CollectionsSidebar). */
export interface CollectionRequest {
  name: string
  method: string
  url: string
  headers?: Record<string, string>
  body?: string
  auth?: AuthConfig
}

/** Mirrors internal/collection.Collection (Go), narrowed to what's shown. */
export interface Collection {
  name: string
  requests: CollectionRequest[]
}

/** Mirrors internal/environment.Environment (Go). */
export interface Environment {
  name: string
  variables: Record<string, string>
}

/** Mirrors internal/personas.Persona (Go) — one CSV row keyed by column name;
 * each key becomes {{persona.<key>}} in a request. */
export type Persona = Record<string, string>

/**
 * Mirrors app/services.RunCollectionInput (Go). workers/iterations/durationMs/rps
 * are the same load-testing knobs as the CLI's `reqx run -c/-n/-d/--rps` —
 * all optional, all default to a single-pass run when omitted. personas is
 * the CLI's `--personas` CSV, round-tripped from OpenPersonas.
 */
export interface RunCollectionInput {
  collection: Collection
  envVariables?: Record<string, string>
  workers?: number
  iterations?: number
  durationMs?: number
  rps?: number
  personas?: Persona[]
}

/** Mirrors app/services.RequestStat (Go) — aggregated across every iteration
 * a request ran in, not one row per iteration (unreadable at load-test scale). */
export interface RequestStat {
  name: string
  totalRuns: number
  successes: number
  failures: number
  avgLatencyMs: number
  p95LatencyMs: number
  topError?: string
}

/** Mirrors app/services.RunCollectionSummary (Go). */
export interface RunCollectionSummary {
  totalRequests: number
  totalSuccess: number
  totalFailures: number
  successRate: number
  avgLatencyMs: number
  p95LatencyMs: number
  totalDurationMs: number
}

/** Mirrors app/services.RunCollectionOutput (Go). */
export interface RunCollectionOutput {
  stats: RequestStat[]
  summary: RunCollectionSummary
}
