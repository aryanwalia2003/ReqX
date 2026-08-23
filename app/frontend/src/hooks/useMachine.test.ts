import { act, renderHook } from '@testing-library/react'
import { describe, expect, it } from 'vitest'

import { useMachine } from '@/hooks/useMachine'
import { createMachine } from '@/lib/machine'

type ToggleState = 'off' | 'on'
type ToggleEvent = { type: 'TOGGLE' }

const toggleMachine = createMachine<undefined, ToggleState, ToggleEvent>({
  initial: 'off',
  context: undefined,
  states: {
    off: { on: { TOGGLE: 'on' } },
    on: { on: { TOGGLE: 'off' } },
  },
})

describe('useMachine', () => {
  it('starts at the machine initial state', () => {
    const { result } = renderHook(() => useMachine(toggleMachine))
    expect(result.current.state).toBe('off')
    expect(result.current.matches('off')).toBe(true)
  })

  it('applies a transition on send', () => {
    const { result } = renderHook(() => useMachine(toggleMachine))

    act(() => result.current.send({ type: 'TOGGLE' }))
    expect(result.current.state).toBe('on')

    act(() => result.current.send({ type: 'TOGGLE' }))
    expect(result.current.state).toBe('off')
  })

  it('reports which events the current state accepts', () => {
    const { result } = renderHook(() => useMachine(toggleMachine))
    expect(result.current.can('TOGGLE')).toBe(true)
  })

  it('keeps context in sync with the state', () => {
    type FetchState = 'idle' | 'loading' | 'success'
    type FetchEvent = { type: 'FETCH' } | { type: 'RESOLVE'; data: string }

    const fetchMachine = createMachine<{ data?: string }, FetchState, FetchEvent>({
      initial: 'idle',
      context: {},
      states: {
        idle: { on: { FETCH: 'loading' } },
        loading: {
          on: {
            RESOLVE: { target: 'success', assign: (ctx, event) => ({ ...ctx, data: event.data }) },
          },
        },
        success: {},
      },
    })

    const { result } = renderHook(() => useMachine(fetchMachine))
    act(() => result.current.send({ type: 'FETCH' }))
    act(() => result.current.send({ type: 'RESOLVE', data: 'hello' }))

    expect(result.current.state).toBe('success')
    expect(result.current.context).toEqual({ data: 'hello' })
  })
})
