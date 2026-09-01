import { useState } from 'react'

import {
  Alert,
  Badge,
  Button,
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  Checkbox,
  EmptyState,
  Field,
  Input,
  Select,
} from '@/components/ui'
import { pickFile, pickSaveFile } from '@/features/collection-runner/api'
import { OpenInEditorButtons } from '@/features/collection-runner/components/OpenInEditorButtons'
import { useOpenEnvironment } from '@/features/collection-runner/hooks/useOpenEnvironment'
import { useOpenPersonas } from '@/features/collection-runner/hooks/useOpenPersonas'
import { useRunCollection } from '@/features/collection-runner/hooks/useRunCollection'
import type { Collection } from '@/features/collection-runner/types'
import type { DagNodeRow } from '@/features/history'
import { getErrorMessage } from '@/lib/errors'
import { reportError } from '@/lib/reportError'

export interface CollectionRunnerPanelProps {
  /** Collection chuni gayi Collections sidebar se — null jab tak kuch select na ho. */
  selected: { path: string; collection: Collection } | null
}

const INJECT_METHODS = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE']

function parsePositiveInt(s: string): number | undefined {
  const n = parseInt(s, 10)
  return n > 0 ? n : undefined
}

function parsePositiveFloat(s: string): number | undefined {
  const n = parseFloat(s)
  return n > 0 ? n : undefined
}

/** DAG nodes ko level_idx se group karta, level order me sorted. */
function groupByLevel(nodes: DagNodeRow[]): [number, DagNodeRow[]][] {
  const byLevel = new Map<number, DagNodeRow[]>()
  for (const node of nodes) {
    const group = byLevel.get(node.level_idx) ?? []
    group.push(node)
    byLevel.set(node.level_idx, group)
  }
  return [...byLevel.entries()].sort(([a], [b]) => a - b)
}

/** Postman-lite ka doosra half: selected collection (+ optional environment, load-test knobs) run karo, results dekho. */
export function CollectionRunnerPanel({ selected }: CollectionRunnerPanelProps) {
  const [envPath, setEnvPath] = useState('')
  const [personasPath, setPersonasPath] = useState('')
  const [workers, setWorkers] = useState('')
  const [iterations, setIterations] = useState('')
  const [durationSec, setDurationSec] = useState('')
  const [rps, setRps] = useState('')

  const [showAdvanced, setShowAdvanced] = useState(false)
  const [stages, setStages] = useState('')
  const [noCookies, setNoCookies] = useState(false)
  const [clearCookies, setClearCookies] = useState(false)
  const [graphqlCheck, setGraphqlCheck] = useState(false)
  const [exportPath, setExportPath] = useState('')
  const [injectIndex, setInjectIndex] = useState('')
  const [injectName, setInjectName] = useState('')
  const [injectMethod, setInjectMethod] = useState('GET')
  const [injectUrl, setInjectUrl] = useState('')

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

  async function handleBrowseExport() {
    const result = await pickSaveFile('Export results', 'results.ndjson')
    if (!result.ok) return reportError(result.error, 'browse export save location')
    if (!result.value) return // user cancelled
    setExportPath(result.value)
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
      stages: stages.trim() || undefined,
      noCookies: noCookies || undefined,
      clearCookies: clearCookies || undefined,
      graphql: graphqlCheck || undefined,
      exportPath: exportPath.trim() || undefined,
      injectIndex: injectIndex.trim() || undefined,
      injectName: injectName.trim() || undefined,
      injectMethod: injectName.trim() ? injectMethod : undefined,
      injectUrl: injectUrl.trim() || undefined,
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
                <div className="flex items-center gap-2">
                  <OpenInEditorButtons path={selected.path} />
                  <Button variant="primary" isLoading={isRunning} onClick={() => void handleRun()}>
                    Run
                  </Button>
                </div>
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

              <div className="flex flex-col gap-2">
                <button
                  type="button"
                  onClick={() => setShowAdvanced((v) => !v)}
                  aria-label="Advanced options"
                  aria-expanded={showAdvanced}
                  className="text-fg-muted hover:text-fg flex items-center gap-1 text-xs font-medium"
                >
                  <span aria-hidden="true">{showAdvanced ? '▾' : '▸'}</span>
                  Advanced options
                </button>

                {showAdvanced && (
                  <div className="border-border bg-surface flex flex-col gap-3 rounded-md border p-3">
                    <Field
                      label="Stages (--stages)"
                      hint='e.g. "10s:5,30s:20,10s:0" — ramps workers over time, overrides Workers/Duration above'
                    >
                      <Input
                        value={stages}
                        onChange={(e) => setStages(e.target.value)}
                        placeholder="10s:5,30s:20,10s:0"
                        aria-label="Stages"
                      />
                    </Field>

                    <div className="flex flex-wrap gap-4">
                      <label className="flex items-center gap-1.5 text-sm">
                        <Checkbox
                          checked={noCookies}
                          onChange={(e) => setNoCookies(e.target.checked)}
                        />
                        No cookies
                      </label>
                      <label className="flex items-center gap-1.5 text-sm">
                        <Checkbox
                          checked={clearCookies}
                          onChange={(e) => setClearCookies(e.target.checked)}
                        />
                        Clear cookies per request
                      </label>
                      <label className="flex items-center gap-1.5 text-sm">
                        <Checkbox
                          checked={graphqlCheck}
                          onChange={(e) => setGraphqlCheck(e.target.checked)}
                        />
                        GraphQL error detection
                      </label>
                    </div>

                    <div className="flex gap-2">
                      <Input
                        value={exportPath}
                        onChange={(e) => setExportPath(e.target.value)}
                        placeholder="Export results to… (optional)"
                        className="flex-1"
                        aria-label="Export path"
                      />
                      <Button variant="outline" onClick={() => void handleBrowseExport()}>
                        Browse (save)…
                      </Button>
                    </div>

                    <div className="flex flex-col gap-1.5">
                      <span className="text-fg text-xs font-medium">
                        Inject a request (temporary — collection file untouched)
                      </span>
                      <div className="flex gap-2">
                        <Input
                          type="number"
                          min={1}
                          value={injectIndex}
                          onChange={(e) => setInjectIndex(e.target.value)}
                          placeholder="#"
                          className="w-14"
                          aria-label="Inject position"
                        />
                        <Select
                          value={injectMethod}
                          onChange={(e) => setInjectMethod(e.target.value)}
                          className="w-28"
                          aria-label="Inject method"
                        >
                          {INJECT_METHODS.map((m) => (
                            <option key={m} value={m}>
                              {m}
                            </option>
                          ))}
                        </Select>
                        <Input
                          value={injectName}
                          onChange={(e) => setInjectName(e.target.value)}
                          placeholder="Name"
                          className="flex-1"
                          aria-label="Inject name"
                        />
                        <Input
                          value={injectUrl}
                          onChange={(e) => setInjectUrl(e.target.value)}
                          placeholder="URL"
                          className="flex-1"
                          aria-label="Inject URL"
                        />
                      </div>
                    </div>
                  </div>
                )}
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
            <OpenInEditorButtons path={envPath.trim() || undefined} />
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

      {output?.dagNodes && output.dagNodes.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Execution graph</CardTitle>
          </CardHeader>
          <CardContent className="flex flex-col gap-3">
            {groupByLevel(output.dagNodes).map(([level, nodes]) => (
              <div key={level} className="flex flex-col gap-1">
                <span className="text-fg-subtle text-xs font-medium">Level {level}</span>
                <div className="flex flex-wrap gap-2">
                  {nodes.map((node) => (
                    <Badge
                      key={node.name}
                      variant={
                        node.status === 'failed'
                          ? 'outline'
                          : node.status === 'skipped'
                            ? 'subtle'
                            : 'default'
                      }
                    >
                      {node.name} · {node.status} · {node.duration_ms} ms
                    </Badge>
                  ))}
                </div>
              </div>
            ))}
          </CardContent>
        </Card>
      )}
    </div>
  )
}
