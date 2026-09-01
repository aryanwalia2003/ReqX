import { useState } from 'react'

import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui'
import { CollectionRunnerPanel, CollectionsSidebar } from '@/features/collection-runner'
import type { Collection } from '@/features/collection-runner'
import { HistoryPanel } from '@/features/history'
import { SendRequestPanel } from '@/features/send-request'
import type { SendRequestInput } from '@/features/send-request'
import { SocketDebuggerPanel } from '@/features/socket-debugger'

function App() {
  const [activeTab, setActiveTab] = useState('send')
  const [selected, setSelected] = useState<{ path: string; collection: Collection } | null>(null)
  const [loadedRequest, setLoadedRequest] = useState<SendRequestInput | null>(null)
  // SendRequestPanel only reads loadRequest on mount — bump this to force a
  // remount (fresh initial state) even when Send tab is already open.
  const [loadedRequestNonce, setLoadedRequestNonce] = useState(0)

  function handleSelectCollection(path: string, collection: Collection) {
    setSelected({ path, collection })
    setActiveTab('run')
  }

  function handleOpenRequest(request: SendRequestInput) {
    setLoadedRequest(request)
    setLoadedRequestNonce((n) => n + 1)
    setActiveTab('send')
  }

  return (
    <div className="flex h-screen">
      <CollectionsSidebar
        selectedPath={selected?.path ?? null}
        onSelectCollection={handleSelectCollection}
        onOpenRequest={handleOpenRequest}
      />

      <main className="flex-1 overflow-y-auto">
        <div className="mx-auto flex max-w-3xl flex-col gap-6 p-8">
          <header className="flex flex-col gap-1">
            <h1 className="text-2xl font-semibold tracking-tight">ReqX</h1>
            <p className="text-fg-muted text-sm">
              Send a request, run a collection, or check history.
            </p>
          </header>

          <Tabs value={activeTab} onValueChange={setActiveTab}>
            <TabsList>
              <TabsTrigger value="send">Send request</TabsTrigger>
              <TabsTrigger value="run">Run collection</TabsTrigger>
              <TabsTrigger value="socket">Socket debugger</TabsTrigger>
              <TabsTrigger value="history">History</TabsTrigger>
            </TabsList>
            <TabsContent value="send" className="pt-4">
              <SendRequestPanel key={loadedRequestNonce} loadRequest={loadedRequest} />
            </TabsContent>
            <TabsContent value="run" className="pt-4">
              <CollectionRunnerPanel selected={selected} />
            </TabsContent>
            <TabsContent value="socket" className="pt-4">
              <SocketDebuggerPanel />
            </TabsContent>
            <TabsContent value="history" className="pt-4">
              <HistoryPanel />
            </TabsContent>
          </Tabs>
        </div>
      </main>
    </div>
  )
}

export default App
