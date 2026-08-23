import { renderHook } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'

import { useKeyboardShortcut } from '@/hooks/useKeyboardShortcut'

function fireKey(init: KeyboardEventInit) {
  window.dispatchEvent(new KeyboardEvent('keydown', init))
}

describe('useKeyboardShortcut', () => {
  it('fires only when the exact combo matches', () => {
    const handler = vi.fn()
    renderHook(() => useKeyboardShortcut('mod+enter', handler))

    fireKey({ key: 'Enter' })
    expect(handler).not.toHaveBeenCalled()

    fireKey({ key: 'Enter', metaKey: true })
    expect(handler).toHaveBeenCalledOnce()
  })

  it('picks up the latest handler without re-registering the listener', () => {
    const first = vi.fn()
    const second = vi.fn()
    const { rerender } = renderHook(({ fn }) => useKeyboardShortcut('shift+a', fn), {
      initialProps: { fn: first },
    })

    rerender({ fn: second })
    fireKey({ key: 'a', shiftKey: true })
    expect(first).not.toHaveBeenCalled()
    expect(second).toHaveBeenCalledOnce()
  })
})
