import { openCollection } from '@/features/collection-runner/api'
import { useAsync } from '@/hooks'

/** `run(path)` bhejta; `data` loaded Collection hai. */
export function useOpenCollection() {
  return useAsync(openCollection)
}
