import { createContext, useContext, useEffect, useRef } from 'react'
import type { HTMLAttributes, ReactNode } from 'react'

import { Button } from '@/components/ui/Button'
import { XIcon } from '@/components/ui/icons'
import { cn } from '@/lib/cn'

const DialogContext = createContext<{ close: () => void } | null>(null)

export interface DialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  children: ReactNode
  className?: string
}

/**
 * Native <dialog> pe bana modal — apna backdrop/focus-trap/Escape sab
 * browser handle karta hai, koi portal library nahi chahiye.
 */
export function Dialog({ open, onOpenChange, children, className }: DialogProps) {
  const ref = useRef<HTMLDialogElement>(null)

  useEffect(() => {
    const dialog = ref.current
    if (!dialog) return
    if (open && !dialog.open) dialog.showModal()
    if (!open && dialog.open) dialog.close()
  }, [open])

  return (
    <DialogContext.Provider value={{ close: () => onOpenChange(false) }}>
      <dialog
        ref={ref}
        onClose={() => onOpenChange(false)}
        onCancel={() => onOpenChange(false)}
        onClick={(e) => {
          if (e.target === ref.current) onOpenChange(false)
        }}
        className={cn(
          // Tailwind preflight margin:0 UA ka auto-centering todta hai — m-auto se wapas.
          'border-border bg-surface text-fg m-auto w-full max-w-md rounded-lg border p-0',
          className,
        )}
      >
        {children}
      </dialog>
    </DialogContext.Provider>
  )
}

export function DialogHeader({ className, children, ...props }: HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn('border-border flex items-start justify-between gap-3 border-b p-4', className)}
      {...props}
    >
      <div className="flex flex-col gap-1">{children}</div>
      <DialogClose />
    </div>
  )
}

export function DialogTitle({ className, ...props }: HTMLAttributes<HTMLHeadingElement>) {
  return <h2 className={cn('text-fg text-sm font-semibold', className)} {...props} />
}

export function DialogDescription({ className, ...props }: HTMLAttributes<HTMLParagraphElement>) {
  return <p className={cn('text-fg-muted text-sm', className)} {...props} />
}

export function DialogContent({ className, ...props }: HTMLAttributes<HTMLDivElement>) {
  return <div className={cn('p-4', className)} {...props} />
}

export function DialogFooter({ className, ...props }: HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn('border-border flex items-center justify-end gap-2 border-t p-4', className)}
      {...props}
    />
  )
}

export function DialogClose({ className }: { className?: string }) {
  const ctx = useContext(DialogContext)
  return (
    <Button
      type="button"
      variant="ghost"
      size="icon"
      aria-label="Close"
      onClick={() => ctx?.close()}
      className={cn('-m-1', className)}
    >
      <XIcon />
    </Button>
  )
}
