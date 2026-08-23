import { renderHook } from '@testing-library/react'
import { describe, expect, it } from 'vitest'

import { usePrevious } from '@/hooks/usePrevious'

describe('usePrevious', () => {
  it('returns undefined on the first render, then the prior value', () => {
    const { result, rerender } = renderHook(({ value }) => usePrevious(value), {
      initialProps: { value: 1 },
    })
    expect(result.current).toBeUndefined()

    rerender({ value: 2 })
    expect(result.current).toBe(1)

    rerender({ value: 3 })
    expect(result.current).toBe(2)
  })
})
