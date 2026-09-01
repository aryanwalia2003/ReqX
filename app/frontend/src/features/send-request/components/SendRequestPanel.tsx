import { useRef, useState } from 'react'

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
import type { AuthFieldsEditorHandle } from '@/features/send-request/components/AuthFieldsEditor'
import { AuthFieldsEditor } from '@/features/send-request/components/AuthFieldsEditor'
import { useSendRequest } from '@/features/send-request/hooks/useSendRequest'
import type { SendRequestInput } from '@/features/send-request/types'
import { getErrorMessage } from '@/lib/errors'
import { emptyRow, rowSetters, rowsFrom, rowsToRecord } from '@/lib/keyValueRows'
import type { KeyValueRow } from '@/lib/keyValueRows'

const METHODS = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'OPTIONS']

/** JSON Content-Type ho to pretty-print, warna raw string laut deta. */
function formatBody(body: string, contentType: string | undefined): string {
  if (!contentType?.includes('json')) return body
  try {
    return JSON.stringify(JSON.parse(body), null, 2)
  } catch {
    return body
  }
}

export interface SendRequestPanelProps {
  /**
   * Sidebar se ek request yahan bhej do — form usi se initialize hota. Sirf
   * MOUNT pe padha jaata (initial state), baad me prop badalne se form
   * apne aap reset nahi hota — caller ko naya `key` dena chahiye (e.g. ek
   * counter) taaki React fresh mount kare aur naya request load ho jaaye.
   */
  loadRequest?: SendRequestInput | null
}

/** Postman-lite: ek request banao, bhejo, response dekho — end-to-end Wails wiring. */
export function SendRequestPanel({ loadRequest }: SendRequestPanelProps = {}) {
  const [method, setMethod] = useState(loadRequest?.method || 'GET')
  const [url, setUrl] = useState(loadRequest?.url ?? '')
  const [headerRows, setHeaderRows] = useState<KeyValueRow[]>(() => rowsFrom(loadRequest?.headers))
  const [body, setBody] = useState(loadRequest?.body ?? '')
  const authRef = useRef<AuthFieldsEditorHandle>(null)

  const { data, error, isLoading, run } = useSendRequest()

  const headerHandlers = rowSetters(setHeaderRows)

  async function handleSend() {
    const input: SendRequestInput = {
      method,
      url,
      headers: rowsToRecord(headerRows),
      body: body || undefined,
      auth: authRef.current?.getValue(),
    }

    const result = await run(input)
    if (result.ok && headerRows.at(-1)?.key.trim()) {
      setHeaderRows((rows) => [...rows, emptyRow()])
    }
  }

  const formattedResponseBody = data ? formatBody(data.body, data.headers['Content-Type']) : ''

  return (
    <div className="flex flex-col gap-4">
      <Card>
        <CardHeader>
          <CardTitle>Send a request</CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col gap-4">
          <div className="flex gap-2">
            <Select
              value={method}
              onChange={(e) => setMethod(e.target.value)}
              className="w-32 shrink-0"
              aria-label="HTTP method"
            >
              {METHODS.map((m) => (
                <option key={m} value={m}>
                  {m}
                </option>
              ))}
            </Select>
            <Input
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              placeholder="https://api.example.com/users"
              className="flex-1"
              aria-label="URL"
            />
            <Button
              variant="primary"
              isLoading={isLoading}
              disabled={!url.trim()}
              onClick={() => void handleSend()}
            >
              Send
            </Button>
          </div>

          <AuthFieldsEditor ref={authRef} initial={loadRequest?.auth} />

          <div className="flex flex-col gap-1.5">
            <span className="text-fg text-sm font-medium">Headers</span>
            {headerRows.map((row) => (
              <div key={row.id} className="flex gap-2">
                <Input
                  value={row.key}
                  onChange={(e) => headerHandlers.update(row.id, { key: e.target.value })}
                  placeholder="Header-Name"
                  className="flex-1"
                  aria-label="Header name"
                />
                <Input
                  value={row.value}
                  onChange={(e) => headerHandlers.update(row.id, { value: e.target.value })}
                  placeholder="value"
                  className="flex-1"
                  aria-label="Header value"
                />
                <Button
                  variant="ghost"
                  size="icon"
                  aria-label="Remove header"
                  onClick={() => headerHandlers.remove(row.id)}
                >
                  ×
                </Button>
              </div>
            ))}
          </div>

          <Field label="Body">
            <Textarea
              mono
              rows={6}
              value={body}
              onChange={(e) => setBody(e.target.value)}
              placeholder="{}"
            />
          </Field>
        </CardContent>
      </Card>

      {error && (
        <Alert variant="danger" title={`Request failed (${error.kind})`}>
          {getErrorMessage(error)}
        </Alert>
      )}

      {data && (
        <Card>
          <CardHeader>
            <CardTitle className="flex flex-wrap items-center gap-2">
              <Badge>{data.statusCode}</Badge>
              <span>{data.status}</span>
              <Badge variant="outline">{data.durationMs} ms</Badge>
              <Badge variant="subtle">{data.bytesReceived} B</Badge>
            </CardTitle>
          </CardHeader>
          <CardContent className="flex flex-col gap-3">
            <div className="flex flex-col gap-1">
              <span className="text-fg-muted text-xs font-medium">Headers</span>
              <div className="border-border bg-surface flex flex-col gap-0.5 rounded-md border p-2 font-mono text-xs">
                {Object.entries(data.headers).map(([key, value]) => (
                  <div key={key}>
                    <span className="text-fg-muted">{key}:</span> {value}
                  </div>
                ))}
              </div>
            </div>
            <div className="flex flex-col gap-1">
              <span className="text-fg-muted text-xs font-medium">Body</span>
              <pre className="border-border bg-surface max-h-96 overflow-auto rounded-md border p-2 font-mono text-xs whitespace-pre-wrap">
                {formattedResponseBody || <span className="text-fg-subtle">(empty)</span>}
              </pre>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
