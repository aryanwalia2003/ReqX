import { Component } from 'react'
import type { ErrorInfo, ReactNode } from 'react'

import { AlertCircleIcon } from '@/components/ui/icons'
import { cn } from '@/lib/cn'
import { reportError } from '@/lib/reportError'

export interface ErrorBoundaryProps {
  children: ReactNode
  /** `reportError`/logs me is boundary ko pehchaanne ke liye. */
  context?: string
  /** Badalte hi boundary khud reset ho jata hai — e.g. active tab id. */
  resetKey?: unknown
  fallback?: (error: Error, reset: () => void) => ReactNode
  className?: string
}

interface ErrorBoundaryState {
  error: Error | null
  resetKey: unknown
}

/** Render crash ko poore app ko blank hone se rokta — fallback UI dikhata hai. */
export class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  state: ErrorBoundaryState = { error: null, resetKey: this.props.resetKey }

  static getDerivedStateFromError(error: Error): Partial<ErrorBoundaryState> {
    return { error }
  }

  static getDerivedStateFromProps(
    props: ErrorBoundaryProps,
    state: ErrorBoundaryState,
  ): Partial<ErrorBoundaryState> | null {
    if (props.resetKey !== state.resetKey) {
      return { error: null, resetKey: props.resetKey }
    }
    return null
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    reportError(error, this.props.context ?? 'render')
    if (import.meta.env.DEV) {
      console.error(info.componentStack)
    }
  }

  reset = (): void => this.setState({ error: null })

  render() {
    const { error } = this.state
    if (!error) return this.props.children

    if (this.props.fallback) return this.props.fallback(error, this.reset)

    return (
      <div
        role="alert"
        className={cn(
          'border-border-strong bg-surface flex flex-col items-center gap-3 rounded-lg border border-dashed p-8 text-center',
          this.props.className,
        )}
      >
        <AlertCircleIcon className="text-fg-muted size-6" />
        <div className="flex flex-col gap-1">
          <p className="text-fg text-sm font-medium">Something broke.</p>
          <p className="text-fg-muted max-w-sm text-sm">{error.message}</p>
        </div>
        <button
          type="button"
          onClick={this.reset}
          className="text-fg-muted hover:text-fg text-sm underline underline-offset-2"
        >
          Try again
        </button>
      </div>
    )
  }
}
