import { forwardRef } from 'react'
import type { InputHTMLAttributes, ReactNode } from 'react'

import { cn } from '@/lib/cn'

export interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  invalid?: boolean
  leftElement?: ReactNode
  rightElement?: ReactNode
}

export const Input = forwardRef<HTMLInputElement, InputProps>(function Input(
  { invalid, leftElement, rightElement, disabled, className, ...props },
  ref,
) {
  return (
    <div
      className={cn(
        'bg-surface flex h-9 items-center gap-2 rounded-md border px-3 text-sm',
        'border-border focus-within:border-border-strong focus-within:outline-ring focus-within:outline-2 focus-within:outline-offset-1',
        invalid && 'border-fg-subtle border-dashed',
        disabled && 'pointer-events-none opacity-50',
        className,
      )}
    >
      {leftElement}
      <input
        ref={ref}
        disabled={disabled}
        aria-invalid={invalid}
        className="text-fg placeholder:text-fg-subtle w-full bg-transparent outline-none"
        {...props}
      />
      {rightElement}
    </div>
  )
})
