import { act, renderHook } from '@testing-library/react'
import { beforeEach, describe, expect, it } from 'vitest'

import { useRecentCollections } from './useRecentCollections'

const demo = { name: 'Demo', requests: [] }
const other = { name: 'Other', requests: [{ name: 'Ping', method: 'GET', url: 'https://x' }] }

describe('useRecentCollections', () => {
  beforeEach(() => localStorage.clear())

  it('starts empty and adds an entry via upsert', () => {
    const { result } = renderHook(() => useRecentCollections())
    expect(result.current.recents).toEqual([])

    act(() => result.current.upsert('/a/demo.json', demo))

    expect(result.current.recents).toHaveLength(1)
    expect(result.current.recents[0]).toMatchObject({ path: '/a/demo.json', collection: demo })
  })

  it('moves an existing path to the front instead of duplicating it', () => {
    const { result } = renderHook(() => useRecentCollections())

    act(() => result.current.upsert('/a/demo.json', demo))
    act(() => result.current.upsert('/b/other.json', other))
    act(() => result.current.upsert('/a/demo.json', demo))

    expect(result.current.recents.map((r) => r.path)).toEqual(['/a/demo.json', '/b/other.json'])
  })

  it('removes an entry', () => {
    const { result } = renderHook(() => useRecentCollections())

    act(() => result.current.upsert('/a/demo.json', demo))
    act(() => result.current.remove('/a/demo.json'))

    expect(result.current.recents).toEqual([])
  })

  it('caps the list at 8 entries', () => {
    const { result } = renderHook(() => useRecentCollections())

    for (let i = 0; i < 10; i++) {
      act(() => result.current.upsert(`/c${i}.json`, demo))
    }

    expect(result.current.recents).toHaveLength(8)
    expect(result.current.recents[0]?.path).toBe('/c9.json')
  })

  it('persists to localStorage and a fresh hook instance picks it up', () => {
    const { result, unmount } = renderHook(() => useRecentCollections())
    act(() => result.current.upsert('/a/demo.json', demo))
    unmount()

    const { result: fresh } = renderHook(() => useRecentCollections())
    expect(fresh.current.recents).toHaveLength(1)
    expect(fresh.current.recents[0]?.path).toBe('/a/demo.json')
  })
})
