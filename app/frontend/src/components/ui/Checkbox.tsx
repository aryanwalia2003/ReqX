import { forwardRef } from 'react'
import type { InputHTMLAttributes } from 'react'

import { CheckIcon } from '@/components/ui/icons'
import { cn } from '@/lib/cn'

export type CheckboxProps = Omit<InputHTMLAttributes<HTMLInputElement>, 'type' | 'size'>

/** Native checkbox pe hi styling — checked state ka indicator peer se aata hai. */
export const Checkbox = forwardRef<HTMLInputElement, CheckboxProps>(function Checkbox(
  { className, ...props },
  ref,
) {
  return (
    <span className={cn('relative inline-flex size-4 shrink-0', className)}>
      <input
        ref={ref}
        type="checkbox"
        className="peer border-border bg-surface checked:border-fg checked:bg-fg focus-visible:outline-ring size-4 shrink-0 cursor-pointer appearance-none rounded border focus-visible:outline-2 focus-visible:outline-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
        {...props}
      />
      <CheckIcon className="text-bg pointer-events-none absolute inset-0 m-auto hidden size-3 peer-checked:block" />
    </span>
  )
})
