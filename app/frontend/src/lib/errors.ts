export type ErrorKind =
  | 'invalid_input'
  | 'not_found'
  | 'forbidden'
  | 'unauthorized'
  | 'conflict'
  | 'database_error'
  | 'external_service_error'
  | 'internal_error'
  | 'unknown'

export interface AppError {
  kind: ErrorKind
  message: string
  cause?: unknown
}

const KNOWN_KINDS: readonly ErrorKind[] = [
  'invalid_input',
  'not_found',
  'forbidden',
  'unauthorized',
  'conflict',
  'database_error',
  'external_service_error',
  'internal_error',
]

function isErrorKind(value: unknown): value is ErrorKind {
  return typeof value === 'string' && (KNOWN_KINDS as readonly string[]).includes(value)
}

function hasStringMessage(value: unknown): value is { message: string; kind?: unknown } {
  return (
    !!value &&
    typeof value === 'object' &&
    'message' in value &&
    typeof (value as { message: unknown }).message === 'string'
  )
}

/**
 * Kisi bhi thrown/rejected value ko ek consistent AppError shape me badalta.
 * Aaj Wails sirf plain string message bhejta hai — kal Go side JSON-encoded
 * `{ kind, message }` bheje, to wo bhi yahin se parse ho jayega, bina kisi
 * call-site ko badle.
 */
export function toAppError(cause: unknown): AppError {
  if (typeof cause === 'string') {
    try {
      const parsed: unknown = JSON.parse(cause)
      if (hasStringMessage(parsed)) {
        return {
          kind: isErrorKind(parsed.kind) ? parsed.kind : 'unknown',
          message: parsed.message,
          cause,
        }
      }
    } catch {
      // JSON nahi hai — aaj ke Wails jaisa plain string message
    }
    return { kind: 'unknown', message: cause, cause }
  }

  if (hasStringMessage(cause)) {
    return {
      kind: isErrorKind(cause.kind) ? cause.kind : 'unknown',
      message: cause.message,
      cause,
    }
  }

  return { kind: 'unknown', message: String(cause), cause }
}

const FALLBACK_MESSAGE = 'Something went wrong. Try again.'

/** UI me dikhane layak message — khaali message ho to fallback text. */
export function getErrorMessage(error: AppError): string {
  return error.message.trim() || FALLBACK_MESSAGE
}
