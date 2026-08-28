import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import {
  openCollection,
  openEnvironment,
  pickFile,
  runCollection,
} from '@/features/collection-runner/api'
import { err, ok } from '@/lib/result'

import { CollectionRunnerPanel } from './CollectionRunnerPanel'

vi.mock('@/features/collection-runner/api', () => ({
  pickFile: vi.fn(),
  openCollection: vi.fn(),
  openEnvironment: vi.fn(),
  runCollection: vi.fn(),
}))

const mockedPickFile = vi.mocked(pickFile)
const mockedOpen = vi.mocked(openCollection)
const mockedOpenEnv = vi.mocked(openEnvironment)
const mockedRun = vi.mocked(runCollection)

describe('CollectionRunnerPanel', () => {
  beforeEach(() => vi.clearAllMocks())

  it('opens a collection, lists its requests, then runs it and shows results', async () => {
    mockedOpen.mockResolvedValue(
      ok({
        name: 'Demo',
        requests: [{ name: 'Get user', method: 'GET', url: 'https://api.example.com/user' }],
      }),
    )
    mockedRun.mockResolvedValue(
      ok({
        results: [
          {
            name: 'Get user',
            protocol: 'HTTP',
            statusCode: 200,
            statusString: '200 OK',
            durationMs: 12,
            bytesReceived: 34,
          },
        ],
        summary: {
          totalRequests: 1,
          totalSuccess: 1,
          totalFailures: 0,
          successRate: 100,
          avgLatencyMs: 12,
          p95LatencyMs: 12,
          totalDurationMs: 12,
        },
      }),
    )

    render(<CollectionRunnerPanel />)
    await userEvent.type(screen.getByLabelText('Collection file path'), 'collection.json')
    await userEvent.click(screen.getByRole('button', { name: 'Open' }))

    expect(mockedOpen).toHaveBeenCalledWith('collection.json')
    await waitFor(() => expect(screen.getByText(/Demo — 1 request/)).toBeInTheDocument())
    expect(screen.getByText('Get user')).toBeInTheDocument()

    await userEvent.click(screen.getByRole('button', { name: 'Run' }))

    expect(mockedRun).toHaveBeenCalledWith({
      collection: {
        name: 'Demo',
        requests: [{ name: 'Get user', method: 'GET', url: 'https://api.example.com/user' }],
      },
    })
    await waitFor(() => expect(screen.getByText('1 ok')).toBeInTheDocument())
    expect(screen.getByText('0 failed')).toBeInTheDocument()
  })

  it('loads an environment and passes its variables into Run', async () => {
    mockedOpen.mockResolvedValue(
      ok({
        name: 'Demo',
        requests: [{ name: 'Get user', method: 'GET', url: 'https://{{host}}/user' }],
      }),
    )
    mockedOpenEnv.mockResolvedValue(ok({ name: 'dev', variables: { host: 'api.example.com' } }))
    mockedRun.mockResolvedValue(
      ok({
        results: [],
        summary: {
          totalRequests: 0,
          totalSuccess: 0,
          totalFailures: 0,
          successRate: 0,
          avgLatencyMs: 0,
          p95LatencyMs: 0,
          totalDurationMs: 0,
        },
      }),
    )

    render(<CollectionRunnerPanel />)
    await userEvent.type(screen.getByLabelText('Collection file path'), 'collection.json')
    await userEvent.click(screen.getByRole('button', { name: 'Open' }))
    await waitFor(() => expect(screen.getByText(/Demo — 1 request/)).toBeInTheDocument())

    await userEvent.type(screen.getByLabelText('Environment file path'), 'dev.json')
    await userEvent.click(screen.getByRole('button', { name: 'Load env' }))

    expect(mockedOpenEnv).toHaveBeenCalledWith('dev.json')
    await waitFor(() => expect(screen.getByText(/dev — 1 variable/)).toBeInTheDocument())
    expect(screen.getByText('api.example.com')).toBeInTheDocument()

    await userEvent.click(screen.getByRole('button', { name: 'Run' }))

    await waitFor(() =>
      expect(mockedRun).toHaveBeenCalledWith(
        expect.objectContaining({ envVariables: { host: 'api.example.com' } }),
      ),
    )
  })

  it('browsing for a file fills the path and auto-opens it', async () => {
    mockedPickFile.mockResolvedValue(ok('/home/dev/collection.json'))
    mockedOpen.mockResolvedValue(
      ok({
        name: 'Demo',
        requests: [{ name: 'Get user', method: 'GET', url: 'https://api.example.com/user' }],
      }),
    )

    render(<CollectionRunnerPanel />)
    await userEvent.click(screen.getByRole('button', { name: 'Browse for collection file' }))

    expect(mockedPickFile).toHaveBeenCalledWith('Open collection')
    await waitFor(() => expect(mockedOpen).toHaveBeenCalledWith('/home/dev/collection.json'))
    expect(screen.getByLabelText('Collection file path')).toHaveValue('/home/dev/collection.json')
    await waitFor(() => expect(screen.getByText(/Demo — 1 request/)).toBeInTheDocument())
  })

  it('does nothing when the file picker is cancelled', async () => {
    mockedPickFile.mockResolvedValue(ok(''))

    render(<CollectionRunnerPanel />)
    await userEvent.click(screen.getByRole('button', { name: 'Browse for collection file' }))

    await waitFor(() => expect(mockedPickFile).toHaveBeenCalled())
    expect(mockedOpen).not.toHaveBeenCalled()
    expect(screen.getByLabelText('Collection file path')).toHaveValue('')
  })

  it('shows an alert when opening fails', async () => {
    mockedOpen.mockResolvedValue(err({ kind: 'invalid_input', message: 'no such file' }))

    render(<CollectionRunnerPanel />)
    await userEvent.type(screen.getByLabelText('Collection file path'), 'missing.json')
    await userEvent.click(screen.getByRole('button', { name: 'Open' }))

    await waitFor(() => expect(screen.getByText('no such file')).toBeInTheDocument())
  })

  it('disables Run until a collection is loaded', () => {
    render(<CollectionRunnerPanel />)
    expect(screen.getByRole('button', { name: 'Run' })).toBeDisabled()
  })
})
