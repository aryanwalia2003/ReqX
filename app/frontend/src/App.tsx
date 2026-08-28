import { SendRequestPanel } from '@/features/send-request'

function App() {
  return (
    <main className="mx-auto flex max-w-3xl flex-col gap-6 p-8">
      <header className="flex flex-col gap-1">
        <h1 className="text-2xl font-semibold tracking-tight">ReqX</h1>
        <p className="text-fg-muted text-sm">Send a request and see what comes back.</p>
      </header>

      <SendRequestPanel />
    </main>
  )
}

export default App
