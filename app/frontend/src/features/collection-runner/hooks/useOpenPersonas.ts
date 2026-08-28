import { openPersonas } from '@/features/collection-runner/api'
import { useAsync } from '@/hooks'

/** `run(path)` bhejta; `data` loaded Personas (CSV rows) hai. */
export function useOpenPersonas() {
  return useAsync(openPersonas)
}
