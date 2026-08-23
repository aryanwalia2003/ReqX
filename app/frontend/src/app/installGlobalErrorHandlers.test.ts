import { afterEach, describe, expect, it, vi } from 'vitest'

import { installGlobalErrorHandlers } from '@/app/installGlobalErrorHandlers'
import { onErrorReported } from '@/lib/reportError'

describe('installGlobalErrorHandlers', () => {
  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('funnels window error events into reportError', () => {
    vi.spyOn(console, 'error').mockImplementation(() => {})
    const listener = vi.fn()
    const unsubscribeListener = onErrorReported(listener)
    const uninstall = installGlobalErrorHandlers()

    window.dispatchEvent(new ErrorEvent('error', { error: new Error('boom'), message: 'boom' }))
    expect(listener).toHaveBeenCalledOnce()

    uninstall()
    unsubscribeListener()
  })

  it('funnels unhandledrejection events into reportError', () => {
    vi.spyOn(console, 'error').mockImplementation(() => {})
    const listener = vi.fn()
    const unsubscribeListener = onErrorReported(listener)
    const uninstall = installGlobalErrorHandlers()

    const event = new Event('unhandledrejection') as PromiseRejectionEvent
    Object.defineProperty(event, 'reason', { value: 'boom' })
    window.dispatchEvent(event)
    expect(listener).toHaveBeenCalledOnce()

    uninstall()
    unsubscribeListener()
  })

  it('stops listening after uninstall', () => {
    vi.spyOn(console, 'error').mockImplementation(() => {})
    const listener = vi.fn()
    const unsubscribeListener = onErrorReported(listener)
    const uninstall = installGlobalErrorHandlers()
    uninstall()

    const event = new Event('unhandledrejection') as PromiseRejectionEvent
    Object.defineProperty(event, 'reason', { value: 'boom' })
    window.dispatchEvent(event)
    expect(listener).not.toHaveBeenCalled()

    unsubscribeListener()
  })
})
