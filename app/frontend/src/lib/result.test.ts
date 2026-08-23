import { describe, expect, it } from 'vitest'

import { err, ok, toResult } from '@/lib/result'

describe('result', () => {
  it('wraps a resolved promise as ok', async () => {
    expect(await toResult(Promise.resolve(42))).toEqual(ok(42))
  })

  it('normalizes a rejected promise into an AppError', async () => {
    const result = await toResult(Promise.reject(new Error('boom')))
    expect(result).toEqual(err({ kind: 'unknown', message: 'boom', cause: expect.any(Error) }))
  })

  it('parses a JSON-encoded backend error into its kind', async () => {
    const raw = JSON.stringify({ kind: 'not_found', message: 'collection missing' })
    const result = await toResult(Promise.reject(raw))
    expect(result).toEqual(err({ kind: 'not_found', message: 'collection missing', cause: raw }))
  })
})
