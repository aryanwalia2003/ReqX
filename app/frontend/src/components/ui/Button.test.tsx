import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'

import { Button } from '@/components/ui/Button'

describe('Button', () => {
  it('renders children and handles clicks', async () => {
    const onClick = vi.fn()
    render(<Button onClick={onClick}>Send</Button>)
    await userEvent.click(screen.getByRole('button', { name: 'Send' }))
    expect(onClick).toHaveBeenCalledOnce()
  })

  it('disables the button and shows a spinner while loading', () => {
    render(<Button isLoading>Send</Button>)
    const button = screen.getByRole('button', { name: /send/i })
    expect(button).toBeDisabled()
    expect(screen.getByRole('status', { name: 'Loading' })).toBeInTheDocument()
  })

  it('does not fire onClick when disabled', async () => {
    const onClick = vi.fn()
    render(
      <Button disabled onClick={onClick}>
        Send
      </Button>,
    )
    await userEvent.click(screen.getByRole('button', { name: 'Send' }))
    expect(onClick).not.toHaveBeenCalled()
  })
})
