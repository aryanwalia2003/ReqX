// Hand-written stand-in for what `wails generate module` would generate —
// the Wails CLI isn't installed in this environment (see app/README.md), so
// that command has never actually run here. Delete this file once it has;
// the real generated file (same shape) is a drop-in replacement.
//
// window.go.services.* is injected at runtime by the live Go backend (per
// app/wailsapp's BindTargets) when running inside the actual Wails webview —
// this only proxies to it, typed via src/types/wails.d.ts.
import type { SendRequestInput, SendRequestOutput } from '@/features/send-request/types'

export function Send(input: SendRequestInput): Promise<SendRequestOutput> {
  // window.go only exists inside the real Wails webview — guard so a plain
  // browser (e.g. `npm run dev` standalone, per app/frontend/README.md)
  // gets a rejected promise through the normal Result<T> path instead of a
  // synchronous throw that would escape toResult's try/catch.
  if (typeof window === 'undefined' || !window.go?.services?.RequestService) {
    return Promise.reject(new Error('Not running inside the Wails desktop app.'))
  }
  return window.go.services.RequestService.Send(input)
}
