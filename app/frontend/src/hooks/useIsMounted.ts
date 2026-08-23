import { useCallback, useEffect, useRef } from 'react'

/** Async callback me setState se pehle mount-check ke liye. */
export function useIsMounted(): () => boolean {
  const mounted = useRef(false)

  useEffect(() => {
    mounted.current = true
    return () => {
      mounted.current = false
    }
  }, [])

  return useCallback(() => mounted.current, [])
}
