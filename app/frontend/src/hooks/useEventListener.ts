import { useEffect, useRef } from 'react'
import type { RefObject } from 'react'

type Target = Window | HTMLElement

/** Latest handler ref me rakhta — listener baar baar re-attach nahi hota. */
export function useEventListener<K extends keyof WindowEventMap>(
  eventName: K,
  handler: (event: WindowEventMap[K]) => void,
  target?: RefObject<Target | null>,
): void {
  const savedHandler = useRef(handler)

  useEffect(() => {
    savedHandler.current = handler
  }, [handler])

  useEffect(() => {
    const node = target ? target.current : window
    if (!node) return

    const listener = (event: Event) => savedHandler.current(event as WindowEventMap[K])
    node.addEventListener(eventName, listener)
    return () => node.removeEventListener(eventName, listener)
  }, [eventName, target])
}
