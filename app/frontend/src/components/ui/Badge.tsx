import type { HTMLAttributes } from 'react'

import { cn } from '@/lib/cn'

export type BadgeVariant = 'default' | 'outline' | 'subtle'

export interface BadgeProps extends HTMLAttributes<HTMLSpanElement> {
  variant?: BadgeVariant
}

const variants: Record<BadgeVariant, string> = {
  default: 'bg-fg text-bg',
  outline: 'border border-border text-fg',
  subtle: 'bg-surface-2 text-fg-muted',
}

export function Badge({ variant = 'subtle', className, ...props }: BadgeProps) {
  return (
    <span
      className={cn(
        'inline-flex items-center rounded-md px-2 py-0.5 text-xs font-medium',
        variants[variant],
        className,
      )}
      {...props}
    />
  )
}
