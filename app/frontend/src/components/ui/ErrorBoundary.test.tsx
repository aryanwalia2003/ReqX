import { render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'

import { ErrorBoundary } from '@/components/ui/ErrorBoundary'

function Bomb(): never {
  throw new Error('kaboom')
}

describe('ErrorBoundary', () => {
  it('renders children when nothing throws', () => {
    render(
      <ErrorBoundary>
        <p>all good</p>
      </ErrorBoundary>,
    )
    expect(screen.getByText('all good')).toBeInTheDocument()
  })

  it('renders a fallback instead of crashing the tree', () => {
    vi.spyOn(console, 'error').mockImplementation(() => {})

    render(
      <ErrorBoundary>
        <Bomb />
      </ErrorBoundary>,
    )

    expect(screen.getByText('kaboom')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Try again' })).toBeInTheDocument()
  })

  it('resets when resetKey changes', () => {
    vi.spyOn(console, 'error').mockImplementation(() => {})

    const { rerender } = render(
      <ErrorBoundary resetKey="a">
        <Bomb />
      </ErrorBoundary>,
    )
    expect(screen.getByText('kaboom')).toBeInTheDocument()

    rerender(
      <ErrorBoundary resetKey="b">
        <p>recovered</p>
      </ErrorBoundary>,
    )
    expect(screen.getByText('recovered')).toBeInTheDocument()
  })

  it('supports a custom fallback renderer with manual reset', async () => {
    vi.spyOn(console, 'error').mockImplementation(() => {})

    render(
      <ErrorBoundary fallback={(error, reset) => <button onClick={reset}>{error.message}</button>}>
        <Bomb />
      </ErrorBoundary>,
    )

    expect(screen.getByRole('button', { name: 'kaboom' })).toBeInTheDocument()
  })
})
