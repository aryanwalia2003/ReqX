import { describe, expect, it } from 'vitest'

import { createMachine } from '@/lib/machine'

type FetchState = 'idle' | 'loading' | 'success' | 'error'
type FetchEvent =
  | { type: 'FETCH' }
  | { type: 'RESOLVE'; data: string }
  | { type: 'REJECT'; message: string }
  | { type: 'RETRY' }

function createFetchMachine() {
  return createMachine<{ data?: string; error?: string }, FetchState, FetchEvent>({
    initial: 'idle',
    context: {},
    states: {
      idle: { on: { FETCH: 'loading' } },
      loading: {
        on: {
          RESOLVE: { target: 'success', assign: (ctx, event) => ({ ...ctx, data: event.data }) },
          REJECT: { target: 'error', assign: (ctx, event) => ({ ...ctx, error: event.message }) },
        },
      },
      success: { on: { FETCH: 'loading' } },
      error: { on: { RETRY: 'loading' } },
    },
  })
}

describe('createMachine', () => {
  it('starts at the configured initial snapshot', () => {
    const machine = createFetchMachine()
    expect(machine.initial).toEqual({ state: 'idle', context: {} })
  })

  it('follows a string-shorthand transition without touching context', () => {
    const machine = createFetchMachine()
    const next = machine.transition(machine.initial, { type: 'FETCH' })
    expect(next).toEqual({ state: 'loading', context: {} })
  })

  it('assigns context via an object transition', () => {
    const machine = createFetchMachine()
    const loading = machine.transition(machine.initial, { type: 'FETCH' })
    const resolved = machine.transition(loading, { type: 'RESOLVE', data: 'hello' })
    expect(resolved).toEqual({ state: 'success', context: { data: 'hello' } })
  })

  it('takes the error branch and stores the message', () => {
    const machine = createFetchMachine()
    const loading = machine.transition(machine.initial, { type: 'FETCH' })
    const errored = machine.transition(loading, { type: 'REJECT', message: 'network down' })
    expect(errored).toEqual({ state: 'error', context: { error: 'network down' } })
  })

  it('allows recovering only via the transition the errored state declares', () => {
    const machine = createFetchMachine()
    const errored = { state: 'error' as const, context: { error: 'network down' } }
    const retried = machine.transition(errored, { type: 'RETRY' })
    expect(retried).toEqual({ state: 'loading', context: { error: 'network down' } })
  })

  it('ignores an event the current state does not declare, preserving identity', () => {
    const machine = createFetchMachine()
    const next = machine.transition(machine.initial, { type: 'RETRY' })
    expect(next).toBe(machine.initial)
  })

  it('reports whether a state can handle a given event', () => {
    const machine = createFetchMachine()
    expect(machine.can('idle', 'FETCH')).toBe(true)
    expect(machine.can('idle', 'RETRY')).toBe(false)
    expect(machine.can('error', 'RETRY')).toBe(true)
  })
})
