import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'

import { sendRequest } from '@/features/send-request/api'
import { err, ok } from '@/lib/result'

import { SendRequestPanel } from './SendRequestPanel'

vi.mock('@/features/send-request/api', () => ({
  sendRequest: vi.fn(),
}))

const mockedSendRequest = vi.mocked(sendRequest)

describe('SendRequestPanel', () => {
  it('sends the typed URL/method and renders the response', async () => {
    mockedSendRequest.mockResolvedValue(
      ok({
        statusCode: 200,
        status: '200 OK',
        headers: { 'Content-Type': 'application/json' },
        body: '{"id":1}',
        durationMs: 42,
        bytesSent: 0,
        bytesReceived: 8,
      }),
    )

    render(<SendRequestPanel />)
    await userEvent.type(screen.getByLabelText('URL'), 'https://api.example.com/users')
    await userEvent.click(screen.getByRole('button', { name: 'Send' }))

    expect(mockedSendRequest).toHaveBeenCalledWith(
      expect.objectContaining({ method: 'GET', url: 'https://api.example.com/users' }),
    )
    await waitFor(() => expect(screen.getByText('200 OK')).toBeInTheDocument())
    expect(screen.getByText('42 ms')).toBeInTheDocument()
  })

  it('shows an alert when the request fails', async () => {
    mockedSendRequest.mockResolvedValue(
      err({ kind: 'external_service_error', message: 'connection refused' }),
    )

    render(<SendRequestPanel />)
    await userEvent.type(screen.getByLabelText('URL'), 'http://localhost:1')
    await userEvent.click(screen.getByRole('button', { name: 'Send' }))

    await waitFor(() => expect(screen.getByText('connection refused')).toBeInTheDocument())
  })

  it('only sends non-empty header rows', async () => {
    mockedSendRequest.mockResolvedValue(
      ok({
        statusCode: 204,
        status: '204 No Content',
        headers: {},
        body: '',
        durationMs: 5,
        bytesSent: 0,
        bytesReceived: 0,
      }),
    )

    render(<SendRequestPanel />)
    await userEvent.type(screen.getByLabelText('URL'), 'https://api.example.com')
    await userEvent.type(screen.getByLabelText('Header name'), 'X-Test')
    await userEvent.type(screen.getByLabelText('Header value'), 'value')
    await userEvent.click(screen.getByRole('button', { name: 'Send' }))

    await waitFor(() =>
      expect(mockedSendRequest).toHaveBeenCalledWith(
        expect.objectContaining({ headers: { 'X-Test': 'value' } }),
      ),
    )
  })
})
