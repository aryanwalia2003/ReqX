# Features

One folder per feature-slice — e.g. `collections/`, `run/`, `socket-debugger/`,
`history/`. Don't hand-write a new one: run `npm run gen:feature` from
`app/frontend/` and answer the prompt.

Each generated feature folder gets:

```
<feature>/
  components/   # feature-local components (not shared — put those in src/components/ui)
  hooks/        # feature-local hooks
  api/          # thin wrappers around @wails/go/services/* bindings, returning Result<T>
  types.ts      # feature-local types
  index.ts      # the ONLY thing other features may import — re-export the public API here
```

Nothing outside `index.ts`'s exports is a stable import target — treat every
other file in a feature folder as private to it.
