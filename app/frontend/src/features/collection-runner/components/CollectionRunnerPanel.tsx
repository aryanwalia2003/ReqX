import { useState } from 'react'

import {
  Alert,
  Badge,
  Button,
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  EmptyState,
  Field,
  Input,
} from '@/components/ui'
import { pickFile } from '@/features/collection-runner/api'
import { useOpenEnvironment } from '@/features/collection-runner/hooks/useOpenEnvironment'
import { useOpenPersonas } from '@/features/collection-runner/hooks/useOpenPersonas'
import { useRunCollection } from '@/features/collection-runner/hooks/useRunCollection'
import type { Collection } from '@/features/collection-runner/types'
import { getErrorMessage } from '@/lib/errors'
import { reportError } from '@/lib/reportError'

export interface CollectionRunnerPanelProps {
  /** Collection chuni gayi Collections sidebar se — null jab tak kuch select na ho. */
  selected: { path: string; collection: Collection } | null
}

function parsePositiveInt(s: string): number | undefined {
  const n = parseInt(s, 10)
  return n > 0 ? n : undefined
}

function parsePositiveFloat(s: string): number | undefined {
  const n = parseFloat(s)
  return n > 0 ? n : undefined
}

/** Postman-lite ka doosra half: selected collection (+ optional environment, load-test knobs) run karo, results dekho. */
export function CollectionRunnerPanel({ selected }: CollectionRunnerPanelProps) {
  const [envPath, setEnvPath] = useState('')
  const [personasPath, setPersonasPath] = useState('')
  const [workers, setWorkers] = useState('')
  const [iterations, setIterations] = useState('')
  const [durationSec, setDurationSec] = useState('')
  const [rps, setRps] = useState('')

  const {
    data: environment,
    error: envError,
    isLoading: isOpeningEnv,
    run: openEnv,
  } = useOpenEnvironment()
  const {
    data: personas,
    error: personasError,
    isLoading: isOpeningPersonas,
    run: openPersonas,
  } = useOpenPersonas()
  const {
    data: output,
    error: runError,
    isLoading: isRunning,
    run: runCollection,
  } = useRunCollection()

  async function handleOpenEnv() {
    await openEnv(envPath)
  }

  async function handleBrowseEnv() {
    const result = await pickFile('Open environment', 'json')
    if (!result.ok) return reportError(result.error, 'browse environment file')
    if (!result.value) return // user cancelled
    setEnvPath(result.value)
    await openEnv(result.value)
  }

  async function handleOpenPersonas() {
    await openPersonas(personasPath)
  }

  async function handleBrowsePersonas() {
    const result = await pickFile('Open personas', 'csv')
    if (!result.ok) return reportError(result.error, 'browse personas file')
    if (!result.value) return // user cancelled
    setPersonasPath(result.value)
    await openPersonas(result.value)
  }

  async function handleRun() {
    if (!selected) return
    const durationSeconds = parsePositiveFloat(durationSec)
    await runCollection({
      collection: selected.collection,
      envVariables: environment?.variables,
      workers: parsePositiveInt(workers),
      iterations: parsePositiveInt(iterations),
      durationMs: durationSeconds ? Math.round(durationSeconds * 1000) : undefined,
      rps: parsePositiveFloat(rps),
      personas,
    })
  }

  return (
    <div className="flex flex-col gap-4">
      <Card>
        <CardHeader>
          <CardTitle>Run a collection</CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col gap-4">
          {!selected && (
            <EmptyState
              title="No collection selected"
              description="Pick one from the Collections sidebar, or click + Add to browse for one."
            />
          )}

          {selected && (
            <>
              <div className="flex items-center justify-between">
                <span className="text-fg text-sm font-medium">
                  {selected.collection.name || 'Untitled collection'} —{' '}
                  {selected.collection.requests.length} request
                  {selected.collection.requests.length === 1 ? '' : 's'}
                </span>
                <Button variant="primary" isLoading={isRunning} onClick={() => void handleRun()}>
                  Run
                </Button>
              </div>
              <div className="border-border bg-surface flex flex-col divide-y rounded-md border">
                {selected.collection.requests.map((req, i) => (
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

              <div className="grid grid-cols-4 gap-2">
                <Field label="Workers (-c)">
                  <Input
                    type="number"
                    min={1}
                    value={workers}
                    onChange={(e) => setWorkers(e.target.value)}
                    placeholder="1"
                    aria-label="Workers"
                  />
                </Field>
                <Field label="Iterations (-n)">
                  <Input
                    type="number"
                    min={1}
                    value={iterations}
                    onChange={(e) => setIterations(e.target.value)}
                    placeholder="1"
                    aria-label="Iterations"
                  />
                </Field>
                <Field label="Duration, sec (-d)">
                  <Input
                    type="number"
                    min={0}
                    value={durationSec}
                    onChange={(e) => setDurationSec(e.target.value)}
                    placeholder="none"
                    aria-label="Duration in seconds"
                  />
                </Field>
                <Field label="Rate limit (rps)">
                  <Input
                    type="number"
                    min={0}
                    value={rps}
                    onChange={(e) => setRps(e.target.value)}
                    placeholder="unlimited"
                    aria-label="Requests per second"
                  />
                </Field>
              </div>
            </>
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

          <div className="flex gap-2">
            <Input
              value={personasPath}
              onChange={(e) => setPersonasPath(e.target.value)}
              placeholder="personas.csv (optional, {{persona.<col>}})"
              className="flex-1"
              aria-label="Personas file path"
            />
            <Button
              variant="outline"
              aria-label="Browse for personas file"
              onClick={() => void handleBrowsePersonas()}
            >
              Browse…
            </Button>
            <Button
              variant="secondary"
              isLoading={isOpeningPersonas}
              disabled={!personasPath.trim()}
              onClick={() => void handleOpenPersonas()}
            >
              Load personas
            </Button>
          </div>

          {personasError && (
            <Alert variant="danger" title={`Could not load personas (${personasError.kind})`}>
              {getErrorMessage(personasError)}
            </Alert>
          )}

          {personas && (
            <span className="text-fg text-sm font-medium">
              {personas.length} persona{personas.length === 1 ? '' : 's'} loaded
              {personas.length > 1 ? ' — assigned round-robin across workers' : ''}
            </span>
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
              {output.stats.map((s, i) => (
                <div key={`${s.name}-${i}`} className="flex items-center gap-2 px-3 py-1.5 text-sm">
                  <Badge
                    variant={s.failures > 0 ? 'outline' : 'default'}
                    className="shrink-0 justify-center"
                  >
                    {s.successes}/{s.totalRuns}
                  </Badge>
                  <span className="text-fg-muted truncate">{s.name}</span>
                  <span className="text-fg-subtle ml-auto shrink-0 text-xs">
                    {s.topError ?? `${s.avgLatencyMs} ms avg / ${s.p95LatencyMs} ms p95`}
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
