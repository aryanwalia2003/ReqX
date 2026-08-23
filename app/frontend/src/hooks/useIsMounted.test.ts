import { renderHook } from '@testing-library/react'
import { describe, expect, it } from 'vitest'

import { useIsMounted } from '@/hooks/useIsMounted'

describe('useIsMounted', () => {
  it('is true while mounted and false after unmount', () => {
    const { result, unmount } = renderHook(() => useIsMounted())
    expect(result.current()).toBe(true)

    unmount()
    expect(result.current()).toBe(false)
  })
})
