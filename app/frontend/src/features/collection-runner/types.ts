// Types local to the "collection-runner" feature.

/** Only the fields this feature displays — internal/collection.Request (Go)
 * carries more (scripts, DAG deps, socket events); those pass through
 * unedited when the loaded Collection is round-tripped back into Run. */
export interface CollectionRequest {
  name: string
  method: string
  url: string
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

/** Mirrors app/services.RunCollectionInput (Go). */
export interface RunCollectionInput {
  collection: Collection
  envVariables?: Record<string, string>
}

/** Mirrors app/services.RequestResult (Go). */
export interface RequestResult {
  name: string
  protocol: string
  statusCode: number
  statusString: string
  durationMs: number
  bytesReceived: number
  errorMessage?: string
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
  results: RequestResult[]
  summary: RunCollectionSummary
}
