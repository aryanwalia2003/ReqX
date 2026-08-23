import { useCallback, useState } from 'react'

type SetValue<T> = T | ((prev: T) => T)

/** localStorage-backed state — JSON serialize/parse apne aap ho jata hai. */
export function useLocalStorage<T>(key: string, initialValue: T) {
  const [value, setValue] = useState<T>(() => {
    try {
      const item = window.localStorage.getItem(key)
      return item ? (JSON.parse(item) as T) : initialValue
    } catch {
      return initialValue
    }
  })

  const set = useCallback(
    (next: SetValue<T>) => {
      setValue((prev) => {
        const resolved = next instanceof Function ? next(prev) : next
        try {
          window.localStorage.setItem(key, JSON.stringify(resolved))
        } catch {
          // quota full ya storage disabled — state phir bhi update hota hai
        }
        return resolved
      })
    },
    [key],
  )

  return [value, set] as const
}
