import type { HTMLAttributes } from 'react'

import { cn } from '@/lib/cn'

/** Keyboard shortcut dikhane ke liye — jaise Cmd+Enter ya Ctrl+K. */
export function Kbd({ className, ...props }: HTMLAttributes<HTMLElement>) {
  return (
    <kbd
      className={cn(
        'border-border bg-surface-2 text-fg-muted rounded border px-1.5 py-0.5 font-mono text-xs',
        className,
      )}
      {...props}
    />
  )
}
