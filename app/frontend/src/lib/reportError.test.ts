import { afterEach, describe, expect, it, vi } from 'vitest'

import { onErrorReported, reportError } from '@/lib/reportError'

describe('reportError', () => {
  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('notifies subscribers with a normalized AppError', () => {
    vi.spyOn(console, 'error').mockImplementation(() => {})
    const listener = vi.fn()
    const unsubscribe = onErrorReported(listener)

    reportError(new Error('boom'), 'test')

    expect(listener).toHaveBeenCalledWith({
      error: { kind: 'unknown', message: 'boom', cause: expect.any(Error) },
      context: 'test',
    })
    unsubscribe()
  })

  it('stops notifying after unsubscribing', () => {
    vi.spyOn(console, 'error').mockImplementation(() => {})
    const listener = vi.fn()
    const unsubscribe = onErrorReported(listener)
    unsubscribe()

    reportError('boom')
    expect(listener).not.toHaveBeenCalled()
  })
})
