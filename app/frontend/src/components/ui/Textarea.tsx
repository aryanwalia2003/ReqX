import { forwardRef } from 'react'
import type { TextareaHTMLAttributes } from 'react'

import { cn } from '@/lib/cn'

export interface TextareaProps extends TextareaHTMLAttributes<HTMLTextAreaElement> {
  invalid?: boolean
  /** JSON/body editors ke liye monospace font. */
  mono?: boolean
}

export const Textarea = forwardRef<HTMLTextAreaElement, TextareaProps>(function Textarea(
  { invalid, mono, disabled, className, ...props },
  ref,
) {
  return (
    <textarea
      ref={ref}
      disabled={disabled}
      aria-invalid={invalid}
      className={cn(
        'bg-surface min-h-20 w-full rounded-md border px-3 py-2 text-sm',
        'text-fg placeholder:text-fg-subtle outline-none',
        'border-border focus:border-border-strong focus:outline-ring focus:outline-2 focus:outline-offset-1',
        invalid && 'border-fg-subtle border-dashed',
        disabled && 'pointer-events-none opacity-50',
        mono && 'font-mono',
        className,
      )}
      {...props}
    />
  )
})
