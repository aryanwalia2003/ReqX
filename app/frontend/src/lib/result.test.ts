import { describe, expect, it } from 'vitest'

import { err, ok, toResult } from '@/lib/result'

describe('result', () => {
  it('wraps a resolved promise as ok', async () => {
    expect(await toResult(Promise.resolve(42))).toEqual(ok(42))
  })

  it('wraps a rejected promise as err', async () => {
    const result = await toResult(Promise.reject(new Error('boom')))
    expect(result).toEqual(err(new Error('boom')))
  })
})
