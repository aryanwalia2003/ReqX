import { reportError } from '@/lib/reportError'

/**
 * `window.onerror` + unhandled promise rejections ko `reportError` se jodta.
 * Ye un errors ke liye hai jo kisi `toResult`/try-catch se bach jaate hain.
 */
export function installGlobalErrorHandlers(): () => void {
  const onError = (event: ErrorEvent) => reportError(event.error ?? event.message, 'window.error')
  const onRejection = (event: PromiseRejectionEvent) =>
    reportError(event.reason, 'unhandledrejection')

  window.addEventListener('error', onError)
  window.addEventListener('unhandledrejection', onRejection)

  return () => {
    window.removeEventListener('error', onError)
    window.removeEventListener('unhandledrejection', onRejection)
  }
}
