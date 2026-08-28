import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { getRunStats, listRuns } from '@/features/history/api'
import { err, ok } from '@/lib/result'

import { HistoryPanel } from './HistoryPanel'

vi.mock('@/features/history/api', () => ({
  listRuns: vi.fn(),
  getRunStats: vi.fn(),
}))

const mockedListRuns = vi.mocked(listRuns)
const mockedGetRunStats = vi.mocked(getRunStats)

describe('HistoryPanel', () => {
  beforeEach(() => vi.clearAllMocks())

  it('loads runs on mount and shows them', async () => {
    mockedListRuns.mockResolvedValue(
      ok([
        {
          id: '1',
          ts: '2026-08-28 10:00:00',
          collection: 'demo.json',
          total_reqs: 5,
          rps: 12.3,
          p95_ms: 80,
          error_pct: 20,
        },
      ]),
    )

    render(<HistoryPanel />)

    await waitFor(() => expect(mockedListRuns).toHaveBeenCalled())
    expect(await screen.findByText('demo.json')).toBeInTheDocument()
    expect(screen.getByText('5 reqs')).toBeInTheDocument()
    expect(screen.getByText('20% err')).toBeInTheDocument()
  })

  it('shows an empty state when there is no history', async () => {
    mockedListRuns.mockResolvedValue(ok([]))

    render(<HistoryPanel />)

    expect(await screen.findByText('No runs yet')).toBeInTheDocument()
  })

  it('shows an alert when loading history fails', async () => {
    mockedListRuns.mockResolvedValue(err({ kind: 'database_error', message: 'disk full' }))

    render(<HistoryPanel />)

    expect(await screen.findByText('disk full')).toBeInTheDocument()
  })

  it('clicking a run loads and shows its per-request breakdown', async () => {
    mockedListRuns.mockResolvedValue(
      ok([
        {
          id: 'run-1',
          ts: '2026-08-28 10:00:00',
          collection: 'demo.json',
          total_reqs: 1,
          rps: 1,
          p95_ms: 10,
          error_pct: 0,
        },
      ]),
    )
    mockedGetRunStats.mockResolvedValue(
      ok([{ name: 'Get user', successes: 1, failures: 0, p95_ms: 10, avg_ms: 8 }]),
    )

    render(<HistoryPanel />)
    const row = await screen.findByText('demo.json')
    await userEvent.click(row)

    expect(mockedGetRunStats).toHaveBeenCalledWith('run-1')
    expect(await screen.findByText('Get user')).toBeInTheDocument()
    expect(screen.getByText('1 ok')).toBeInTheDocument()
  })
})
