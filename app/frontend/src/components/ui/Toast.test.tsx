import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it } from 'vitest'

import { ToastProvider, useToast } from '@/components/ui/Toast'

function Example() {
  const toast = useToast()
  return <button onClick={() => toast({ title: 'Collection saved', duration: 300 })}>Save</button>
}

describe('ToastProvider / useToast', () => {
  it('shows a toast on demand and auto-dismisses it', async () => {
    render(
      <ToastProvider>
        <Example />
      </ToastProvider>,
    )

    await userEvent.click(screen.getByRole('button', { name: 'Save' }))
    expect(await screen.findByText('Collection saved')).toBeInTheDocument()

    await waitFor(() => expect(screen.queryByText('Collection saved')).not.toBeInTheDocument())
  })

  it('throws when used outside a provider', () => {
    function Broken() {
      useToast()
      return null
    }
    expect(() => render(<Broken />)).toThrow(/ToastProvider/)
  })
})
