import { useState } from 'react'

import { Alert, Badge, Button, ChevronDownIcon } from '@/components/ui'
import { openCollection, pickFile } from '@/features/collection-runner/api'
import { useRecentCollections } from '@/features/collection-runner/hooks/useRecentCollections'
import type { Collection, CollectionRequest } from '@/features/collection-runner/types'
import type { SendRequestInput } from '@/features/send-request'
import { cn } from '@/lib/cn'
import { getErrorMessage } from '@/lib/errors'
import { reportError } from '@/lib/reportError'

export interface CollectionsSidebarProps {
  selectedPath: string | null
  onSelectCollection: (path: string, collection: Collection) => void
  onOpenRequest: (request: SendRequestInput) => void
}

const COLLAPSED_KEY = 'reqx.sidebarCollapsed'

function loadCollapsed(): boolean {
  try {
    return localStorage.getItem(COLLAPSED_KEY) === 'true'
  } catch {
    return false
  }
}

/** Postman-jaisa left sidebar: recent collections + expandable request tree. */
export function CollectionsSidebar({
  selectedPath,
  onSelectCollection,
  onOpenRequest,
}: CollectionsSidebarProps) {
  const { recents, upsert, remove } = useRecentCollections()
  const [expanded, setExpanded] = useState<Set<string>>(new Set())
  const [isAdding, setIsAdding] = useState(false)
  const [addError, setAddError] = useState<string | null>(null)
  const [isCollapsed, setIsCollapsed] = useState(loadCollapsed)

  function toggleCollapsed() {
    setIsCollapsed((prev) => {
      const next = !prev
      try {
        localStorage.setItem(COLLAPSED_KEY, String(next))
      } catch {
        // Private-mode/quota — collapse still works this session, just doesn't persist.
      }
      return next
    })
  }

  function toggleExpand(path: string) {
    setExpanded((prev) => {
      const next = new Set(prev)
      if (next.has(path)) next.delete(path)
      else next.add(path)
      return next
    })
  }

  function handleRowClick(path: string, collection: Collection) {
    toggleExpand(path)
    upsert(path, collection) // bump to front, refresh lastOpenedAt
    onSelectCollection(path, collection)
  }

  function handleOpenRequest(req: CollectionRequest) {
    onOpenRequest({
      method: req.method,
      url: req.url,
      headers: req.headers,
      body: req.body,
      auth: req.auth,
    })
  }

  async function handleAdd() {
    setAddError(null)
    setIsAdding(true)
    const picked = await pickFile('Open collection', 'json')
    if (!picked.ok) {
      setIsAdding(false)
      reportError(picked.error, 'browse collection file')
      return
    }
    if (!picked.value) {
      setIsAdding(false)
      return // user cancelled
    }
    const opened = await openCollection(picked.value)
    setIsAdding(false)
    if (!opened.ok) {
      setAddError(getErrorMessage(opened.error))
      return
    }
    upsert(picked.value, opened.value)
    setExpanded((prev) => new Set(prev).add(picked.value))
    onSelectCollection(picked.value, opened.value)
  }

  if (isCollapsed) {
    return (
      <aside className="border-border bg-surface flex h-full w-10 shrink-0 flex-col items-center border-r py-3">
        <Button variant="ghost" size="icon" aria-label="Expand sidebar" onClick={toggleCollapsed}>
          <ChevronDownIcon className="-rotate-90" />
        </Button>
      </aside>
    )
  }

  return (
    <aside className="border-border bg-surface flex h-full w-64 shrink-0 flex-col gap-2 overflow-y-auto border-r p-3">
      <div className="flex items-center justify-between">
        <span className="text-fg text-sm font-semibold">Collections</span>
        <div className="flex items-center gap-1">
          <Button variant="ghost" size="sm" isLoading={isAdding} onClick={() => void handleAdd()}>
            + Add
          </Button>
          <Button
            variant="ghost"
            size="icon"
            aria-label="Collapse sidebar"
            onClick={toggleCollapsed}
          >
            <ChevronDownIcon className="rotate-90" />
          </Button>
        </div>
      </div>

      {addError && (
        <Alert variant="danger" title="Could not open collection">
          {addError}
        </Alert>
      )}

      {recents.length === 0 && (
        <p className="text-fg-subtle text-xs">No collections yet — click Add to browse for one.</p>
      )}

      <div className="flex flex-col gap-0.5">
        {recents.map((recent) => {
          const isExpanded = expanded.has(recent.path)
          const isSelected = selectedPath === recent.path
          return (
            <div key={recent.path} className="flex flex-col">
              <div className="group flex items-center gap-0.5">
                <button
                  type="button"
                  onClick={() => handleRowClick(recent.path, recent.collection)}
                  title={recent.path}
                  className={cn(
                    'flex min-w-0 flex-1 items-center gap-1.5 rounded-md px-1.5 py-1 text-left text-sm',
                    isSelected
                      ? 'bg-surface-2 text-fg'
                      : 'text-fg-muted hover:bg-surface-2 hover:text-fg',
                  )}
                >
                  <span className="text-fg-subtle w-3 shrink-0 text-xs">
                    {isExpanded ? '▾' : '▸'}
                  </span>
                  <span className="min-w-0 flex-1 truncate">
                    {recent.collection.name || 'Untitled collection'}
                  </span>
                  <Badge variant="subtle" className="shrink-0">
                    {recent.collection.requests.length}
                  </Badge>
                </button>
                <Button
                  variant="ghost"
                  size="icon"
                  aria-label={`Remove ${recent.collection.name || 'collection'} from recents`}
                  className="shrink-0 opacity-0 group-hover:opacity-100"
                  onClick={() => remove(recent.path)}
                >
                  ×
                </Button>
              </div>
              {isExpanded && (
                <div className="border-border ml-4 flex flex-col gap-0.5 border-l pl-2">
                  {recent.collection.requests.length === 0 && (
                    <span className="text-fg-subtle px-1.5 py-1 text-xs">No requests</span>
                  )}
                  {recent.collection.requests.map((req, i) => (
                    <button
                      key={`${req.name}-${i}`}
                      type="button"
                      onClick={() => handleOpenRequest(req)}
                      className="hover:bg-surface-2 flex items-center gap-1.5 rounded-md px-1.5 py-1 text-left text-xs"
                    >
                      <Badge variant="outline" className="shrink-0 px-1 py-0 text-[10px]">
                        {req.method}
                      </Badge>
                      <span className="text-fg-muted truncate">{req.name}</span>
                    </button>
                  ))}
                </div>
              )}
            </div>
          )
        })}
      </div>
    </aside>
  )
}
