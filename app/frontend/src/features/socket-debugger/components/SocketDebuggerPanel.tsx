import { useState } from 'react'

import {
  Alert,
  Badge,
  Button,
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  Field,
  Input,
  Select,
  Textarea,
} from '@/components/ui'
import { useSocketDebugger } from '@/features/socket-debugger/hooks/useSocketDebugger'
import type { SocketProtocol } from '@/features/socket-debugger/types'

function parseHeaders(text: string): Record<string, string> | undefined {
  const headers: Record<string, string> = {}
  for (const line of text.split('\n')) {
    const idx = line.indexOf(':')
    if (idx === -1) continue
    const key = line.slice(0, idx).trim()
    if (key) headers[key] = line.slice(idx + 1).trim()
  }
  return Object.keys(headers).length > 0 ? headers : undefined
}

/** JSON ho to pretty-print, warna raw string. */
function formatData(data: string): string {
  try {
    return JSON.stringify(JSON.parse(data), null, 2)
  } catch {
    return data
  }
}

const DIRECTION_LABEL: Record<string, string> = { in: '⬇ IN', out: '⬆ OUT', system: 'SYSTEM' }

/** Postman-lite WS/Socket.IO REPL — `reqx ws`/`reqx sio` ka GUI roop. */
export function SocketDebuggerPanel() {
  const [url, setUrl] = useState('')
  const [protocol, setProtocol] = useState<SocketProtocol>('ws')
  const [headersText, setHeadersText] = useState('')
  const [sendText, setSendText] = useState('')
  const [eventName, setEventName] = useState('')
  const [eventPayload, setEventPayload] = useState('')

  const { connected, connecting, error, messages, connect, send, emit, disconnect } =
    useSocketDebugger()

  async function handleConnect() {
    await connect({ url, protocol, headers: parseHeaders(headersText) })
  }

  async function handleSend() {
    if (!sendText.trim()) return
    await send(sendText)
    setSendText('')
  }

  async function handleEmit() {
    if (!eventName.trim()) return
    await emit(eventName.trim(), eventPayload)
    setEventPayload('')
  }

  return (
    <div className="flex flex-col gap-4">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            Socket debugger
            <Badge variant={connected ? 'default' : 'subtle'}>
              {connected ? 'Connected' : connecting ? 'Connecting…' : 'Disconnected'}
            </Badge>
          </CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col gap-3">
          <div className="flex gap-2">
            <Select
              value={protocol}
              onChange={(e) => setProtocol(e.target.value as SocketProtocol)}
              className="w-40 shrink-0"
              aria-label="Protocol"
              disabled={connected}
            >
              <option value="ws">WebSocket</option>
              <option value="sio">Socket.IO</option>
            </Select>
            <Input
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              placeholder={protocol === 'ws' ? 'ws://localhost:8080' : 'http://localhost:3000'}
              className="flex-1"
              aria-label="Socket URL"
              disabled={connected}
            />
            {connected ? (
              <Button variant="secondary" onClick={() => void disconnect()}>
                Disconnect
              </Button>
            ) : (
              <Button
                variant="primary"
                isLoading={connecting}
                disabled={!url.trim()}
                onClick={() => void handleConnect()}
              >
                Connect
              </Button>
            )}
          </div>

          <Field label="Headers" hint="Ek line ek header — Key: Value">
            <Textarea
              mono
              rows={2}
              value={headersText}
              onChange={(e) => setHeadersText(e.target.value)}
              placeholder="Authorization: Bearer token"
              disabled={connected}
            />
          </Field>
        </CardContent>
      </Card>

      {error && (
        <Alert variant="danger" title="Socket error">
          {error}
        </Alert>
      )}

      {connected && (
        <Card>
          <CardHeader>
            <CardTitle>Send</CardTitle>
          </CardHeader>
          <CardContent className="flex flex-col gap-3">
            {protocol === 'ws' ? (
              <div className="flex gap-2">
                <Input
                  value={sendText}
                  onChange={(e) => setSendText(e.target.value)}
                  placeholder="message"
                  className="flex-1"
                  aria-label="Message to send"
                  onKeyDown={(e) => e.key === 'Enter' && void handleSend()}
                />
                <Button onClick={() => void handleSend()} disabled={!sendText.trim()}>
                  Send
                </Button>
              </div>
            ) : (
              <div className="flex gap-2">
                <Input
                  value={eventName}
                  onChange={(e) => setEventName(e.target.value)}
                  placeholder="event name"
                  className="w-48 shrink-0"
                  aria-label="Event name"
                />
                <Input
                  value={eventPayload}
                  onChange={(e) => setEventPayload(e.target.value)}
                  placeholder='{"key":"value"}'
                  className="flex-1"
                  aria-label="Event payload"
                  onKeyDown={(e) => e.key === 'Enter' && void handleEmit()}
                />
                <Button onClick={() => void handleEmit()} disabled={!eventName.trim()}>
                  Emit
                </Button>
              </div>
            )}
          </CardContent>
        </Card>
      )}

      <Card>
        <CardHeader>
          <CardTitle>Transcript</CardTitle>
        </CardHeader>
        <CardContent>
          {messages.length === 0 ? (
            <p className="text-fg-subtle text-sm">No messages yet.</p>
          ) : (
            <div className="flex max-h-96 flex-col gap-1 overflow-y-auto font-mono text-xs">
              {messages.map((msg, i) => (
                <div
                  key={i}
                  className="border-border flex flex-col gap-0.5 border-b py-1 last:border-0"
                >
                  <div className="flex items-center gap-2">
                    <Badge variant={msg.direction === 'out' ? 'default' : 'outline'}>
                      {DIRECTION_LABEL[msg.direction]}
                    </Badge>
                    {msg.eventName && <span className="text-fg-muted">{msg.eventName}</span>}
                    <span className="text-fg-subtle">
                      {new Date(msg.timestamp).toLocaleTimeString()}
                    </span>
                  </div>
                  <pre className="whitespace-pre-wrap">{formatData(msg.data)}</pre>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
