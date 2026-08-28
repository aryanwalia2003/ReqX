import { useCallback, useState } from 'react'

import type { Collection } from '@/features/collection-runner/types'

export interface RecentCollection {
  path: string
  collection: Collection
  lastOpenedAt: number
}

const STORAGE_KEY = 'reqx.recentCollections'
const MAX_RECENTS = 8

function loadFromStorage(): RecentCollection[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return []
    const parsed: unknown = JSON.parse(raw)
    return Array.isArray(parsed) ? (parsed as RecentCollection[]) : []
  } catch {
    // Private-mode/quota/disabled storage — start empty, still works this session.
    return []
  }
}

function saveToStorage(recents: RecentCollection[]) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(recents))
  } catch {
    // Same as above — nothing to recover, next upsert just tries again.
  }
}

/**
 * Postman-jaisi "recent collections" list — localStorage me persist, taaki
 * har baar Browse dobara na karna pade. Sirf path + last-loaded Collection
 * snapshot rakhta (koi backend/history.db involvement nahi — ye purely FE
 * convenience hai, run history se alag concern).
 */
export function useRecentCollections() {
  const [recents, setRecents] = useState<RecentCollection[]>(() => loadFromStorage())

  const upsert = useCallback((path: string, collection: Collection) => {
    setRecents((prev) => {
      const next = [
        { path, collection, lastOpenedAt: Date.now() },
        ...prev.filter((r) => r.path !== path),
      ].slice(0, MAX_RECENTS)
      saveToStorage(next)
      return next
    })
  }, [])

  const remove = useCallback((path: string) => {
    setRecents((prev) => {
      const next = prev.filter((r) => r.path !== path)
      saveToStorage(next)
      return next
    })
  }, [])

  return { recents, upsert, remove }
}
