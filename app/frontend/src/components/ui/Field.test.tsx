import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'

import { Field } from '@/components/ui/Field'
import { Input } from '@/components/ui/Input'

describe('Field', () => {
  it('connects the label to the control via id', () => {
    render(
      <Field label="Request name">
        <Input />
      </Field>,
    )
    expect(screen.getByLabelText('Request name')).toBeInTheDocument()
  })

  it('marks the control invalid and describes it with the error text', () => {
    render(
      <Field label="Body" error="Valid JSON nahi hai.">
        <Input />
      </Field>,
    )
    const input = screen.getByLabelText('Body')
    expect(input).toHaveAttribute('aria-invalid', 'true')
    expect(input).toHaveAccessibleDescription('Valid JSON nahi hai.')
  })
})
