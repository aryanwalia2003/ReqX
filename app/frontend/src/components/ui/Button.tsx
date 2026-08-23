import { forwardRef } from 'react'
import type { ButtonHTMLAttributes, ReactNode } from 'react'

import { Spinner } from '@/components/ui/Spinner'
import { cn } from '@/lib/cn'

export type ButtonVariant = 'primary' | 'secondary' | 'outline' | 'ghost' | 'destructive'
export type ButtonSize = 'sm' | 'md' | 'lg' | 'icon'

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant
  size?: ButtonSize
  isLoading?: boolean
  leftIcon?: ReactNode
  rightIcon?: ReactNode
}

const base =
  'inline-flex items-center justify-center gap-2 rounded-md font-medium transition-colors ' +
  'focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring ' +
  'disabled:pointer-events-none disabled:opacity-50'

// Har variant sirf grayscale weight se alag — koi hue kahin nahi.
const variants: Record<ButtonVariant, string> = {
  primary: 'bg-fg text-bg hover:bg-fg-muted',
  secondary: 'bg-surface-2 text-fg hover:bg-surface-3 border border-border',
  outline: 'bg-transparent text-fg border border-border hover:bg-surface-2',
  ghost: 'bg-transparent text-fg hover:bg-surface-2',
  // Dashed border = "sochke aage badho" — danger bina color ke bhi dikhta hai.
  destructive: 'bg-transparent text-fg border border-dashed border-fg-subtle hover:bg-surface-2',
}

const sizes: Record<ButtonSize, string> = {
  sm: 'h-8 px-3 text-sm',
  md: 'h-9 px-4 text-sm',
  lg: 'h-11 px-6 text-base',
  icon: 'size-9',
}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(function Button(
  {
    variant = 'secondary',
    size = 'md',
    isLoading,
    leftIcon,
    rightIcon,
    disabled,
    className,
    children,
    ...props
  },
  ref,
) {
  return (
    <button
      ref={ref}
      disabled={disabled || isLoading}
      className={cn(base, variants[variant], sizes[size], className)}
      {...props}
    >
      {isLoading ? <Spinner size="sm" /> : leftIcon}
      {children}
      {!isLoading && rightIcon}
    </button>
  )
})
