import { toAppError } from '@/lib/errors'
import type { AppError } from '@/lib/errors'

export interface ErrorReport {
  error: AppError
  context?: string
}

type Listener = (report: ErrorReport) => void

const listeners = new Set<Listener>()

/** Naya subscriber register karta — telemetry/logging aage yahin judega. */
export function onErrorReported(listener: Listener): () => void {
  listeners.add(listener)
  return () => listeners.delete(listener)
}

/** App me kahin se bhi error report karne ka single entry point. */
export function reportError(cause: unknown, context?: string): void {
  const report: ErrorReport = { error: toAppError(cause), context }
  console.error(context ? `[${context}]` : '[error]', report.error.message, report.error.cause)
  listeners.forEach((listener) => listener(report))
}
