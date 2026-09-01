import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'

import {
  connectSocket,
  disconnectSocket,
  emitSocketEvent,
  sendSocketMessage,
} from '@/features/socket-debugger/api'
import { err, ok } from '@/lib/result'
import { EventsOn } from '@wails/runtime/runtime'

import { SocketDebuggerPanel } from './SocketDebuggerPanel'

vi.mock('@/features/socket-debugger/api', () => ({
  connectSocket: vi.fn(),
  sendSocketMessage: vi.fn(),
  emitSocketEvent: vi.fn(),
  disconnectSocket: vi.fn(),
}))

vi.mock('@wails/runtime/runtime', () => ({
  EventsOn: vi.fn(() => vi.fn()),
}))

const mockedConnect = vi.mocked(connectSocket)
const mockedSend = vi.mocked(sendSocketMessage)
const mockedEmit = vi.mocked(emitSocketEvent)
const mockedDisconnect = vi.mocked(disconnectSocket)
const mockedEventsOn = vi.mocked(EventsOn)

describe('SocketDebuggerPanel', () => {
  it('connects with the typed URL and protocol, then shows the Send card', async () => {
    mockedConnect.mockResolvedValue(ok(undefined))

    render(<SocketDebuggerPanel />)
    await userEvent.type(screen.getByLabelText('Socket URL'), 'ws://localhost:8080')
    await userEvent.click(screen.getByRole('button', { name: 'Connect' }))

    expect(mockedConnect).toHaveBeenCalledWith({
      url: 'ws://localhost:8080',
      protocol: 'ws',
      headers: undefined,
    })
    await waitFor(() => expect(screen.getByText('Connected')).toBeInTheDocument())
    expect(screen.getByLabelText('Message to send')).toBeInTheDocument()
  })

  it('parses newline-separated headers into a record', async () => {
    mockedConnect.mockResolvedValue(ok(undefined))

    render(<SocketDebuggerPanel />)
    await userEvent.type(screen.getByLabelText('Socket URL'), 'ws://localhost:8080')
    await userEvent.type(
      screen.getByLabelText('Headers'),
      'Authorization: Bearer tok{enter}X-Test: value',
    )
    await userEvent.click(screen.getByRole('button', { name: 'Connect' }))

    await waitFor(() =>
      expect(mockedConnect).toHaveBeenCalledWith({
        url: 'ws://localhost:8080',
        protocol: 'ws',
        headers: { Authorization: 'Bearer tok', 'X-Test': 'value' },
      }),
    )
  })

  it('shows a connection error and stays disconnected on failure', async () => {
    mockedConnect.mockResolvedValue(err({ kind: 'external_service_error', message: 'refused' }))

    render(<SocketDebuggerPanel />)
    await userEvent.type(screen.getByLabelText('Socket URL'), 'ws://localhost:8080')
    await userEvent.click(screen.getByRole('button', { name: 'Connect' }))

    await waitFor(() => expect(screen.getByText('refused')).toBeInTheDocument())
    expect(screen.getByText('Disconnected')).toBeInTheDocument()
  })

  it('sends a raw message once connected over plain WebSocket', async () => {
    mockedConnect.mockResolvedValue(ok(undefined))
    mockedSend.mockResolvedValue(ok(undefined))

    render(<SocketDebuggerPanel />)
    await userEvent.type(screen.getByLabelText('Socket URL'), 'ws://localhost:8080')
    await userEvent.click(screen.getByRole('button', { name: 'Connect' }))
    await waitFor(() => expect(screen.getByText('Connected')).toBeInTheDocument())

    await userEvent.type(screen.getByLabelText('Message to send'), 'hello')
    await userEvent.click(screen.getByRole('button', { name: 'Send' }))

    await waitFor(() => expect(mockedSend).toHaveBeenCalledWith('hello'))
  })

  it('emits a Socket.IO event with name and payload once connected', async () => {
    mockedConnect.mockResolvedValue(ok(undefined))
    mockedEmit.mockResolvedValue(ok(undefined))

    render(<SocketDebuggerPanel />)
    await userEvent.selectOptions(screen.getByLabelText('Protocol'), 'sio')
    await userEvent.type(screen.getByLabelText('Socket URL'), 'http://localhost:3000')
    await userEvent.click(screen.getByRole('button', { name: 'Connect' }))
    await waitFor(() => expect(screen.getByText('Connected')).toBeInTheDocument())

    await userEvent.type(screen.getByLabelText('Event name'), 'greet')
    fireEvent.change(screen.getByLabelText('Event payload'), {
      target: { value: '{"name":"dev"}' },
    })
    await userEvent.click(screen.getByRole('button', { name: 'Emit' }))

    await waitFor(() => expect(mockedEmit).toHaveBeenCalledWith('greet', '{"name":"dev"}'))
  })

  it('disconnects and returns to the Disconnected state', async () => {
    mockedConnect.mockResolvedValue(ok(undefined))
    mockedDisconnect.mockResolvedValue(ok(undefined))

    render(<SocketDebuggerPanel />)
    await userEvent.type(screen.getByLabelText('Socket URL'), 'ws://localhost:8080')
    await userEvent.click(screen.getByRole('button', { name: 'Connect' }))
    await waitFor(() => expect(screen.getByText('Connected')).toBeInTheDocument())

    await userEvent.click(screen.getByRole('button', { name: 'Disconnect' }))

    await waitFor(() => expect(screen.getByText('Disconnected')).toBeInTheDocument())
    expect(mockedDisconnect).toHaveBeenCalled()
  })

  it('appends incoming socket:message events to the transcript', async () => {
    render(<SocketDebuggerPanel />)

    expect(mockedEventsOn).toHaveBeenCalledWith('socket:message', expect.any(Function))
    const callback = mockedEventsOn.mock.calls.at(-1)?.[1]

    callback?.({ direction: 'in', data: '{"ok":true}', eventName: 'greet', timestamp: 1000 })

    await waitFor(() => expect(screen.getByText('greet')).toBeInTheDocument())
    expect(screen.getByText(/"ok": true/)).toBeInTheDocument()
  })
})
