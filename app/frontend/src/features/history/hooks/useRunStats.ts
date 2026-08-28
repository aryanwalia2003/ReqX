import { getRunStats } from '@/features/history/api'
import { useAsync } from '@/hooks'

/** `run(runId)` bhejta; `data` us run ka per-request breakdown hai. */
export function useRunStats() {
  return useAsync(getRunStats)
}
