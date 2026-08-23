import type { HTMLAttributes } from 'react'

import { cn } from '@/lib/cn'

/** Loading placeholder block — width/height className se control karo. */
export function Skeleton({ className, ...props }: HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      aria-hidden="true"
      className={cn('bg-surface-2 animate-pulse rounded-md', className)}
      {...props}
    />
  )
}
