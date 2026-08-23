import { describe, expect, it } from 'vitest'

import { getErrorMessage, toAppError } from '@/lib/errors'

describe('toAppError', () => {
  it('wraps a plain string message (todays Wails behavior)', () => {
    expect(toAppError('boom')).toEqual({ kind: 'unknown', message: 'boom', cause: 'boom' })
  })

  it('extracts kind + message from a JSON-encoded string', () => {
    const raw = JSON.stringify({ kind: 'invalid_input', message: 'name is required' })
    expect(toAppError(raw)).toEqual({
      kind: 'invalid_input',
      message: 'name is required',
      cause: raw,
    })
  })

  it('falls back to unknown for an unrecognized kind', () => {
    const raw = JSON.stringify({ kind: 'made_up', message: 'oops' })
    expect(toAppError(raw)).toEqual({ kind: 'unknown', message: 'oops', cause: raw })
  })

  it('wraps a native Error', () => {
    const error = new Error('native failure')
    expect(toAppError(error)).toEqual({ kind: 'unknown', message: 'native failure', cause: error })
  })

  it('wraps a totally unrecognized thrown value', () => {
    expect(toAppError(42)).toEqual({ kind: 'unknown', message: '42', cause: 42 })
  })
})

describe('getErrorMessage', () => {
  it('returns the message when present', () => {
    expect(getErrorMessage({ kind: 'unknown', message: 'oops' })).toBe('oops')
  })

  it('falls back for an empty message', () => {
    expect(getErrorMessage({ kind: 'unknown', message: '' })).toBe(
      'Something went wrong. Try again.',
    )
  })
})
