import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { useState } from 'react'
import { describe, expect, it, vi } from 'vitest'

import { Dialog, DialogHeader, DialogTitle } from '@/components/ui/Dialog'

function Example() {
  const [open, setOpen] = useState(true)
  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTitle>Delete collection?</DialogTitle>
    </Dialog>
  )
}

describe('Dialog', () => {
  it('renders its content while open', () => {
    render(<Example />)
    expect(screen.getByText('Delete collection?')).toBeInTheDocument()
  })

  it('calls onOpenChange(false) via the built-in close button', async () => {
    const onOpenChange = vi.fn()
    render(
      <Dialog open onOpenChange={onOpenChange}>
        <DialogHeader>
          <DialogTitle>Hello</DialogTitle>
        </DialogHeader>
      </Dialog>,
    )
    await userEvent.click(screen.getByRole('button', { name: 'Close' }))
    expect(onOpenChange).toHaveBeenCalledWith(false)
  })
})
