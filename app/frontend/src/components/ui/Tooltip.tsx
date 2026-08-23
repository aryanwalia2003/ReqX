import type { ReactNode } from 'react'

import { cn } from '@/lib/cn'

export interface TooltipProps {
  content: ReactNode
  side?: 'top' | 'bottom' | 'left' | 'right'
  children: ReactNode
  className?: string
}

const sides: Record<NonNullable<TooltipProps['side']>, string> = {
  top: 'bottom-full left-1/2 mb-2 -translate-x-1/2',
  bottom: 'top-full left-1/2 mt-2 -translate-x-1/2',
  left: 'right-full top-1/2 mr-2 -translate-y-1/2',
  right: 'left-full top-1/2 ml-2 -translate-y-1/2',
}

/** Hover/focus pe chhota label — koi JS positioning library nahi chahiye. */
export function Tooltip({ content, side = 'top', children, className }: TooltipProps) {
  return (
    <span className={cn('group relative inline-flex', className)}>
      {children}
      <span
        role="tooltip"
        className={cn(
          'border-border bg-surface-2 text-fg pointer-events-none absolute z-50 rounded-md border px-2 py-1 text-xs whitespace-nowrap opacity-0 transition-opacity delay-150 group-focus-within:opacity-100 group-hover:opacity-100',
          sides[side],
        )}
      >
        {content}
      </span>
    </span>
  )
}
