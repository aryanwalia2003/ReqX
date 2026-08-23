import { useState } from 'react'

import {
  Alert,
  Badge,
  Button,
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
  Checkbox,
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  Field,
  Input,
  Kbd,
  Select,
  Separator,
  Switch,
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
  Textarea,
  useToast,
} from '@/components/ui'

/**
 * Living demo of the shared UI kit — har naya component yahan ek line me
 * dikhna chahiye. `npm run gen:feature` se real feature banao.
 */
function App() {
  const [dialogOpen, setDialogOpen] = useState(false)
  const [notifyOn, setNotifyOn] = useState(true)
  const toast = useToast()

  return (
    <main className="mx-auto flex max-w-3xl flex-col gap-8 p-8">
      <header className="flex flex-col gap-1">
        <h1 className="text-2xl font-semibold tracking-tight">ReqX</h1>
        <p className="text-fg-muted text-sm">
          Shared UI kit demo — see{' '}
          <code className="bg-surface rounded px-1.5 py-0.5">src/components/ui</code>.
        </p>
      </header>

      <section className="flex flex-wrap items-center gap-2">
        <Button variant="primary">Send request</Button>
        <Button variant="secondary">Save</Button>
        <Button variant="outline">Duplicate</Button>
        <Button variant="ghost">Cancel</Button>
        <Button variant="destructive" onClick={() => setDialogOpen(true)}>
          Delete collection
        </Button>
        <Button variant="secondary" isLoading>
          Sending
        </Button>
      </section>

      <Separator />

      <section className="flex flex-col gap-4">
        <Field label="Request name" hint="Sirf tumhe dikhega, team ko nahi.">
          <Input placeholder="Get user profile" />
        </Field>
        <Field label="Environment" required>
          <Select defaultValue="local">
            <option value="local">Local</option>
            <option value="staging">Staging</option>
            <option value="prod">Production</option>
          </Select>
        </Field>
        <Field label="Body" error="Valid JSON nahi hai.">
          <Textarea mono rows={4} defaultValue={'{\n  "id": 1\n}'} />
        </Field>
        <div className="flex items-center gap-2">
          <Checkbox id="verify-ssl" defaultChecked />
          <label htmlFor="verify-ssl" className="text-fg text-sm">
            Verify SSL certificates
          </label>
        </div>
        <div className="flex items-center gap-2">
          <Switch checked={notifyOn} onChange={(e) => setNotifyOn(e.target.checked)} />
          <span className="text-fg text-sm">Notify on response error</span>
        </div>
      </section>

      <Separator />

      <section className="flex flex-wrap items-center gap-2">
        <Badge>GET</Badge>
        <Badge variant="outline">200 OK</Badge>
        <Badge variant="subtle">3 headers</Badge>
        <Kbd>⌘</Kbd>
        <Kbd>Enter</Kbd>
        <span className="text-fg-muted text-sm">to send</span>
      </section>

      <Alert variant="warning" title="Token expires soon">
        Refresh your auth token before this request stops working.
      </Alert>

      <Card>
        <CardHeader>
          <CardTitle>Collections</CardTitle>
        </CardHeader>
        <CardContent>
          <Tabs defaultValue="requests">
            <TabsList>
              <TabsTrigger value="requests">Requests</TabsTrigger>
              <TabsTrigger value="env">Environments</TabsTrigger>
            </TabsList>
            <TabsContent value="requests" className="text-fg-muted pt-3 text-sm">
              12 requests in this collection.
            </TabsContent>
            <TabsContent value="env" className="text-fg-muted pt-3 text-sm">
              3 environments configured.
            </TabsContent>
          </Tabs>
        </CardContent>
        <CardFooter>
          <Button
            size="sm"
            variant="secondary"
            onClick={() => toast({ title: 'Collection saved', variant: 'success' })}
          >
            Save changes
          </Button>
        </CardFooter>
      </Card>

      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogHeader>
          <DialogTitle>Delete collection?</DialogTitle>
          <DialogDescription>This removes every request inside it. No undo.</DialogDescription>
        </DialogHeader>
        <DialogContent />
        <DialogFooter>
          <Button variant="ghost" onClick={() => setDialogOpen(false)}>
            Cancel
          </Button>
          <Button
            variant="destructive"
            onClick={() => {
              setDialogOpen(false)
              toast({ title: 'Collection deleted', variant: 'danger' })
            }}
          >
            Delete
          </Button>
        </DialogFooter>
      </Dialog>
    </main>
  )
}

export default App
