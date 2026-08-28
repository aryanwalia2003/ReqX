/**
 * Two stand-ins for what `wails generate module` produces once the Wails
 * CLI actually runs (not installed in this environment — see app/README.md):
 *
 * 1. The `declare global { interface Window { go: ... } }` block below —
 *    ADD to it (don't replace) as more services get bound in
 *    app/wailsapp/app_struct.go, so it stays the one place declaring the
 *    shape of window.go.
 * 2. wailsjs/go/services/<Service>.ts — one hand-written proxy per bound
 *    service (e.g. wailsjs/go/services/RequestService.ts), each just
 *    calling through to window.go.services.<Service>.<Method>, typed
 *    against the global declared here.
 *
 * DELETE all of it — this block and the whole wailsjs/ folder — once
 * `wails generate module` has actually run; its own output is a drop-in
 * replacement with the same shape.
 */
import type { SendRequestInput, SendRequestOutput } from '@/features/send-request/types'

declare global {
  interface Window {
    go: {
      services: {
        RequestService: {
          Send(input: SendRequestInput): Promise<SendRequestOutput>
        }
      }
    }
  }
}

// Fallback for any @wails/* import not covered by a real wailsjs/ file yet.
declare module '@wails/*' {
  const value: unknown
  export default value
}
