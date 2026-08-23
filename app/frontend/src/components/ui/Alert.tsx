import type { HTMLAttributes, ReactNode } from 'react'

import { AlertCircleIcon, AlertTriangleIcon, CheckIcon, InfoIcon } from '@/components/ui/icons'
import { cn } from '@/lib/cn'

export type AlertVariant = 'info' | 'success' | 'warning' | 'danger'

export interface AlertProps extends Omit<HTMLAttributes<HTMLDivElement>, 'title'> {
  variant?: AlertVariant
  title?: ReactNode
}

const icons: Record<AlertVariant, typeof InfoIcon> = {
  info: InfoIcon,
  success: CheckIcon,
  warning: AlertTriangleIcon,
  danger: AlertCircleIcon,
}

// Rang nahi — border-style + icon se hi info/warning/danger alag dikhte hain.
const borders: Record<AlertVariant, string> = {
  info: 'border-border',
  success: 'border-border',
  warning: 'border-dashed border-border-strong',
  danger: 'border-dashed border-border-strong',
}

export function Alert({ variant = 'info', title, className, children, ...props }: AlertProps) {
  const Icon = icons[variant]
  return (
    <div
      role={variant === 'danger' ? 'alert' : 'status'}
      className={cn(
        'bg-surface flex gap-3 rounded-lg border p-3 text-sm',
        borders[variant],
        className,
      )}
      {...props}
    >
      <Icon className="text-fg-muted mt-0.5 size-4 shrink-0" />
      <div className="flex flex-col gap-0.5">
        {title && <p className="text-fg font-medium">{title}</p>}
        {children && <div className="text-fg-muted">{children}</div>}
      </div>
    </div>
  )
}
