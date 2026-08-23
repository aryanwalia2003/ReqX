import { act, renderHook } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'

import { useLocalStorage } from '@/hooks/useLocalStorage'

describe('useLocalStorage', () => {
  afterEach(() => window.localStorage.clear())

  it('reads the initial value and persists updates', () => {
    const { result } = renderHook(() => useLocalStorage('theme', 'dark'))
    expect(result.current[0]).toBe('dark')

    act(() => result.current[1]('light'))
    expect(result.current[0]).toBe('light')
    expect(window.localStorage.getItem('theme')).toBe('"light"')
  })

  it('supports a functional updater', () => {
    window.localStorage.setItem('count', '1')
    const { result } = renderHook(() => useLocalStorage('count', 0))

    act(() => result.current[1]((prev) => prev + 1))
    expect(result.current[0]).toBe(2)
  })

  it('falls back to the initial value when stored JSON is invalid', () => {
    window.localStorage.setItem('broken', '{not json')
    const { result } = renderHook(() => useLocalStorage('broken', 'fallback'))
    expect(result.current[0]).toBe('fallback')
  })
})
