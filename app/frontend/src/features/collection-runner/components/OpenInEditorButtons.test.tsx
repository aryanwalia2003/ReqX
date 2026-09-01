import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'

import { openInEditor } from '@/features/collection-runner/api'
import { err, ok } from '@/lib/result'

import { OpenInEditorButtons } from './OpenInEditorButtons'

vi.mock('@/features/collection-runner/api', () => ({
  openInEditor: vi.fn(),
}))

const mockedOpenInEditor = vi.mocked(openInEditor)

describe('OpenInEditorButtons', () => {
  it('disables every button when no path is given', () => {
    render(<OpenInEditorButtons path={undefined} />)

    expect(screen.getByRole('button', { name: 'VS Code' })).toBeDisabled()
    expect(screen.getByRole('button', { name: 'Vim' })).toBeDisabled()
    expect(screen.getByRole('button', { name: 'Text editor' })).toBeDisabled()
  })

  it('opens the given path in VS Code', async () => {
    mockedOpenInEditor.mockResolvedValue(ok(undefined))

    render(<OpenInEditorButtons path="/tmp/collection.json" />)
    await userEvent.click(screen.getByRole('button', { name: 'VS Code' }))

    expect(mockedOpenInEditor).toHaveBeenCalledWith('vscode', '/tmp/collection.json')
  })

  it('opens the given path in Vim', async () => {
    mockedOpenInEditor.mockResolvedValue(ok(undefined))

    render(<OpenInEditorButtons path="/tmp/collection.json" />)
    await userEvent.click(screen.getByRole('button', { name: 'Vim' }))

    expect(mockedOpenInEditor).toHaveBeenCalledWith('vim', '/tmp/collection.json')
  })

  it('opens the given path in the system default text editor', async () => {
    mockedOpenInEditor.mockResolvedValue(ok(undefined))

    render(<OpenInEditorButtons path="/tmp/collection.json" />)
    await userEvent.click(screen.getByRole('button', { name: 'Text editor' }))

    expect(mockedOpenInEditor).toHaveBeenCalledWith('system', '/tmp/collection.json')
  })

  it('shows an error message when opening fails', async () => {
    mockedOpenInEditor.mockResolvedValue(
      err({ kind: 'invalid_input', message: 'VS Code CLI not found' }),
    )

    render(<OpenInEditorButtons path="/tmp/collection.json" />)
    await userEvent.click(screen.getByRole('button', { name: 'VS Code' }))

    await waitFor(() =>
      expect(
        screen.getByText(/Could not open in VS Code: VS Code CLI not found/),
      ).toBeInTheDocument(),
    )
  })
})
