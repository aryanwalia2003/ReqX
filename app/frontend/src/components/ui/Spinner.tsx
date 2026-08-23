import { cn } from '@/lib/cn'

export interface SpinnerProps {
  size?: 'sm' | 'md' | 'lg'
  className?: string
}

const sizes = {
  sm: 'size-3.5',
  md: 'size-5',
  lg: 'size-7',
} as const

/** Chhota loading indicator — buttons, cards, page states sab me reusable. */
export function Spinner({ size = 'md', className }: SpinnerProps) {
  return (
    <svg
      viewBox="0 0 24 24"
      fill="none"
      role="status"
      aria-label="Loading"
      className={cn('animate-spin', sizes[size], className)}
    >
      <circle cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="3" opacity="0.25" />
      <path
        d="M22 12a10 10 0 0 1-10 10"
        stroke="currentColor"
        strokeWidth="3"
        strokeLinecap="round"
      />
    </svg>
  )
}
