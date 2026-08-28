import { useEffect, useState } from 'react'

import {
  Alert,
  Badge,
  Button,
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  EmptyState,
  Spinner,
} from '@/components/ui'
import { useListRuns } from '@/features/history/hooks/useListRuns'
import { useRunStats } from '@/features/history/hooks/useRunStats'
import { getErrorMessage } from '@/lib/errors'

/** CLI ke `reqx ui` dashboard jaisa hi — past runs list, click karke drilldown. */
export function HistoryPanel() {
  const [selectedRunId, setSelectedRunId] = useState<string | null>(null)

  const { data: runs, error: listError, isLoading: isListing, run: list } = useListRuns()
  const {
    data: stats,
    error: statsError,
    isLoading: isLoadingStats,
    run: loadStats,
  } = useRunStats()

  useEffect(() => {
    void list()
    // eslint-disable-next-line react-hooks/exhaustive-deps -- ek baar mount pe load
  }, [])

  function handleSelect(runId: string) {
    setSelectedRunId(runId)
    void loadStats(runId)
  }

  return (
    <div className="flex flex-col gap-4">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center justify-between">
            <span>Run history</span>
            <Button variant="outline" size="sm" isLoading={isListing} onClick={() => void list()}>
              Refresh
            </Button>
          </CardTitle>
        </CardHeader>
        <CardContent>
          {listError && (
            <Alert variant="danger" title={`Could not load history (${listError.kind})`}>
              {getErrorMessage(listError)}
            </Alert>
          )}

          {!listError && isListing && !runs && (
            <div className="flex justify-center py-6">
              <Spinner />
            </div>
          )}

          {!listError && runs && runs.length === 0 && (
            <EmptyState
              title="No runs yet"
              description="Run a collection from the Run collection tab — it'll show up here."
            />
          )}

          {runs && runs.length > 0 && (
            <div className="border-border bg-surface flex flex-col divide-y rounded-md border">
              {runs.map((run) => (
                <button
                  key={run.id}
                  type="button"
                  onClick={() => handleSelect(run.id)}
                  className={`hover:bg-surface-2 flex w-full items-center gap-2 px-3 py-1.5 text-left text-sm transition-colors ${
                    selectedRunId === run.id ? 'bg-surface-2' : ''
                  }`}
                >
                  <span className="text-fg-muted truncate">{run.collection}</span>
                  <span className="text-fg-subtle shrink-0 text-xs">{run.ts}</span>
                  <span className="ml-auto flex shrink-0 items-center gap-1.5">
                    <Badge variant="outline">{run.total_reqs} reqs</Badge>
                    <Badge variant="subtle">{run.p95_ms} ms p95</Badge>
                    <Badge variant={run.error_pct > 0 ? 'outline' : 'default'}>
                      {run.error_pct.toFixed(0)}% err
                    </Badge>
                  </span>
                </button>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {selectedRunId && (
        <Card>
          <CardHeader>
            <CardTitle>Per-request breakdown</CardTitle>
          </CardHeader>
          <CardContent>
            {statsError && (
              <Alert variant="danger" title={`Could not load run detail (${statsError.kind})`}>
                {getErrorMessage(statsError)}
              </Alert>
            )}
            {isLoadingStats && (
              <div className="flex justify-center py-6">
                <Spinner />
              </div>
            )}
            {stats && stats.length === 0 && !isLoadingStats && (
              <EmptyState title="No per-request stats for this run" />
            )}
            {stats && stats.length > 0 && (
              <div className="border-border bg-surface flex flex-col divide-y rounded-md border">
                {stats.map((s, i) => (
                  <div
                    key={`${s.name}-${i}`}
                    className="flex items-center gap-2 px-3 py-1.5 text-sm"
                  >
                    <span className="text-fg-muted truncate">{s.name}</span>
                    <span className="ml-auto flex shrink-0 items-center gap-1.5 text-xs">
                      <Badge variant="default">{s.successes} ok</Badge>
                      <Badge variant={s.failures > 0 ? 'outline' : 'subtle'}>
                        {s.failures} failed
                      </Badge>
                      <Badge variant="subtle">{s.avg_ms} ms avg</Badge>
                      <Badge variant="subtle">{s.p95_ms} ms p95</Badge>
                    </span>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      )}
    </div>
  )
}
