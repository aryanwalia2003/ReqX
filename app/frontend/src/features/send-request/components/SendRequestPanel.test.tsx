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

  it('sends basic auth credentials', async () => {
    mockedSendRequest.mockResolvedValue(
      ok({
        statusCode: 200,
        status: '200 OK',
        headers: {},
        body: '',
        durationMs: 1,
        bytesSent: 0,
        bytesReceived: 0,
      }),
    )

    render(<SendRequestPanel />)
    await userEvent.type(screen.getByLabelText('URL'), 'https://api.example.com')
    await userEvent.selectOptions(screen.getByLabelText('Auth type'), 'basic')
    await userEvent.type(screen.getByLabelText('Basic auth username'), 'admin')
    await userEvent.type(screen.getByLabelText('Basic auth password'), 'pw')
    await userEvent.click(screen.getByRole('button', { name: 'Send' }))

    await waitFor(() =>
      expect(mockedSendRequest).toHaveBeenCalledWith(
        expect.objectContaining({
          auth: { type: 'basic', username: 'admin', password: 'pw' },
        }),
      ),
    )
  })

  it('sends an API key in the configured location', async () => {
    mockedSendRequest.mockResolvedValue(
      ok({
        statusCode: 200,
        status: '200 OK',
        headers: {},
        body: '',
        durationMs: 1,
        bytesSent: 0,
        bytesReceived: 0,
      }),
    )

    render(<SendRequestPanel />)
    await userEvent.type(screen.getByLabelText('URL'), 'https://api.example.com')
    await userEvent.selectOptions(screen.getByLabelText('Auth type'), 'apikey')
    await userEvent.type(screen.getByLabelText('API key name'), 'X-Api-Key')
    await userEvent.type(screen.getByLabelText('API key value'), 'secret')
    await userEvent.selectOptions(screen.getByLabelText('API key location'), 'query')
    await userEvent.click(screen.getByRole('button', { name: 'Send' }))

    await waitFor(() =>
      expect(mockedSendRequest).toHaveBeenCalledWith(
        expect.objectContaining({
          auth: { type: 'apikey', key: 'X-Api-Key', value: 'secret', in: 'query' },
        }),
      ),
    )
  })

  it('sends cookie auth rows', async () => {
    mockedSendRequest.mockResolvedValue(
      ok({
        statusCode: 200,
        status: '200 OK',
        headers: {},
        body: '',
        durationMs: 1,
        bytesSent: 0,
        bytesReceived: 0,
      }),
    )

    render(<SendRequestPanel />)
    await userEvent.type(screen.getByLabelText('URL'), 'https://api.example.com')
    await userEvent.selectOptions(screen.getByLabelText('Auth type'), 'cookie')
    await userEvent.type(screen.getByLabelText('Cookie name'), 'session')
    await userEvent.type(screen.getByLabelText('Cookie value'), 'abc123')
    await userEvent.click(screen.getByRole('button', { name: 'Send' }))

    await waitFor(() =>
      expect(mockedSendRequest).toHaveBeenCalledWith(
        expect.objectContaining({
          auth: { type: 'cookie', cookies: { session: 'abc123' } },
        }),
      ),
    )
  })

  it('initializes the form from loadRequest (e.g. a request opened from the sidebar)', () => {
    render(
      <SendRequestPanel
        loadRequest={{
          method: 'POST',
          url: 'https://api.example.com/posts',
          headers: { 'X-Test': 'value' },
          body: '{"title":"hi"}',
          auth: { type: 'bearer', token: 'secret' },
        }}
      />,
    )

    expect(screen.getByLabelText('URL')).toHaveValue('https://api.example.com/posts')
    expect(screen.getByLabelText('HTTP method')).toHaveValue('POST')
    expect(screen.getByLabelText('Bearer token')).toHaveValue('secret')
    expect(screen.getByLabelText('Header name')).toHaveValue('X-Test')
    expect(screen.getByLabelText('Header value')).toHaveValue('value')
    expect(screen.getByLabelText('Body')).toHaveValue('{"title":"hi"}')
  })
})
