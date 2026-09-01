// Hand-written stand-in for what `wails generate module` would generate for
// the runtime's Events API — see wailsjs/go/services/RequestService.ts for
// why these files exist and when to delete them (once real codegen runs).
//
// Wails exposes this as window.runtime.EventsOn/EventsOff at runtime; a
// missing window.runtime (plain browser / dev mode) makes On() a no-op that
// returns a no-op unsubscribe, same "degrade gracefully outside Wails"
// contract as callWailsMethod.
export function EventsOn(eventName: string, callback: (...data: unknown[]) => void): () => void {
  const off = window.runtime?.EventsOn(eventName, callback)
  return off ?? (() => {})
}

export function EventsOff(eventName: string): void {
  window.runtime?.EventsOff(eventName)
}
