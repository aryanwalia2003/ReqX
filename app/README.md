# ReqX desktop app — backend

The [Wails](https://wails.io) v2 shell around ReqX's existing Go engine
(`../internal/*`). This is scaffolding only right now: `services/` holds one
template service, nothing real is wired up yet.

## Layout

```
app/
  main.go          # Wails bootstrap — embeds frontend/dist, binds services
  wailsapp/        # App struct + service wiring (app_struct/ctor/method.go)
  services/        # one Wails-bound service per feature
  frontend/        # the React/Vite UI — see frontend/README.md
  wails.json       # Wails CLI project config
```

Same file-per-concept convention as the rest of the repo
(`internal/runner/collection_runner_{struct,ctor,method}.go` etc.) —
`app/wailsapp/app_{struct,ctor,method}.go`, and each service gets its own
`*_struct.go` / `*_ctor.go` / `*_method.go` trio. Copy
`services/example_service_*.go` as the starting point for a new one.

## Conventions

- **Errors go through `internal/errs`** — the same `errs.Wrap`/`errs.Kind*`
  convention the CLI already uses (see `cmd/run_cmd_ctor.go`). Never return a
  bare `fmt.Errorf` from a bound method.
- **Bound methods take and return plain structs**, never bare primitives —
  Wails generates matching TypeScript types from them for the frontend.
- **Go → JS naming**: Wails exposes a `PascalCase` Go method as `camelCase`
  in the generated JS binding (`ExampleService.Ping` →
  `window.go.services.ExampleService.ping(...)`). Don't be surprised by the
  case change.
- **No business logic here.** `app/services/*.go` should only translate
  between the frontend and `internal/*` — the engine itself stays in
  `internal/`, shared with the CLI.

## Getting started

Requires the [Wails CLI](https://wails.io/docs/gettingstarted/installation)
(`go install github.com/wailsapp/wails/v2/cmd/wails@latest`) — not installed
in this environment yet, so `wails dev` / `wails build` are unverified here.
`go build ./...` at the repo root does work today: a placeholder
`frontend/dist/index.html` is committed specifically so the `//go:embed`
in `main.go` has something to embed before the frontend has ever been built.

```sh
make dev     # from the repo root — see ../Makefile
# or, once the Wails CLI is installed:
cd app && wails dev
```
