import { runCollection } from '@/features/collection-runner/api'
import { useAsync } from '@/hooks'

/** `run(input)` bhejta; `data` RunCollectionOutput (results + summary) hai. */
export function useRunCollection() {
  return useAsync(runCollection)
}
