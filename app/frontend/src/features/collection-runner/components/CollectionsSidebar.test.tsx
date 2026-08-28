import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { openCollection, pickFile } from '@/features/collection-runner/api'
import { ok } from '@/lib/result'

import { CollectionsSidebar } from './CollectionsSidebar'

vi.mock('@/features/collection-runner/api', () => ({
  pickFile: vi.fn(),
  openCollection: vi.fn(),
}))

const mockedPickFile = vi.mocked(pickFile)
const mockedOpenCollection = vi.mocked(openCollection)

const demo = {
  name: 'Demo',
  requests: [
    {
      name: 'Get user',
      method: 'GET',
      url: 'https://api.example.com/user',
      headers: { Accept: 'json' },
    },
  ],
}

function renderSidebar(overrides: Partial<Parameters<typeof CollectionsSidebar>[0]> = {}) {
  const onSelectCollection = vi.fn()
  const onOpenRequest = vi.fn()
  const { unmount } = render(
    <CollectionsSidebar
      selectedPath={null}
      onSelectCollection={onSelectCollection}
      onOpenRequest={onOpenRequest}
      {...overrides}
    />,
  )
  return { onSelectCollection, onOpenRequest, unmount }
}

describe('CollectionsSidebar', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
  })

  it('shows an empty hint when there are no recents', () => {
    renderSidebar()
    expect(screen.getByText(/No collections yet/)).toBeInTheDocument()
  })

  it('adding a collection opens it, adds it to recents, and selects it', async () => {
    mockedPickFile.mockResolvedValue(ok('/home/dev/demo.json'))
    mockedOpenCollection.mockResolvedValue(ok(demo))
    const { onSelectCollection } = renderSidebar()

    await userEvent.click(screen.getByRole('button', { name: '+ Add' }))

    expect(mockedPickFile).toHaveBeenCalledWith('Open collection', 'json')
    await waitFor(() => expect(mockedOpenCollection).toHaveBeenCalledWith('/home/dev/demo.json'))
    expect(await screen.findByText('Demo')).toBeInTheDocument()
    expect(onSelectCollection).toHaveBeenCalledWith('/home/dev/demo.json', demo)
    // Added collections auto-expand — its request should already be visible.
    expect(screen.getByText('Get user')).toBeInTheDocument()
  })

  it('does nothing when the picker is cancelled', async () => {
    mockedPickFile.mockResolvedValue(ok(''))
    renderSidebar()

    await userEvent.click(screen.getByRole('button', { name: '+ Add' }))

    await waitFor(() => expect(mockedPickFile).toHaveBeenCalled())
    expect(mockedOpenCollection).not.toHaveBeenCalled()
    expect(screen.getByText(/No collections yet/)).toBeInTheDocument()
  })

  it('clicking a recent row toggles its request tree and selects it', async () => {
    mockedPickFile.mockResolvedValue(ok('/home/dev/demo.json'))
    mockedOpenCollection.mockResolvedValue(ok(demo))
    const { onSelectCollection } = renderSidebar()
    await userEvent.click(screen.getByRole('button', { name: '+ Add' }))
    await screen.findByText('Get user') // auto-expanded after add

    onSelectCollection.mockClear()
    await userEvent.click(screen.getByText('Demo'))
    // collapsed now
    expect(screen.queryByText('Get user')).not.toBeInTheDocument()
    expect(onSelectCollection).toHaveBeenCalledWith('/home/dev/demo.json', demo)

    await userEvent.click(screen.getByText('Demo'))
    expect(screen.getByText('Get user')).toBeInTheDocument()
  })

  it('clicking a request opens it with its headers', async () => {
    mockedPickFile.mockResolvedValue(ok('/home/dev/demo.json'))
    mockedOpenCollection.mockResolvedValue(ok(demo))
    const { onOpenRequest } = renderSidebar()
    await userEvent.click(screen.getByRole('button', { name: '+ Add' }))
    await userEvent.click(await screen.findByText('Get user'))

    expect(onOpenRequest).toHaveBeenCalledWith({
      method: 'GET',
      url: 'https://api.example.com/user',
      headers: { Accept: 'json' },
      body: undefined,
      auth: undefined,
    })
  })

  it('removing a collection drops it from the list', async () => {
    mockedPickFile.mockResolvedValue(ok('/home/dev/demo.json'))
    mockedOpenCollection.mockResolvedValue(ok(demo))
    renderSidebar()
    await userEvent.click(screen.getByRole('button', { name: '+ Add' }))
    await screen.findByText('Demo')

    await userEvent.click(screen.getByRole('button', { name: 'Remove Demo from recents' }))

    expect(screen.queryByText('Demo')).not.toBeInTheDocument()
    expect(screen.getByText(/No collections yet/)).toBeInTheDocument()
  })

  it('shows an alert when opening the picked file fails', async () => {
    mockedPickFile.mockResolvedValue(ok('/home/dev/broken.json'))
    mockedOpenCollection.mockResolvedValue({
      ok: false,
      error: { kind: 'invalid_input', message: 'bad json' },
    })
    renderSidebar()

    await userEvent.click(screen.getByRole('button', { name: '+ Add' }))

    expect(await screen.findByText('bad json')).toBeInTheDocument()
  })

  it('collapses and expands, hiding and restoring the collection list', async () => {
    renderSidebar()
    expect(screen.getByText('Collections')).toBeInTheDocument()

    await userEvent.click(screen.getByRole('button', { name: 'Collapse sidebar' }))

    expect(screen.queryByText('Collections')).not.toBeInTheDocument()
    expect(screen.queryByRole('button', { name: '+ Add' })).not.toBeInTheDocument()

    await userEvent.click(screen.getByRole('button', { name: 'Expand sidebar' }))
    expect(screen.getByText('Collections')).toBeInTheDocument()
  })

  it('remembers the collapsed state across remounts', async () => {
    const { unmount } = renderSidebar()
    await userEvent.click(screen.getByRole('button', { name: 'Collapse sidebar' }))
    unmount()

    renderSidebar()
    expect(screen.queryByText('Collections')).not.toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Expand sidebar' })).toBeInTheDocument()
  })
})
