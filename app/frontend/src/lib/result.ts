/**
 * The project's error-handling convention: functions that can fail return a
 * `Result` instead of throwing. Throwing is reserved for truly exceptional,
 * unrecoverable states (a programmer error, not an expected failure mode).
 *
 * Every Wails-bound Go method returns `(T, error)` — wrap the generated
 * binding call in `toResult` so the frontend handles it the same way
 * everywhere, instead of scattering try/catch around call sites.
 *
 * @example
 * const result = await toResult(window.go.services.ExampleService.ping({ name: 'dev' }))
 * if (!result.ok) {
 *   showError(result.error)
 *   return
 * }
 * console.log(result.value.message)
 */
export type Result<T, E = Error> = { ok: true; value: T } | { ok: false; error: E }

export function ok<T>(value: T): Result<T, never> {
  return { ok: true, value }
}

export function err<E>(error: E): Result<never, E> {
  return { ok: false, error }
}

/** Wraps a promise that may reject (e.g. a generated Wails binding call). */
export async function toResult<T>(promise: Promise<T>): Promise<Result<T>> {
  try {
    return ok(await promise)
  } catch (cause) {
    return err(cause instanceof Error ? cause : new Error(String(cause)))
  }
}
