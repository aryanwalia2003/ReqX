import { useState } from 'react'

import {
  Alert,
  Badge,
  Button,
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  Input,
} from '@/components/ui'
import { pickFile } from '@/features/collection-runner/api'
import { useOpenCollection } from '@/features/collection-runner/hooks/useOpenCollection'
import { useOpenEnvironment } from '@/features/collection-runner/hooks/useOpenEnvironment'
import { useRunCollection } from '@/features/collection-runner/hooks/useRunCollection'
import { getErrorMessage } from '@/lib/errors'
import { reportError } from '@/lib/reportError'

/** Postman-lite ka doosra half: ek poori collection (+ optional environment) load karo, run karo, results dekho. */
export function CollectionRunnerPanel() {
  const [path, setPath] = useState('')
  const [envPath, setEnvPath] = useState('')

  const {
    data: collection,
    error: openError,
    isLoading: isOpening,
    run: open,
  } = useOpenCollection()
  const {
    data: environment,
    error: envError,
    isLoading: isOpeningEnv,
    run: openEnv,
  } = useOpenEnvironment()
  const {
    data: output,
    error: runError,
    isLoading: isRunning,
    run: runCollection,
  } = useRunCollection()

  async function handleOpen() {
    await open(path)
  }

  async function handleOpenEnv() {
    await openEnv(envPath)
  }

  async function handleBrowseCollection() {
    const result = await pickFile('Open collection')
    if (!result.ok) return reportError(result.error, 'browse collection file')
    if (!result.value) return // user cancelled
    setPath(result.value)
    await open(result.value)
  }

  async function handleBrowseEnv() {
    const result = await pickFile('Open environment')
    if (!result.ok) return reportError(result.error, 'browse environment file')
    if (!result.value) return // user cancelled
    setEnvPath(result.value)
    await openEnv(result.value)
  }

  async function handleRun() {
    if (!collection) return
    await runCollection({ collection, envVariables: environment?.variables })
  }

  return (
    <div className="flex flex-col gap-4">
      <Card>
        <CardHeader>
          <CardTitle>Run a collection</CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col gap-4">
          <div className="flex gap-2">
            <Input
              value={path}
              onChange={(e) => setPath(e.target.value)}
              placeholder="collection.json"
              className="flex-1"
              aria-label="Collection file path"
            />
            <Button
              variant="outline"
              aria-label="Browse for collection file"
              onClick={() => void handleBrowseCollection()}
            >
              Browse…
            </Button>
            <Button
              variant="secondary"
              isLoading={isOpening}
              disabled={!path.trim()}
              onClick={() => void handleOpen()}
            >
              Open
            </Button>
            <Button
              variant="primary"
              isLoading={isRunning}
              disabled={!collection}
              onClick={() => void handleRun()}
            >
              Run
            </Button>
          </div>

          {openError && (
            <Alert variant="danger" title={`Could not open collection (${openError.kind})`}>
              {getErrorMessage(openError)}
            </Alert>
          )}

          <div className="flex gap-2">
            <Input
              value={envPath}
              onChange={(e) => setEnvPath(e.target.value)}
              placeholder="env.json (optional)"
              className="flex-1"
              aria-label="Environment file path"
            />
            <Button
              variant="outline"
              aria-label="Browse for environment file"
              onClick={() => void handleBrowseEnv()}
            >
              Browse…
            </Button>
            <Button
              variant="secondary"
              isLoading={isOpeningEnv}
              disabled={!envPath.trim()}
              onClick={() => void handleOpenEnv()}
            >
              Load env
            </Button>
          </div>

          {envError && (
            <Alert variant="danger" title={`Could not load environment (${envError.kind})`}>
              {getErrorMessage(envError)}
            </Alert>
          )}

          {environment && (
            <div className="flex flex-col gap-1.5">
              <span className="text-fg text-sm font-medium">
                {environment.name || 'Untitled environment'} —{' '}
                {Object.keys(environment.variables).length} variable
                {Object.keys(environment.variables).length === 1 ? '' : 's'}
              </span>
              <div className="border-border bg-surface flex flex-col gap-0.5 rounded-md border p-2 font-mono text-xs">
                {Object.entries(environment.variables).map(([key, value]) => (
                  <div key={key}>
                    <span className="text-fg-muted">{key}:</span> {value}
                  </div>
                ))}
              </div>
            </div>
          )}

          {collection && (
            <div className="flex flex-col gap-1.5">
              <span className="text-fg text-sm font-medium">
                {collection.name || 'Untitled collection'} — {collection.requests.length} request
                {collection.requests.length === 1 ? '' : 's'}
              </span>
              <div className="border-border bg-surface flex flex-col divide-y rounded-md border">
                {collection.requests.map((req, i) => (
                  <div
                    key={`${req.name}-${i}`}
                    className="flex items-center gap-2 px-3 py-1.5 text-sm"
                  >
                    <Badge variant="outline" className="w-16 shrink-0 justify-center">
                      {req.method}
                    </Badge>
                    <span className="text-fg-muted truncate">{req.name}</span>
                    <span className="text-fg-subtle ml-auto truncate text-xs">{req.url}</span>
                  </div>
                ))}
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      {runError && (
        <Alert variant="danger" title={`Run failed (${runError.kind})`}>
          {getErrorMessage(runError)}
        </Alert>
      )}

      {output && (
        <Card>
          <CardHeader>
            <CardTitle className="flex flex-wrap items-center gap-2">
              <Badge>{output.summary.totalSuccess} ok</Badge>
              <Badge variant="outline">{output.summary.totalFailures} failed</Badge>
              <Badge variant="subtle">{output.summary.avgLatencyMs} ms avg</Badge>
              <Badge variant="subtle">{output.summary.p95LatencyMs} ms p95</Badge>
              <Badge variant="subtle">{output.summary.totalDurationMs} ms total</Badge>
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="border-border bg-surface flex flex-col divide-y rounded-md border">
              {output.results.map((r, i) => (
                <div key={`${r.name}-${i}`} className="flex items-center gap-2 px-3 py-1.5 text-sm">
                  <Badge
                    variant={r.errorMessage ? 'outline' : 'default'}
                    className="w-12 shrink-0 justify-center"
                  >
                    {r.statusCode || '—'}
                  </Badge>
                  <span className="text-fg-muted truncate">{r.name}</span>
                  <span className="text-fg-subtle ml-auto shrink-0 text-xs">
                    {r.errorMessage ?? `${r.durationMs} ms`}
                  </span>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
