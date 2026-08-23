# Framework-agnostic utilities

Plain TypeScript, no React imports (except where noted) — the kind of thing
any feature might need, with no UI attached.

- **`cn`** — merge conditional classnames + resolve conflicting Tailwind
  utilities.
- **`result`** — the project's error-handling convention: `Result<T, E>`
  instead of throwing. `toResult(promise)` wraps a Wails binding call and
  normalizes any rejection into an `AppError` (see below).
- **`errors`** — `AppError`/`ErrorKind` (mirrors the Go side's
  `internal/errs.Kind`) + `toAppError`, which normalizes any thrown/rejected
  value into that shape, and `getErrorMessage`, a safe display string with a
  fallback for an empty message.
- **`reportError`** — the single entry point for reporting an error from
  anywhere (not just React components): `reportError(cause, context?)`.
  Logs to the console and notifies subscribers registered via
  `onErrorReported` — this is what `<ErrorBoundary>` and the root's
  `installGlobalErrorHandlers` funnel into, and where a future telemetry
  sink would subscribe.
- **`machine`** — `createMachine`, a small dependency-free finite state
  machine: explicit states + the events each one accepts, with an optional
  `context` payload updated via `assign`. Wire one into a component with
  `useMachine` (`src/hooks/README.md`). Shared, cross-feature machine
  configs live under `machines/` — see `src/lib/machines/README.md`.

## Error handling, end to end

1. A Wails call: `toResult(window.go.services.X.y(...))` → `Result<T, AppError>`.
2. A component: `useAsync` wraps that in `{ data, error, isLoading, run }`
   (see `src/hooks/README.md`); render `error` with `<Alert variant="danger">`
   using `getErrorMessage(error)`.
3. A render crash: `<ErrorBoundary>` (see `src/components/ui/README.md`)
   catches it, shows a fallback, and reports it.
4. Anything else that slips through (an unhandled rejection, a stray
   `window.onerror`): `installGlobalErrorHandlers` (wired once in
   `main.tsx`) catches it and reports it — `<GlobalErrorToasts>` turns every
   report into a toast so nothing fails silently.
