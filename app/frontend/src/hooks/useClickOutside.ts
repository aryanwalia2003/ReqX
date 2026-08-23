import { useEffect, useRef } from 'react'
import type { RefObject } from 'react'

/** Ref ke bahar click hone par handler chalao — Dialog/Select dismiss ke liye. */
export function useClickOutside(ref: RefObject<HTMLElement | null>, handler: () => void): void {
  const savedHandler = useRef(handler)

  useEffect(() => {
    savedHandler.current = handler
  }, [handler])

  useEffect(() => {
    const listener = (event: MouseEvent) => {
      const node = ref.current
      if (!node || node.contains(event.target as Node)) return
      savedHandler.current()
    }
    document.addEventListener('mousedown', listener)
    return () => document.removeEventListener('mousedown', listener)
  }, [ref])
}
