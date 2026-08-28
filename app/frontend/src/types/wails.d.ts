/**
 * Two stand-ins for what `wails generate module` produces once the Wails
 * CLI actually runs (not installed in this environment — see app/README.md):
 *
 * 1. The `declare global { interface Window { go: ... } }` block below —
 *    ADD to it (don't replace) as more services get bound in
 *    app/wailsapp/app_struct.go, so it stays the one place declaring the
 *    shape of window.go. `go` and every service on it are optional — that's
 *    the honest type: outside the real Wails webview (a plain browser),
 *    window.go genuinely doesn't exist. wailsjs/go/wailsRuntime.ts's
 *    callWailsMethod is what turns a missing one into a rejected promise.
 * 2. wailsjs/go/services/<Service>.ts — one hand-written proxy per bound
 *    service (e.g. wailsjs/go/services/RequestService.ts), each just
 *    calling through to window.go.services.<Service>.<Method>, typed
 *    against the global declared here.
 *
 * DELETE all of it — this block and the whole wailsjs/ folder — once
 * `wails generate module` has actually run; its own output is a drop-in
 * replacement with the same shape.
 */
import type {
  Collection,
  Environment,
  Persona,
  RunCollectionInput,
  RunCollectionOutput,
} from '@/features/collection-runner/types'
import type { RunRow, StatRow } from '@/features/history/types'
import type { SendRequestInput, SendRequestOutput } from '@/features/send-request/types'

declare global {
  interface Window {
    go?: {
      services?: {
        RequestService?: {
          Send(input: SendRequestInput): Promise<SendRequestOutput>
        }
        CollectionService?: {
          PickFile(title: string, extension: string): Promise<string>
          Open(path: string): Promise<Collection>
          OpenEnvironment(path: string): Promise<Environment>
          OpenPersonas(path: string): Promise<Persona[]>
          Run(input: RunCollectionInput): Promise<RunCollectionOutput>
        }
        HistoryService?: {
          ListRuns(limit: number): Promise<RunRow[]>
          GetRunStats(runId: string): Promise<StatRow[]>
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
