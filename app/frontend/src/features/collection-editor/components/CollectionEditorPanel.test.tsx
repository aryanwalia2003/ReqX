import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'

import { pickSaveFile, saveCollection } from '@/features/collection-runner/api'
import type { Collection } from '@/features/collection-runner/types'
import { err, ok } from '@/lib/result'

import { CollectionEditorPanel } from './CollectionEditorPanel'

vi.mock('@/features/collection-runner/api', () => ({
  pickSaveFile: vi.fn(),
  saveCollection: vi.fn(),
}))

const mockedPickSaveFile = vi.mocked(pickSaveFile)
const mockedSaveCollection = vi.mocked(saveCollection)

const EXISTING: Collection = {
  name: 'Demo',
  requests: [
    { name: 'Get user', method: 'GET', url: 'https://api.example.com/user' },
    { name: 'List users', method: 'GET', url: 'https://api.example.com/users' },
  ],
}

describe('CollectionEditorPanel', () => {
  it('starts empty with a default name when nothing is selected', () => {
    render(<CollectionEditorPanel selected={null} />)

    expect(screen.getByLabelText('Collection name')).toHaveValue('New collection')
    expect(screen.getByText('No requests yet.')).toBeInTheDocument()
  })

  it('loads the selected collection’s name and requests', () => {
    render(<CollectionEditorPanel selected={{ path: '/demo.json', collection: EXISTING }} />)

    expect(screen.getByLabelText('Collection name')).toHaveValue('Demo')
    expect(screen.getByText('Get user')).toBeInTheDocument()
    expect(screen.getByText('List users')).toBeInTheDocument()
  })

  it('adds a new request through the form', async () => {
    render(<CollectionEditorPanel selected={null} />)

    await userEvent.click(screen.getByRole('button', { name: '+ Add request' }))
    await userEvent.type(screen.getByLabelText('Request name'), 'Ping')
    await userEvent.type(screen.getByLabelText('Request URL'), 'https://api.example.com/ping')
    await userEvent.click(screen.getByRole('button', { name: 'Save request' }))

    expect(screen.getByText('Ping')).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'Save request' })).not.toBeInTheDocument()
  })

  it('edits an existing request', async () => {
    render(<CollectionEditorPanel selected={{ path: '/demo.json', collection: EXISTING }} />)

    const row = screen.getByText('Get user').closest('div')
    await userEvent.click(within(row as HTMLElement).getByRole('button', { name: 'Edit' }))

    const nameField = screen.getByLabelText('Request name')
    await userEvent.clear(nameField)
    await userEvent.type(nameField, 'Get user by id')
    await userEvent.click(screen.getByRole('button', { name: 'Save request' }))

    expect(screen.getByText('Get user by id')).toBeInTheDocument()
    expect(screen.queryByText('Get user', { exact: true })).not.toBeInTheDocument()
  })

  it('deletes a request', async () => {
    render(<CollectionEditorPanel selected={{ path: '/demo.json', collection: EXISTING }} />)

    const row = screen.getByText('Get user').closest('div')
    await userEvent.click(
      within(row as HTMLElement).getByRole('button', { name: 'Delete Get user' }),
    )

    expect(screen.queryByText('Get user', { exact: true })).not.toBeInTheDocument()
    expect(screen.getByText('List users')).toBeInTheDocument()
  })

  it('reorders requests with the move buttons', async () => {
    render(<CollectionEditorPanel selected={{ path: '/demo.json', collection: EXISTING }} />)

    await userEvent.click(screen.getByRole('button', { name: 'Move List users up' }))

    const names = screen.getAllByText(/Get user|List users/).map((el) => el.textContent)
    expect(names.indexOf('List users')).toBeLessThan(names.indexOf('Get user'))
  })

  it('saves directly when a path is already known', async () => {
    mockedSaveCollection.mockResolvedValue(ok(undefined))

    render(<CollectionEditorPanel selected={{ path: '/demo.json', collection: EXISTING }} />)
    await userEvent.click(screen.getByRole('button', { name: 'Save' }))

    expect(mockedPickSaveFile).not.toHaveBeenCalled()
    await waitFor(() =>
      expect(mockedSaveCollection).toHaveBeenCalledWith(
        { name: 'Demo', requests: EXISTING.requests },
        '/demo.json',
      ),
    )
    await waitFor(() => expect(screen.getByText('Saved to /demo.json')).toBeInTheDocument())
  })

  it('prompts for a path with the native save dialog for a brand new collection', async () => {
    mockedPickSaveFile.mockResolvedValue(ok('/home/dev/new-collection.json'))
    mockedSaveCollection.mockResolvedValue(ok(undefined))

    render(<CollectionEditorPanel selected={null} />)
    await userEvent.click(screen.getByRole('button', { name: 'Save…' }))

    await waitFor(() =>
      expect(mockedSaveCollection).toHaveBeenCalledWith(
        { name: 'New collection', requests: [] },
        '/home/dev/new-collection.json',
      ),
    )
    await waitFor(() =>
      expect(screen.getByText('Saved to /home/dev/new-collection.json')).toBeInTheDocument(),
    )
  })

  it('shows an error alert when saving fails', async () => {
    mockedSaveCollection.mockResolvedValue(err({ kind: 'internal_error', message: 'disk full' }))

    render(<CollectionEditorPanel selected={{ path: '/demo.json', collection: EXISTING }} />)
    await userEvent.click(screen.getByRole('button', { name: 'Save' }))

    await waitFor(() => expect(screen.getByText('disk full')).toBeInTheDocument())
  })
})
