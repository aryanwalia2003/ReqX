import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import {
  openEnvironment,
  openPersonas,
  pickFile,
  pickSaveFile,
  runCollection,
} from '@/features/collection-runner/api'
import { err, ok } from '@/lib/result'

import { CollectionRunnerPanel } from './CollectionRunnerPanel'

vi.mock('@/features/collection-runner/api', () => ({
  pickFile: vi.fn(),
  pickSaveFile: vi.fn(),
  openEnvironment: vi.fn(),
  openPersonas: vi.fn(),
  runCollection: vi.fn(),
}))

const mockedPickFile = vi.mocked(pickFile)
const mockedPickSaveFile = vi.mocked(pickSaveFile)
const mockedOpenEnv = vi.mocked(openEnvironment)
const mockedOpenPersonas = vi.mocked(openPersonas)
const mockedRun = vi.mocked(runCollection)

const demoCollection = {
  path: '/home/dev/demo.json',
  collection: {
    name: 'Demo',
    requests: [{ name: 'Get user', method: 'GET', url: 'https://api.example.com/user' }],
  },
}

const emptySummary = {
  totalRequests: 0,
  totalSuccess: 0,
  totalFailures: 0,
  successRate: 0,
  avgLatencyMs: 0,
  p95LatencyMs: 0,
  totalDurationMs: 0,
}

describe('CollectionRunnerPanel', () => {
  beforeEach(() => vi.clearAllMocks())

  it('shows an empty state when nothing is selected', () => {
    render(<CollectionRunnerPanel selected={null} />)
    expect(screen.getByText('No collection selected')).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'Run' })).not.toBeInTheDocument()
  })

  it('lists the selected collection and runs it', async () => {
    mockedRun.mockResolvedValue(
      ok({
        stats: [
          {
            name: 'Get user',
            totalRuns: 1,
            successes: 1,
            failures: 0,
            avgLatencyMs: 12,
            p95LatencyMs: 12,
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

    render(<CollectionRunnerPanel selected={demoCollection} />)
    expect(screen.getByText(/Demo — 1 request/)).toBeInTheDocument()
    expect(screen.getByText('Get user')).toBeInTheDocument()

    await userEvent.click(screen.getByRole('button', { name: 'Run' }))

    expect(mockedRun).toHaveBeenCalledWith({
      collection: demoCollection.collection,
      envVariables: undefined,
    })
    await waitFor(() => expect(screen.getByText('1 ok')).toBeInTheDocument())
    expect(screen.getByText('0 failed')).toBeInTheDocument()
  })

  it('shows an alert when the run fails', async () => {
    mockedRun.mockResolvedValue(
      err({ kind: 'external_service_error', message: 'connection refused' }),
    )

    render(<CollectionRunnerPanel selected={demoCollection} />)
    await userEvent.click(screen.getByRole('button', { name: 'Run' }))

    await waitFor(() => expect(screen.getByText('connection refused')).toBeInTheDocument())
  })

  it('loads an environment and passes its variables into Run', async () => {
    mockedOpenEnv.mockResolvedValue(ok({ name: 'dev', variables: { host: 'api.example.com' } }))
    mockedRun.mockResolvedValue(
      ok({
        stats: [],
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

    render(<CollectionRunnerPanel selected={demoCollection} />)
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

  it('browsing for an environment file fills the path and auto-loads it', async () => {
    mockedPickFile.mockResolvedValue(ok('/home/dev/dev.json'))
    mockedOpenEnv.mockResolvedValue(ok({ name: 'dev', variables: {} }))

    render(<CollectionRunnerPanel selected={demoCollection} />)
    await userEvent.click(screen.getByRole('button', { name: 'Browse for environment file' }))

    expect(mockedPickFile).toHaveBeenCalledWith('Open environment', 'json')
    await waitFor(() => expect(mockedOpenEnv).toHaveBeenCalledWith('/home/dev/dev.json'))
    expect(screen.getByLabelText('Environment file path')).toHaveValue('/home/dev/dev.json')
  })

  it('passes workers/iterations/duration/rps through to Run', async () => {
    mockedRun.mockResolvedValue(
      ok({
        stats: [],
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

    render(<CollectionRunnerPanel selected={demoCollection} />)
    await userEvent.type(screen.getByLabelText('Workers'), '10')
    await userEvent.type(screen.getByLabelText('Iterations'), '50')
    await userEvent.type(screen.getByLabelText('Duration in seconds'), '30')
    await userEvent.type(screen.getByLabelText('Requests per second'), '25')
    await userEvent.click(screen.getByRole('button', { name: 'Run' }))

    await waitFor(() =>
      expect(mockedRun).toHaveBeenCalledWith(
        expect.objectContaining({ workers: 10, iterations: 50, durationMs: 30000, rps: 25 }),
      ),
    )
  })

  it('leaves load-test fields undefined when left blank', async () => {
    mockedRun.mockResolvedValue(
      ok({
        stats: [],
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

    render(<CollectionRunnerPanel selected={demoCollection} />)
    await userEvent.click(screen.getByRole('button', { name: 'Run' }))

    await waitFor(() =>
      expect(mockedRun).toHaveBeenCalledWith(
        expect.objectContaining({
          workers: undefined,
          iterations: undefined,
          durationMs: undefined,
          rps: undefined,
        }),
      ),
    )
  })

  it('loads personas via browse and passes them through to Run', async () => {
    mockedPickFile.mockResolvedValue(ok('/home/dev/personas.csv'))
    mockedOpenPersonas.mockResolvedValue(ok([{ name: 'Alice' }, { name: 'Bob' }]))
    mockedRun.mockResolvedValue(
      ok({
        stats: [],
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

    render(<CollectionRunnerPanel selected={demoCollection} />)
    await userEvent.click(screen.getByRole('button', { name: 'Browse for personas file' }))

    expect(mockedPickFile).toHaveBeenCalledWith('Open personas', 'csv')
    await waitFor(() => expect(mockedOpenPersonas).toHaveBeenCalledWith('/home/dev/personas.csv'))
    expect(screen.getByLabelText('Personas file path')).toHaveValue('/home/dev/personas.csv')
    expect(await screen.findByText(/2 personas loaded/)).toBeInTheDocument()

    await userEvent.click(screen.getByRole('button', { name: 'Run' }))

    await waitFor(() =>
      expect(mockedRun).toHaveBeenCalledWith(
        expect.objectContaining({ personas: [{ name: 'Alice' }, { name: 'Bob' }] }),
      ),
    )
  })

  it('does nothing when the environment file picker is cancelled', async () => {
    mockedPickFile.mockResolvedValue(ok(''))

    render(<CollectionRunnerPanel selected={demoCollection} />)
    await userEvent.click(screen.getByRole('button', { name: 'Browse for environment file' }))

    await waitFor(() => expect(mockedPickFile).toHaveBeenCalled())
    expect(mockedOpenEnv).not.toHaveBeenCalled()
    expect(screen.getByLabelText('Environment file path')).toHaveValue('')
  })

  it('advanced options are hidden until toggled, then pass through to Run', async () => {
    mockedRun.mockResolvedValue(ok({ stats: [], summary: emptySummary }))

    render(<CollectionRunnerPanel selected={demoCollection} />)
    expect(screen.queryByLabelText('Stages')).not.toBeInTheDocument()

    await userEvent.click(screen.getByRole('button', { name: 'Advanced options' }))
    expect(screen.getByLabelText('Stages')).toBeInTheDocument()

    await userEvent.type(screen.getByLabelText('Stages'), '10s:5,30s:20,10s:0')
    await userEvent.click(screen.getByText('No cookies'))
    await userEvent.click(screen.getByText('Clear cookies per request'))
    await userEvent.click(screen.getByText('GraphQL error detection'))
    await userEvent.type(screen.getByLabelText('Inject position'), '1')
    await userEvent.type(screen.getByLabelText('Inject name'), 'Health check')
    await userEvent.type(screen.getByLabelText('Inject URL'), 'https://api.example.com/health')

    await userEvent.click(screen.getByRole('button', { name: 'Run' }))

    await waitFor(() =>
      expect(mockedRun).toHaveBeenCalledWith(
        expect.objectContaining({
          stages: '10s:5,30s:20,10s:0',
          noCookies: true,
          clearCookies: true,
          graphql: true,
          injectIndex: '1',
          injectName: 'Health check',
          injectMethod: 'GET',
          injectUrl: 'https://api.example.com/health',
        }),
      ),
    )
  })

  it('browsing for an export location fills the path and passes it into Run', async () => {
    mockedPickSaveFile.mockResolvedValue(ok('/home/dev/results.ndjson'))
    mockedRun.mockResolvedValue(ok({ stats: [], summary: emptySummary }))

    render(<CollectionRunnerPanel selected={demoCollection} />)
    await userEvent.click(screen.getByRole('button', { name: 'Advanced options' }))
    await userEvent.click(screen.getByRole('button', { name: 'Browse (save)…' }))

    expect(mockedPickSaveFile).toHaveBeenCalledWith('Export results', 'results.ndjson')
    await waitFor(() =>
      expect(screen.getByLabelText('Export path')).toHaveValue('/home/dev/results.ndjson'),
    )

    await userEvent.click(screen.getByRole('button', { name: 'Run' }))

    await waitFor(() =>
      expect(mockedRun).toHaveBeenCalledWith(
        expect.objectContaining({ exportPath: '/home/dev/results.ndjson' }),
      ),
    )
  })

  it('shows the execution graph grouped by level when the run has DAG nodes', async () => {
    mockedRun.mockResolvedValue(
      ok({
        stats: [],
        summary: emptySummary,
        dagNodes: [
          { name: 'First', status: 'success', duration_ms: 10, level_idx: 0, depends_on: [] },
          {
            name: 'Second',
            status: 'failed',
            duration_ms: 20,
            level_idx: 1,
            depends_on: ['First'],
          },
        ],
      }),
    )

    render(<CollectionRunnerPanel selected={demoCollection} />)
    await userEvent.click(screen.getByRole('button', { name: 'Run' }))

    expect(await screen.findByText('Execution graph')).toBeInTheDocument()
    expect(screen.getByText('Level 0')).toBeInTheDocument()
    expect(screen.getByText('Level 1')).toBeInTheDocument()
    expect(screen.getByText(/First · success · 10 ms/)).toBeInTheDocument()
    expect(screen.getByText(/Second · failed · 20 ms/)).toBeInTheDocument()
  })

  it('does not show the execution graph for a linear run', async () => {
    mockedRun.mockResolvedValue(ok({ stats: [], summary: emptySummary }))

    render(<CollectionRunnerPanel selected={demoCollection} />)
    await userEvent.click(screen.getByRole('button', { name: 'Run' }))

    await waitFor(() => expect(mockedRun).toHaveBeenCalled())
    expect(screen.queryByText('Execution graph')).not.toBeInTheDocument()
  })
})
