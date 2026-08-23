import type { ReactNode } from 'react'

import { cn } from '@/lib/cn'

export interface EmptyStateProps {
  icon?: ReactNode
  title: string
  description?: string
  action?: ReactNode
  className?: string
}

/** Khaali list/collection ke liye — "no requests yet" jaisi jagah. */
export function EmptyState({ icon, title, description, action, className }: EmptyStateProps) {
  return (
    <div
      className={cn(
        'border-border flex flex-col items-center gap-2 rounded-lg border border-dashed p-8 text-center',
        className,
      )}
    >
      {icon && <div className="text-fg-subtle">{icon}</div>}
      <p className="text-fg text-sm font-medium">{title}</p>
      {description && <p className="text-fg-muted max-w-xs text-sm">{description}</p>}
      {action && <div className="mt-2">{action}</div>}
    </div>
  )
}
