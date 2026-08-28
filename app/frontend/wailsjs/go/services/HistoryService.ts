// Hand-written stand-in for what `wails generate module` would generate —
// see wailsjs/go/services/RequestService.ts for why this file exists and
// when to delete it.
import { callWailsMethod } from '@wails/go/wailsRuntime'

import type { RunRow, StatRow } from '@/features/history/types'

export function ListRuns(limit: number): Promise<RunRow[]> {
  return callWailsMethod(window.go?.services?.HistoryService?.ListRuns(limit))
}

export function GetRunStats(runId: string): Promise<StatRow[]> {
  return callWailsMethod(window.go?.services?.HistoryService?.GetRunStats(runId))
}
