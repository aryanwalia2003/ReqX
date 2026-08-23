import type { LabelHTMLAttributes } from 'react'

import { cn } from '@/lib/cn'

export interface LabelProps extends LabelHTMLAttributes<HTMLLabelElement> {
  required?: boolean
}

export function Label({ required, className, children, ...props }: LabelProps) {
  return (
    <label className={cn('text-fg text-sm font-medium', className)} {...props}>
      {children}
      {required && <span className="text-fg-subtle ml-0.5">*</span>}
    </label>
  )
}
