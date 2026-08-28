import { toResult } from '@/lib/result'
import { Send } from '@wails/go/services/RequestService'

import type { SendRequestInput } from './types'

/** Ek ad-hoc HTTP request bhejta — CLI ke `reqx req` jaisa hi backend path. */
export function sendRequest(input: SendRequestInput) {
  return toResult(Send(input))
}
