import { renderHook } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'

import { useEventListener } from '@/hooks/useEventListener'

describe('useEventListener', () => {
  it('calls the latest handler and cleans up on unmount', () => {
    const handler = vi.fn()
    const { rerender, unmount } = renderHook(({ fn }) => useEventListener('click', fn), {
      initialProps: { fn: handler },
    })

    window.dispatchEvent(new MouseEvent('click'))
    expect(handler).toHaveBeenCalledTimes(1)

    const nextHandler = vi.fn()
    rerender({ fn: nextHandler })
    window.dispatchEvent(new MouseEvent('click'))
    expect(handler).toHaveBeenCalledTimes(1)
    expect(nextHandler).toHaveBeenCalledTimes(1)

    unmount()
    window.dispatchEvent(new MouseEvent('click'))
    expect(nextHandler).toHaveBeenCalledTimes(1)
  })

  it('binds to a ref target instead of window when given one', () => {
    const el = document.createElement('div')
    document.body.appendChild(el)
    const ref = { current: el }
    const handler = vi.fn()

    renderHook(() => useEventListener('click', handler, ref))
    el.dispatchEvent(new MouseEvent('click'))
    expect(handler).toHaveBeenCalledTimes(1)

    document.body.removeChild(el)
  })
})
