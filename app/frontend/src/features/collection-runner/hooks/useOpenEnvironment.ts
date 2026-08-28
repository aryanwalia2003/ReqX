import { openEnvironment } from '@/features/collection-runner/api'
import { useAsync } from '@/hooks'

/** `run(path)` bhejta; `data` loaded Environment (name + variables) hai. */
export function useOpenEnvironment() {
  return useAsync(openEnvironment)
}
