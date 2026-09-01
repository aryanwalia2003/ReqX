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
import { pickSaveFile, saveCollection } from '@/features/collection-runner/api'
import { OpenInEditorButtons } from '@/features/collection-runner/components/OpenInEditorButtons'
import type { Collection, CollectionRequest } from '@/features/collection-runner/types'
import type { AuthFieldsEditorHandle } from '@/features/send-request'
import { AuthFieldsEditor } from '@/features/send-request'
import { getErrorMessage } from '@/lib/errors'
import { rowSetters, rowsFrom, rowsToRecord } from '@/lib/keyValueRows'
import type { KeyValueRow } from '@/lib/keyValueRows'

const METHODS = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'OPTIONS']

export interface CollectionEditorPanelProps {
  /** Sidebar se chuni gayi collection — null matlab "naya banao" se shuru.
   * Sirf MOUNT pe padha jaata; caller `selected?.path ?? 'new'` ko key de
   * taaki collection badalne par form fresh remount ho (SendRequestPanel
   * jaisa hi pattern). */
  selected: { path: string; collection: Collection } | null
}

interface RequestFormProps {
  initial?: CollectionRequest
  onSave: (req: CollectionRequest) => void
  onCancel: () => void
}

/** Ek request add/edit karne ka form — SendRequestPanel jaisa hi shape,
 * bas "Send" ki jagah collection ke requests[] me save karta hai. */
function RequestForm({ initial, onSave, onCancel }: RequestFormProps) {
  const [name, setName] = useState(initial?.name ?? '')
  const [method, setMethod] = useState(initial?.method || 'GET')
  const [url, setUrl] = useState(initial?.url ?? '')
  const [headerRows, setHeaderRows] = useState<KeyValueRow[]>(() => rowsFrom(initial?.headers))
  const [body, setBody] = useState(initial?.body ?? '')
  const authRef = useRef<AuthFieldsEditorHandle>(null)
  const headerHandlers = rowSetters(setHeaderRows)

  function handleSave() {
    if (!name.trim() || !url.trim()) return
    onSave({
      name: name.trim(),
      method,
      url: url.trim(),
      headers: rowsToRecord(headerRows),
      body: body || undefined,
      auth: authRef.current?.getValue(),
    })
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>{initial ? 'Edit request' : 'Add request'}</CardTitle>
      </CardHeader>
      <CardContent className="flex flex-col gap-3">
        <div className="flex gap-2">
          <Select
            value={method}
            onChange={(e) => setMethod(e.target.value)}
            className="w-32 shrink-0"
            aria-label="Request method"
          >
            {METHODS.map((m) => (
              <option key={m} value={m}>
                {m}
              </option>
            ))}
          </Select>
          <Input
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Request name"
            className="flex-1"
            aria-label="Request name"
          />
        </div>
        <Input
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          placeholder="https://api.example.com/users"
          aria-label="Request URL"
        />

        <AuthFieldsEditor ref={authRef} initial={initial?.auth} />

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
            rows={5}
            value={body}
            onChange={(e) => setBody(e.target.value)}
            placeholder="{}"
          />
        </Field>

        <div className="flex justify-end gap-2">
          <Button variant="secondary" onClick={onCancel}>
            Cancel
          </Button>
          <Button variant="primary" disabled={!name.trim() || !url.trim()} onClick={handleSave}>
            Save request
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}

/** Postman-lite collection builder: naya banao ya khuli hui edit karo — add/
 * edit/reorder/delete requests, phir disk par save. */
export function CollectionEditorPanel({ selected }: CollectionEditorPanelProps) {
  const [name, setName] = useState(selected?.collection.name ?? 'New collection')
  const [requests, setRequests] = useState<CollectionRequest[]>(
    () => selected?.collection.requests ?? [],
  )
  const [path, setPath] = useState(selected?.path ?? '')
  const [editingIndex, setEditingIndex] = useState<number | null>(null)
  const [isSaving, setIsSaving] = useState(false)
  const [saveMessage, setSaveMessage] = useState<string | undefined>()
  const [saveError, setSaveError] = useState<string | undefined>()

  function handleSaveRequest(req: CollectionRequest) {
    setRequests((prev) =>
      editingIndex === -1 ? [...prev, req] : prev.map((r, i) => (i === editingIndex ? req : r)),
    )
    setEditingIndex(null)
  }

  function handleDelete(index: number) {
    setRequests((prev) => prev.filter((_, i) => i !== index))
  }

  function handleMove(index: number, dir: -1 | 1) {
    setRequests((prev) => {
      const target = index + dir
      if (target < 0 || target >= prev.length) return prev
      const next = [...prev]
      const moved = next[index]
      const swapped = next[target]
      if (!moved || !swapped) return prev
      next[index] = swapped
      next[target] = moved
      return next
    })
  }

  async function writeTo(savePath: string) {
    setIsSaving(true)
    setSaveError(undefined)
    setSaveMessage(undefined)
    const result = await saveCollection({ name: name.trim() || 'Untitled', requests }, savePath)
    setIsSaving(false)
    if (!result.ok) {
      setSaveError(getErrorMessage(result.error))
      return
    }
    setPath(savePath)
    setSaveMessage(`Saved to ${savePath}`)
  }

  async function handleSave() {
    if (path) {
      await writeTo(path)
      return
    }
    const picked = await pickSaveFile('Save collection', `${name.trim() || 'collection'}.json`)
    if (picked.ok && picked.value) await writeTo(picked.value)
  }

  async function handleSaveAs() {
    const picked = await pickSaveFile('Save collection as', `${name.trim() || 'collection'}.json`)
    if (picked.ok && picked.value) await writeTo(picked.value)
  }

  return (
    <div className="flex flex-col gap-4">
      <Card>
        <CardHeader>
          <CardTitle>Collection editor</CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col gap-3">
          <Field label="Collection name">
            <Input
              value={name}
              onChange={(e) => setName(e.target.value)}
              aria-label="Collection name"
            />
          </Field>

          {requests.length === 0 ? (
            <p className="text-fg-subtle text-sm">No requests yet.</p>
          ) : (
            <div className="flex flex-col gap-1.5">
              {requests.map((req, i) => (
                <div
                  key={i}
                  className="border-border flex items-center gap-2 rounded-md border p-2 text-sm"
                >
                  <Badge variant="outline">{req.method}</Badge>
                  <span className="flex-1 truncate">{req.name}</span>
                  <span className="text-fg-subtle truncate text-xs">{req.url}</span>
                  <Button
                    variant="ghost"
                    size="icon"
                    aria-label={`Move ${req.name} up`}
                    disabled={i === 0}
                    onClick={() => handleMove(i, -1)}
                  >
                    ↑
                  </Button>
                  <Button
                    variant="ghost"
                    size="icon"
                    aria-label={`Move ${req.name} down`}
                    disabled={i === requests.length - 1}
                    onClick={() => handleMove(i, 1)}
                  >
                    ↓
                  </Button>
                  <Button variant="ghost" size="sm" onClick={() => setEditingIndex(i)}>
                    Edit
                  </Button>
                  <Button
                    variant="ghost"
                    size="icon"
                    aria-label={`Delete ${req.name}`}
                    onClick={() => handleDelete(i)}
                  >
                    ×
                  </Button>
                </div>
              ))}
            </div>
          )}

          <div className="flex gap-2">
            <Button variant="secondary" onClick={() => setEditingIndex(-1)}>
              + Add request
            </Button>
            <Button variant="primary" isLoading={isSaving} onClick={() => void handleSave()}>
              {path ? 'Save' : 'Save…'}
            </Button>
            {path && (
              <Button variant="ghost" onClick={() => void handleSaveAs()}>
                Save as…
              </Button>
            )}
            {path && <OpenInEditorButtons path={path} className="ml-auto" />}
          </div>
        </CardContent>
      </Card>

      {saveError && (
        <Alert variant="danger" title="Save failed">
          {saveError}
        </Alert>
      )}
      {saveMessage && !saveError && (
        <Alert variant="success" title="Saved">
          {saveMessage}
        </Alert>
      )}

      {editingIndex !== null && (
        <RequestForm
          key={editingIndex}
          initial={editingIndex === -1 ? undefined : requests[editingIndex]}
          onSave={handleSaveRequest}
          onCancel={() => setEditingIndex(null)}
        />
      )}
    </div>
  )
}
