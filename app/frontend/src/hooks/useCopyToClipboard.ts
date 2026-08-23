import { useCallback, useState } from 'react'

export type CopyStatus = 'idle' | 'copied' | 'error'

/** navigator.clipboard wrap karta — status se copy feedback UI banao. */
export function useCopyToClipboard(): [CopyStatus, (text: string) => Promise<void>] {
  const [status, setStatus] = useState<CopyStatus>('idle')

  const copy = useCallback(async (text: string) => {
    try {
      await navigator.clipboard.writeText(text)
      setStatus('copied')
    } catch {
      setStatus('error')
    }
  }, [])

  return [status, copy]
}
