import { createContext, useCallback, useContext, useRef, useState } from 'react'
import type { ReactNode } from 'react'
import { createPortal } from 'react-dom'

import { XIcon } from '@/components/ui/icons'
import { cn } from '@/lib/cn'

export type ToastVariant = 'default' | 'success' | 'warning' | 'danger'

export interface ToastOptions {
  title: string
  description?: string
  variant?: ToastVariant
  /** ms; default 4000. */
  duration?: number
}

interface ToastItem extends ToastOptions {
  id: number
}

interface ToastContextValue {
  toast: (options: ToastOptions) => void
}

const ToastContext = createContext<ToastContextValue | null>(null)

const DEFAULT_DURATION = 4000

let nextId = 0

/** Root me ek baar wrap karo — phir kahin se bhi useToast() chalega. */
export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<ToastItem[]>([])
  const timers = useRef(new Map<number, ReturnType<typeof setTimeout>>())

  const dismiss = useCallback((id: number) => {
    setToasts((prev) => prev.filter((t) => t.id !== id))
    const timer = timers.current.get(id)
    if (timer) {
      clearTimeout(timer)
      timers.current.delete(id)
    }
  }, [])

  const toast = useCallback(
    (options: ToastOptions) => {
      const id = ++nextId
      setToasts((prev) => [...prev, { id, ...options }])
      const timer = setTimeout(() => dismiss(id), options.duration ?? DEFAULT_DURATION)
      timers.current.set(id, timer)
    },
    [dismiss],
  )

  return (
    <ToastContext.Provider value={{ toast }}>
      {children}
      {createPortal(
        <div className="fixed right-4 bottom-4 z-50 flex w-80 flex-col gap-2">
          {toasts.map((t) => (
            <div
              key={t.id}
              role="status"
              className={cn(
                'bg-surface-2 flex items-start gap-2 rounded-lg border p-3 text-sm shadow-lg',
                t.variant === 'warning' || t.variant === 'danger'
                  ? 'border-border-strong border-dashed'
                  : 'border-border',
              )}
            >
              <div className="flex flex-1 flex-col gap-0.5">
                <p className="text-fg font-medium">{t.title}</p>
                {t.description && <p className="text-fg-muted">{t.description}</p>}
              </div>
              <button
                type="button"
                aria-label="Dismiss"
                onClick={() => dismiss(t.id)}
                className="text-fg-subtle hover:text-fg"
              >
                <XIcon className="size-4" />
              </button>
            </div>
          ))}
        </div>,
        document.body,
      )}
    </ToastContext.Provider>
  )
}

// eslint-disable-next-line react-refresh/only-export-components -- hook context ke saath rehna zaroori hai
export function useToast() {
  const ctx = useContext(ToastContext)
  if (!ctx) throw new Error('useToast must be used inside <ToastProvider>')
  return ctx.toast
}
