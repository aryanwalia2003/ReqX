import { useCallback, useReducer } from 'react'

import type { Machine } from '@/lib/machine'

export interface UseMachineReturn<Ctx, State extends string, Event extends { type: string }> {
  state: State
  context: Ctx
  send: (event: Event) => void
  /** Current state exactly ye value hai kya — conditional render me `if` ke bajaye. */
  matches: (state: State) => boolean
  /** Current state se ye event handle hota hai kya — Button disable karne ke liye. */
  can: (eventType: Event['type']) => boolean
}

/**
 * `createMachine` (`src/lib/machine.ts`) ke config ko React state se jodta.
 * Machine khud stateless hai — module scope me ek baar banao (ya props par
 * depend kare to `useState(() => createMachine(...))`), phir yahan wire karo.
 *
 * @example
 * const { state, context, send, matches } = useMachine(fetchMachine)
 * if (matches('loading')) return <Spinner />
 */
export function useMachine<Ctx, State extends string, Event extends { type: string }>(
  machine: Machine<Ctx, State, Event>,
): UseMachineReturn<Ctx, State, Event> {
  const [snapshot, dispatch] = useReducer(machine.transition, machine.initial)

  const send = useCallback((event: Event) => dispatch(event), [])

  return {
    state: snapshot.state,
    context: snapshot.context,
    send,
    matches: (state) => snapshot.state === state,
    can: (eventType) => machine.can(snapshot.state, eventType),
  }
}
