import { forwardRef } from 'react'
import type { InputHTMLAttributes } from 'react'

import { cn } from '@/lib/cn'

export type RadioProps = Omit<InputHTMLAttributes<HTMLInputElement>, 'type' | 'size'>

export const Radio = forwardRef<HTMLInputElement, RadioProps>(function Radio(
  { className, ...props },
  ref,
) {
  return (
    <span className={cn('relative inline-flex size-4 shrink-0', className)}>
      <input
        ref={ref}
        type="radio"
        className="peer border-border bg-surface checked:border-fg focus-visible:outline-ring size-4 shrink-0 cursor-pointer appearance-none rounded-full border focus-visible:outline-2 focus-visible:outline-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
        {...props}
      />
      <span className="bg-fg pointer-events-none absolute inset-0 m-auto hidden size-1.5 rounded-full peer-checked:block" />
    </span>
  )
})
