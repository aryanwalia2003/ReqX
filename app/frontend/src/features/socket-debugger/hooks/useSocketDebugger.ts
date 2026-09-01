import { useEffect, useState } from 'react'

import {
  connectSocket,
  disconnectSocket,
  emitSocketEvent,
  sendSocketMessage,
} from '@/features/socket-debugger/api'
import type { ConnectSocketInput, SocketMessage } from '@/features/socket-debugger/types'
import { getErrorMessage } from '@/lib/errors'
import { EventsOn } from '@wails/runtime/runtime'

/** Live WS/Socket.IO debugger session — connection state + streamed transcript. */
export function useSocketDebugger() {
  const [connected, setConnected] = useState(false)
  const [connecting, setConnecting] = useState(false)
  const [error, setError] = useState<string | undefined>()
  const [messages, setMessages] = useState<SocketMessage[]>([])

  useEffect(() => {
    return EventsOn('socket:message', (...data) => {
      const msg = data[0] as SocketMessage
      setMessages((prev) => [...prev, msg])
      if (msg.direction === 'system' && msg.data.startsWith('Disconnected')) {
        setConnected(false)
      }
    })
  }, [])

  async function connect(input: ConnectSocketInput) {
    setConnecting(true)
    setError(undefined)
    setMessages([])
    const result = await connectSocket(input)
    setConnecting(false)
    if (!result.ok) {
      setError(getErrorMessage(result.error))
      return
    }
    setConnected(true)
  }

  async function send(text: string) {
    const result = await sendSocketMessage(text)
    if (!result.ok) setError(getErrorMessage(result.error))
  }

  async function emit(eventName: string, payload: string) {
    const result = await emitSocketEvent(eventName, payload)
    if (!result.ok) setError(getErrorMessage(result.error))
  }

  async function disconnect() {
    await disconnectSocket()
    setConnected(false)
  }

  return { connected, connecting, error, messages, connect, send, emit, disconnect }
}
