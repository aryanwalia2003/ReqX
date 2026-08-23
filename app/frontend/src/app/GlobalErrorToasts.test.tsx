import { act, render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'

import { GlobalErrorToasts } from '@/app/GlobalErrorToasts'
import { ToastProvider } from '@/components/ui'
import { reportError } from '@/lib/reportError'

describe('GlobalErrorToasts', () => {
  it('shows a danger toast when an error is reported', () => {
    vi.spyOn(console, 'error').mockImplementation(() => {})

    render(
      <ToastProvider>
        <GlobalErrorToasts />
      </ToastProvider>,
    )

    act(() => {
      reportError(new Error('boom'), 'test')
    })
    expect(screen.getByRole('status')).toHaveTextContent('boom')
  })
})
