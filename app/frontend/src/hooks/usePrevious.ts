import { useState } from 'react'

interface PreviousState<T> {
  value: T
  previous: T | undefined
}

/** Pichhle render ki value deta hai — pehli render par undefined. */
export function usePrevious<T>(value: T): T | undefined {
  const [state, setState] = useState<PreviousState<T>>({ value, previous: undefined })

  if (value !== state.value) {
    setState({ value, previous: state.value })
  }

  return state.previous
}
