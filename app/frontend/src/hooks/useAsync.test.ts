import { act, renderHook } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'

import { useAsync } from '@/hooks/useAsync'
import { err, ok } from '@/lib/result'
import type { Result } from '@/lib/result'

describe('useAsync', () => {
  it('tracks loading state and stores the resolved value', async () => {
    const { result } = renderHook(() => useAsync((n: number) => Promise.resolve(ok(n * 2))))

    let pending: Promise<unknown> = Promise.resolve()
    act(() => {
      pending = result.current.run(21)
    })
    expect(result.current.isLoading).toBe(true)

    await act(async () => {
      await pending
    })
    expect(result.current.isLoading).toBe(false)
    expect(result.current.data).toBe(42)
    expect(result.current.error).toBeUndefined()
  })

  it('stores the error without touching data on failure', async () => {
    const boom = new Error('boom')
    const { result } = renderHook(() => useAsync(() => Promise.resolve(err(boom))))

    await act(async () => {
      await result.current.run()
    })

    expect(result.current.error).toBe(boom)
    expect(result.current.data).toBeUndefined()
  })

  it('ignores a stale call when a newer one resolves first', async () => {
    const resolvers: Array<(value: Result<number>) => void> = []
    const asyncFn = vi.fn(
      () =>
        new Promise<Result<number>>((resolve) => {
          resolvers.push(resolve)
        }),
    )
    const { result } = renderHook(() => useAsync(asyncFn))

    let firstCall: Promise<unknown> = Promise.resolve()
    let secondCall: Promise<unknown> = Promise.resolve()
    act(() => {
      firstCall = result.current.run()
      secondCall = result.current.run()
    })

    await act(async () => {
      resolvers[1]?.(ok(2))
      resolvers[0]?.(ok(1))
      await Promise.all([firstCall, secondCall])
    })

    expect(result.current.data).toBe(2)
  })

  it('invokes the latest onError callback without forcing useCallback on the caller', async () => {
    const boom = new Error('boom')
    const firstOnError = vi.fn()
    const secondOnError = vi.fn()
    const { result, rerender } = renderHook(
      ({ onError }) => useAsync(() => Promise.resolve(err(boom)), { onError }),
      { initialProps: { onError: firstOnError } },
    )

    rerender({ onError: secondOnError })
    await act(async () => {
      await result.current.run()
    })

    expect(firstOnError).not.toHaveBeenCalled()
    expect(secondOnError).toHaveBeenCalledWith(boom)
  })
})
