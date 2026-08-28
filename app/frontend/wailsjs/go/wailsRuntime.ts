// Shared by every hand-written stand-in under wailsjs/go/services/*.ts (see
// RequestService.ts for why those exist). Not itself something
// `wails generate module` produces — delete it along with the rest of
// wailsjs/ once that command has actually run.
//
// Pass the call already-invoked via optional chaining, e.g.
//   callWailsMethod(window.go?.services?.RequestService?.Send(input))
// so a missing window.go short-circuits to `undefined` (optional chaining
// never throws) instead of throwing synchronously — which would escape
// toResult's try/catch, since the throw happens while evaluating the
// argument, before toResult's own try block ever runs.
export function callWailsMethod<T>(result: Promise<T> | undefined): Promise<T> {
  return result ?? Promise.reject(new Error('Not running inside the Wails desktop app.'))
}
