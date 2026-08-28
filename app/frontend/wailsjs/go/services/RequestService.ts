// Hand-written stand-in for what `wails generate module` would generate —
// the Wails CLI isn't installed in this environment (see app/README.md), so
// that command has never actually run here. Delete this file once it has;
// the real generated file (same shape) is a drop-in replacement.
//
// window.go.services.* is injected at runtime by the live Go backend (per
// app/wailsapp's BindTargets) when running inside the actual Wails webview —
// this only proxies to it, typed via src/types/wails.d.ts.
import { callWailsMethod } from '@wails/go/wailsRuntime'

import type { SendRequestInput, SendRequestOutput } from '@/features/send-request/types'

export function Send(input: SendRequestInput): Promise<SendRequestOutput> {
  return callWailsMethod(window.go?.services?.RequestService?.Send(input))
}
