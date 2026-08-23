import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it } from 'vitest'

import { Checkbox } from '@/components/ui/Checkbox'

describe('Checkbox', () => {
  it('toggles checked state on click', async () => {
    render(<Checkbox aria-label="Verify SSL" />)
    const checkbox = screen.getByRole('checkbox', { name: 'Verify SSL' })
    expect(checkbox).not.toBeChecked()
    await userEvent.click(checkbox)
    expect(checkbox).toBeChecked()
  })
})
