import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it } from 'vitest'

import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/Tabs'

function Example() {
  return (
    <Tabs defaultValue="requests">
      <TabsList>
        <TabsTrigger value="requests">Requests</TabsTrigger>
        <TabsTrigger value="env">Environments</TabsTrigger>
      </TabsList>
      <TabsContent value="requests">Requests panel</TabsContent>
      <TabsContent value="env">Env panel</TabsContent>
    </Tabs>
  )
}

describe('Tabs', () => {
  it('shows only the active panel and switches on click', async () => {
    render(<Example />)
    expect(screen.getByText('Requests panel')).toBeInTheDocument()
    expect(screen.queryByText('Env panel')).not.toBeInTheDocument()

    await userEvent.click(screen.getByRole('tab', { name: 'Environments' }))

    expect(screen.getByText('Env panel')).toBeInTheDocument()
    expect(screen.queryByText('Requests panel')).not.toBeInTheDocument()
    expect(screen.getByRole('tab', { name: 'Environments' })).toHaveAttribute(
      'aria-selected',
      'true',
    )
  })
})
