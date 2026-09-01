import { forwardRef, useImperativeHandle, useState } from 'react'

import { Button, Field, Input, Select } from '@/components/ui'
import type { AuthConfig } from '@/features/send-request/types'
import { rowSetters, rowsFrom, rowsToRecord } from '@/lib/keyValueRows'
import type { KeyValueRow } from '@/lib/keyValueRows'

export interface AuthFieldsEditorHandle {
  /** Current auth value — parent reads this only when it actually needs it
   * (submit/save), so this component never needs to re-render the parent
   * on every keystroke. */
  getValue(): AuthConfig | undefined
}

export interface AuthFieldsEditorProps {
  initial?: AuthConfig
}

/** Bearer/Basic/API key/Cookie auth editor — shared by SendRequestPanel and
 * the collection request editor so the four sub-forms exist in one place. */
export const AuthFieldsEditor = forwardRef<AuthFieldsEditorHandle, AuthFieldsEditorProps>(
  function AuthFieldsEditor({ initial }, ref) {
    const [authType, setAuthType] = useState<AuthConfig['type']>(initial?.type ?? 'none')
    const [bearerToken, setBearerToken] = useState(
      initial?.type === 'bearer' ? (initial.token ?? '') : '',
    )
    const [basicUsername, setBasicUsername] = useState(
      initial?.type === 'basic' ? (initial.username ?? '') : '',
    )
    const [basicPassword, setBasicPassword] = useState(
      initial?.type === 'basic' ? (initial.password ?? '') : '',
    )
    const [apiKeyName, setApiKeyName] = useState(
      initial?.type === 'apikey' ? (initial.key ?? '') : '',
    )
    const [apiKeyValue, setApiKeyValue] = useState(
      initial?.type === 'apikey' ? (initial.value ?? '') : '',
    )
    const [apiKeyIn, setApiKeyIn] = useState<'header' | 'query'>(
      initial?.type === 'apikey' ? (initial.in ?? 'header') : 'header',
    )
    const [cookieRows, setCookieRows] = useState<KeyValueRow[]>(() =>
      rowsFrom(initial?.type === 'cookie' ? initial.cookies : undefined),
    )
    const cookieHandlers = rowSetters(setCookieRows)

    useImperativeHandle(ref, () => ({
      getValue(): AuthConfig | undefined {
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
      },
    }))

    return (
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
    )
  },
)
