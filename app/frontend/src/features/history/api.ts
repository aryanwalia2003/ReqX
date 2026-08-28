import { toResult } from '@/lib/result'
import { GetRunStats, ListRuns } from '@wails/go/services/HistoryService'

/** Recent runs, newest first — same query CLI ke `reqx ui` dashboard use karta. */
export function listRuns(limit = 50) {
  return toResult(ListRuns(limit))
}

/** Ek run ka per-request breakdown. */
export function getRunStats(runId: string) {
  return toResult(GetRunStats(runId))
}
