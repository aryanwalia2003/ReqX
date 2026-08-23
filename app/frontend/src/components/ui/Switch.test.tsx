import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it } from 'vitest'

import { Switch } from '@/components/ui/Switch'

describe('Switch', () => {
  it('toggles via the native switch role', async () => {
    render(<Switch aria-label="Notify on error" />)
    const switchEl = screen.getByRole('switch', { name: 'Notify on error' })
    expect(switchEl).not.toBeChecked()
    await userEvent.click(switchEl)
    expect(switchEl).toBeChecked()
  })
})
