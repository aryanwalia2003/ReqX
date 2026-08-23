import { forwardRef } from 'react'
import type { SelectHTMLAttributes } from 'react'

import { ChevronDownIcon } from '@/components/ui/icons'
import { cn } from '@/lib/cn'

export interface SelectProps extends SelectHTMLAttributes<HTMLSelectElement> {
  invalid?: boolean
}

/** Native <select> pe styling — koi extra JS dropdown nahi, hamesha accessible. */
export const Select = forwardRef<HTMLSelectElement, SelectProps>(function Select(
  { invalid, disabled, className, children, ...props },
  ref,
) {
  return (
    <div className={cn('relative', disabled && 'opacity-50')}>
      <select
        ref={ref}
        disabled={disabled}
        aria-invalid={invalid}
        className={cn(
          'bg-surface h-9 w-full appearance-none rounded-md border px-3 pr-8 text-sm',
          'text-fg outline-none',
          'border-border focus:border-border-strong focus:outline-ring focus:outline-2 focus:outline-offset-1',
          invalid && 'border-fg-subtle border-dashed',
          className,
        )}
        {...props}
      >
        {children}
      </select>
      <ChevronDownIcon className="text-fg-subtle pointer-events-none absolute top-1/2 right-2.5 -translate-y-1/2" />
    </div>
  )
})
