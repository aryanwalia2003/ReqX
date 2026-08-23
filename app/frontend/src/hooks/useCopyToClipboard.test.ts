import { act, renderHook } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { useCopyToClipboard } from '@/hooks/useCopyToClipboard'

describe('useCopyToClipboard', () => {
  beforeEach(() => {
    Object.assign(navigator, { clipboard: { writeText: vi.fn() } })
  })

  it('reports copied status on success', async () => {
    const { result } = renderHook(() => useCopyToClipboard())
    expect(result.current[0]).toBe('idle')

    await act(async () => {
      await result.current[1]('hello')
    })
    expect(result.current[0]).toBe('copied')
    expect(navigator.clipboard.writeText).toHaveBeenCalledWith('hello')
  })

  it('reports error status when the clipboard API rejects', async () => {
    vi.mocked(navigator.clipboard.writeText).mockRejectedValueOnce(new Error('denied'))
    const { result } = renderHook(() => useCopyToClipboard())

    await act(async () => {
      await result.current[1]('hello')
    })
    expect(result.current[0]).toBe('error')
  })
})
