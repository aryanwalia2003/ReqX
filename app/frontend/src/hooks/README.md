# Shared hooks

Hooks used by two or more features. Run `npm run gen:hook` to scaffold a new
one (adds the file, a test file, and the barrel export automatically).
Feature-specific hooks belong in that feature's own `hooks/` folder instead.

## Usage

Everything here is re-exported from `@/hooks` — import from there, not from
the individual file:

```tsx
import { useAsync, useDebounce, useDisclosure } from '@/hooks'
```

## What's here

- **`useAsync`** — wraps a `Result<T, E>`-returning async fn (a Wails binding
  call through `toResult`, per `src/lib/result.ts`) in `{ data, error,
isLoading, run }`. `error` is an `AppError` by default (`src/lib/errors.ts`) —
  render it with `<Alert variant="danger">{getErrorMessage(error)}</Alert>`.
  Stale calls — a newer `run()` starting before an older one resolves — are
  ignored automatically. This is the go-to way to call into
  `@wails/go/services/*`:

  ```tsx
  const { data, error, isLoading, run } = useAsync(
    (name: string) => toResult(window.go.services.ExampleService.ping({ name })),
    { onError: (error) => reportError(error, 'ExampleService.ping') },
  )
  useEffect(() => {
    void run('dev')
  }, [run])
  ```

  `onError` is optional — skip it and just render the returned `error` when a
  toast/log isn't needed.

- **`useDebounce(value, delay)`** — debounce a fast-changing value (search
  boxes, URL params typed live).
- **`useEventListener(eventName, handler, ref?)`** — DOM/window listener that
  always calls the latest `handler` without re-attaching; defaults to
  `window`, pass a ref to target an element instead.
- **`useClickOutside(ref, handler)`** — fire `handler` on an outside
  click; used for dismissing dropdowns/popovers/dialogs.
- **`useKeyboardShortcut(combo, handler)`** — global shortcut, e.g.
  `useKeyboardShortcut('mod+enter', sendRequest)` (`mod` = Cmd on Mac, Ctrl
  elsewhere — see `<Kbd>` in `components/ui`).
- **`useDisclosure(initialOpen?)`** — `{ isOpen, open, close, toggle }` for
  dialog/popover open state, instead of a raw boolean `useState`.
- **`useLocalStorage(key, initialValue)`** — persisted `[value, setValue]`
  state, JSON-serialized automatically.
- **`useCopyToClipboard()`** — `[status, copy]`, `status` is
  `'idle' | 'copied' | 'error'`.
- **`usePrevious(value)`** — the value from the previous render (`undefined`
  on the first one).
- **`useIsMounted()`** — returns a `() => boolean`; guard a `setState` inside
  an async callback that may resolve after unmount.
- **`useMachine(machine)`** — wires a `createMachine` config
  (`src/lib/machine.ts`) into React state: `{ state, context, send, matches,
can }`. Reach for a machine instead of a pile of booleans when a
  component's states + valid transitions between them matter — e.g. a
  connection that can only go `idle → connecting → open`, never straight to
  `open`:

  ```tsx
  const { state, send, matches } = useMachine(connectionMachine)
  <Button disabled={!matches('idle')} onClick={() => send({ type: 'CONNECT' })}>
    Connect
  </Button>
  ```

  Run `npm run gen:machine` to scaffold a new shared one under
  `src/lib/machines/` — see `src/lib/machines/README.md`.
