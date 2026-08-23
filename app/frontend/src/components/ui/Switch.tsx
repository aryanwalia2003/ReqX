import { forwardRef } from 'react'
import type { InputHTMLAttributes } from 'react'

import { cn } from '@/lib/cn'

export type SwitchProps = Omit<InputHTMLAttributes<HTMLInputElement>, 'type' | 'size'>

/** Pill toggle — native checkbox hi hai, checked/onChange se hi control hota hai. */
export const Switch = forwardRef<HTMLInputElement, SwitchProps>(function Switch(
  { className, ...props },
  ref,
) {
  return (
    <span className={cn('relative inline-flex h-5 w-9 shrink-0', className)}>
      <input
        ref={ref}
        type="checkbox"
        role="switch"
        className="peer border-border bg-surface-2 checked:bg-fg focus-visible:outline-ring size-full cursor-pointer appearance-none rounded-full border transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
        {...props}
      />
      <span className="bg-fg-muted peer-checked:bg-bg pointer-events-none absolute top-1/2 left-0.5 size-4 -translate-y-1/2 rounded-full transition-transform peer-checked:translate-x-4" />
    </span>
  )
})
