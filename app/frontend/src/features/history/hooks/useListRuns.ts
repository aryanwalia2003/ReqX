import { listRuns } from '@/features/history/api'
import { useAsync } from '@/hooks'

/** `run(limit?)` bhejta; `data` recent runs (newest first) hai. */
export function useListRuns() {
  return useAsync(listRuns)
}
