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
import { useSendRequest } from '@/features/send-request/hooks/useSendRequest'
import type { AuthConfig, SendRequestInput } from '@/features/send-request/types'
import { getErrorMessage } from '@/lib/errors'

const METHODS = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'OPTIONS']

interface HeaderRow {
  id: number
  key: string
  value: string
}

let nextRowId = 1
function emptyRow(): HeaderRow {
  return { id: nextRowId++, key: '', value: '' }
}

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

function headerRowsFrom(headers: Record<string, string> | undefined): HeaderRow[] {
  const entries = Object.entries(headers ?? {})
  return entries.length > 0
    ? entries.map(([key, value]) => ({ id: nextRowId++, key, value }))
    : [emptyRow()]
}

/** Ek hi row-array (headers, cookies) ke liye reusable update/remove — is
 * component me do jagah same shape chahiye tha, isliye ek jagah nikaala. */
function rowSetters(setRows: (updater: (rows: HeaderRow[]) => HeaderRow[]) => void) {
  return {
    update(id: number, patch: Partial<HeaderRow>) {
      setRows((rows) => rows.map((row) => (row.id === id ? { ...row, ...patch } : row)))
    },
    remove(id: number) {
      setRows((rows) => (rows.length > 1 ? rows.filter((row) => row.id !== id) : rows))
    },
  }
}

function rowsToRecord(rows: HeaderRow[]): Record<string, string> | undefined {
  const record: Record<string, string> = {}
  for (const row of rows) {
    if (row.key.trim()) record[row.key.trim()] = row.value
  }
  return Object.keys(record).length > 0 ? record : undefined
}

/** Postman-lite: ek request banao, bhejo, response dekho — end-to-end Wails wiring. */
export function SendRequestPanel({ loadRequest }: SendRequestPanelProps = {}) {
  const [method, setMethod] = useState(loadRequest?.method || 'GET')
  const [url, setUrl] = useState(loadRequest?.url ?? '')
  const [headerRows, setHeaderRows] = useState<HeaderRow[]>(() =>
    headerRowsFrom(loadRequest?.headers),
  )
  const [body, setBody] = useState(loadRequest?.body ?? '')

  const initialAuth = loadRequest?.auth
  const [authType, setAuthType] = useState<AuthConfig['type']>(initialAuth?.type ?? 'none')
  const [bearerToken, setBearerToken] = useState(
    initialAuth?.type === 'bearer' ? (initialAuth.token ?? '') : '',
  )
  const [basicUsername, setBasicUsername] = useState(
    initialAuth?.type === 'basic' ? (initialAuth.username ?? '') : '',
  )
  const [basicPassword, setBasicPassword] = useState(
    initialAuth?.type === 'basic' ? (initialAuth.password ?? '') : '',
  )
  const [apiKeyName, setApiKeyName] = useState(
    initialAuth?.type === 'apikey' ? (initialAuth.key ?? '') : '',
  )
  const [apiKeyValue, setApiKeyValue] = useState(
    initialAuth?.type === 'apikey' ? (initialAuth.value ?? '') : '',
  )
  const [apiKeyIn, setApiKeyIn] = useState<'header' | 'query'>(
    initialAuth?.type === 'apikey' ? (initialAuth.in ?? 'header') : 'header',
  )
  const [cookieRows, setCookieRows] = useState<HeaderRow[]>(() =>
    headerRowsFrom(initialAuth?.type === 'cookie' ? initialAuth.cookies : undefined),
  )

  const { data, error, isLoading, run } = useSendRequest()

  const headerHandlers = rowSetters(setHeaderRows)
  const cookieHandlers = rowSetters(setCookieRows)

  function buildAuth(): AuthConfig | undefined {
    switch (authType) {
      case 'bearer':
        return bearerToken.trim() ? { type: 'bearer', token: bearerToken.trim() } : undefined
      case 'basic':
        return basicUsername.trim()
          ? { type: 'basic', username: basicUsername.trim(), password: basicPassword }
          : undefined
      case 'apikey':
        return apiKeyName.trim()
          ? { type: 'apikey', key: apiKeyName.trim(), value: apiKeyValue, in: apiKeyIn }
          : undefined
      case 'cookie': {
        const cookies = rowsToRecord(cookieRows)
        return cookies ? { type: 'cookie', cookies } : undefined
      }
      default:
        return undefined
    }
  }

  async function handleSend() {
    const input: SendRequestInput = {
      method,
      url,
      headers: rowsToRecord(headerRows),
      body: body || undefined,
      auth: buildAuth(),
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

          <div className="flex flex-col gap-2">
            <Field label="Auth type">
              <Select
                value={authType}
                onChange={(e) => setAuthType(e.target.value as AuthConfig['type'])}
                aria-label="Auth type"
              >
                <option value="none">None</option>
                <option value="bearer">Bearer token</option>
                <option value="basic">Basic auth</option>
                <option value="apikey">API key</option>
                <option value="cookie">Cookie</option>
              </Select>
            </Field>

            {authType === 'bearer' && (
              <Field label="Token">
                <Input
                  value={bearerToken}
                  onChange={(e) => setBearerToken(e.target.value)}
                  placeholder="token"
                  aria-label="Bearer token"
                />
              </Field>
            )}

            {authType === 'basic' && (
              <div className="flex gap-2">
                <Field label="Username" className="flex-1">
                  <Input
                    value={basicUsername}
                    onChange={(e) => setBasicUsername(e.target.value)}
                    aria-label="Basic auth username"
                  />
                </Field>
                <Field label="Password" className="flex-1">
                  <Input
                    type="password"
                    value={basicPassword}
                    onChange={(e) => setBasicPassword(e.target.value)}
                    aria-label="Basic auth password"
                  />
                </Field>
              </div>
            )}

            {authType === 'apikey' && (
              <div className="flex gap-2">
                <Field label="Key" className="flex-1">
                  <Input
                    value={apiKeyName}
                    onChange={(e) => setApiKeyName(e.target.value)}
                    aria-label="API key name"
                  />
                </Field>
                <Field label="Value" className="flex-1">
                  <Input
                    value={apiKeyValue}
                    onChange={(e) => setApiKeyValue(e.target.value)}
                    aria-label="API key value"
                  />
                </Field>
                <Field label="Add to">
                  <Select
                    value={apiKeyIn}
                    onChange={(e) => setApiKeyIn(e.target.value as 'header' | 'query')}
                    aria-label="API key location"
                  >
                    <option value="header">Header</option>
                    <option value="query">Query param</option>
                  </Select>
                </Field>
              </div>
            )}

            {authType === 'cookie' && (
              <div className="flex flex-col gap-1.5">
                <span className="text-fg text-sm font-medium">Cookies</span>
                {cookieRows.map((row) => (
                  <div key={row.id} className="flex gap-2">
                    <Input
                      value={row.key}
                      onChange={(e) => cookieHandlers.update(row.id, { key: e.target.value })}
                      placeholder="Cookie-Name"
                      className="flex-1"
                      aria-label="Cookie name"
                    />
                    <Input
                      value={row.value}
                      onChange={(e) => cookieHandlers.update(row.id, { value: e.target.value })}
                      placeholder="value"
                      className="flex-1"
                      aria-label="Cookie value"
                    />
                    <Button
                      variant="ghost"
                      size="icon"
                      aria-label="Remove cookie"
                      onClick={() => cookieHandlers.remove(row.id)}
                    >
                      ×
                    </Button>
                  </div>
                ))}
              </div>
            )}
          </div>

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
