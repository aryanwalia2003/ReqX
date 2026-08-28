import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui'
import { CollectionRunnerPanel } from '@/features/collection-runner'
import { SendRequestPanel } from '@/features/send-request'

function App() {
  return (
    <main className="mx-auto flex max-w-3xl flex-col gap-6 p-8">
      <header className="flex flex-col gap-1">
        <h1 className="text-2xl font-semibold tracking-tight">ReqX</h1>
        <p className="text-fg-muted text-sm">Send a request, or run a whole collection.</p>
      </header>

      <Tabs defaultValue="send">
        <TabsList>
          <TabsTrigger value="send">Send request</TabsTrigger>
          <TabsTrigger value="run">Run collection</TabsTrigger>
        </TabsList>
        <TabsContent value="send" className="pt-4">
          <SendRequestPanel />
        </TabsContent>
        <TabsContent value="run" className="pt-4">
          <CollectionRunnerPanel />
        </TabsContent>
      </Tabs>
    </main>
  )
}

export default App
