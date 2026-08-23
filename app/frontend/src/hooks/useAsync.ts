import { useCallback, useEffect, useRef, useState } from 'react'

import type { AppError } from '@/lib/errors'
import type { Result } from '@/lib/result'

export interface UseAsyncState<T, E> {
  data: T | undefined
  error: E | undefined
  isLoading: boolean
}

export interface UseAsyncOptions<E> {
  /** Har failure par chalta hai — e.g. `reportError` ya toast dikhane ke liye. */
  onError?: (error: E) => void
}

export interface UseAsyncReturn<T, E, Args extends unknown[]> extends UseAsyncState<T, E> {
  run: (...args: Args) => Promise<Result<T, E>>
}

/**
 * Koi bhi `Result<T, E>` return karne wala async fn (jaise `toResult` +
 * Wails binding) ko loading/data/error state me wrap karta hai. Stale call
 * — yaani ek naya `run` beech me shuru ho jaaye — apne aap ignore ho jata hai.
 *
 * @example
 * const { data, error, isLoading, run } = useAsync((name: string) =>
 *   toResult(window.go.services.ExampleService.ping({ name })),
 * )
 * useEffect(() => { void run('dev') }, [run])
 */
export function useAsync<T, E = AppError, Args extends unknown[] = []>(
  asyncFn: (...args: Args) => Promise<Result<T, E>>,
  options?: UseAsyncOptions<E>,
): UseAsyncReturn<T, E, Args> {
  const [state, setState] = useState<UseAsyncState<T, E>>({
    data: undefined,
    error: undefined,
    isLoading: false,
  })

  const fnRef = useRef(asyncFn)
  useEffect(() => {
    fnRef.current = asyncFn
  }, [asyncFn])

  const onErrorRef = useRef(options?.onError)
  useEffect(() => {
    onErrorRef.current = options?.onError
  }, [options?.onError])

  const callId = useRef(0)

  const run = useCallback(async (...args: Args): Promise<Result<T, E>> => {
    const id = ++callId.current
    setState((prev) => ({ ...prev, isLoading: true }))

    const result = await fnRef.current(...args)
    if (id === callId.current) {
      if (result.ok) {
        setState({ data: result.value, error: undefined, isLoading: false })
      } else {
        setState({ data: undefined, error: result.error, isLoading: false })
        onErrorRef.current?.(result.error)
      }
    }
    return result
  }, [])

  return { ...state, run }
}
