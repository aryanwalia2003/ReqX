import { useEffect, useRef } from 'react'

type Modifier = 'mod' | 'shift' | 'alt'
const MODIFIERS: readonly Modifier[] = ['mod', 'shift', 'alt']

function isModifier(part: string): part is Modifier {
  return (MODIFIERS as readonly string[]).includes(part)
}

// 'mod' = Cmd on Mac, Ctrl elsewhere — jaisa <Kbd> me dikhaya jata hai.
function matchesCombo(event: KeyboardEvent, combo: string): boolean {
  const parts = combo.toLowerCase().split('+')
  const key = parts.pop()
  if (!key) return false

  const modifiers = parts.filter(isModifier)
  const isMod = event.metaKey || event.ctrlKey

  if (modifiers.includes('mod') !== isMod) return false
  if (modifiers.includes('shift') !== event.shiftKey) return false
  if (modifiers.includes('alt') !== event.altKey) return false
  return event.key.toLowerCase() === key
}

/**
 * Global keyboard shortcut — combo jaise 'mod+enter' ya 'shift+a'.
 * @example useKeyboardShortcut('mod+enter', sendRequest)
 */
export function useKeyboardShortcut(combo: string, handler: (event: KeyboardEvent) => void): void {
  const savedHandler = useRef(handler)

  useEffect(() => {
    savedHandler.current = handler
  }, [handler])

  useEffect(() => {
    const listener = (event: KeyboardEvent) => {
      if (!matchesCombo(event, combo)) return
      event.preventDefault()
      savedHandler.current(event)
    }
    window.addEventListener('keydown', listener)
    return () => window.removeEventListener('keydown', listener)
  }, [combo])
}
