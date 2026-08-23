# Contributing to ReqX

## Repo map

| Path | What it is |
|---|---|
| `cmd/` | The CLI (`reqx run`, `reqx ui`, `reqx sio`, ...) |
| `internal/` | The shared engine — HTTP/WS/Socket.IO execution, scripting, scheduler, history. Both the CLI and the desktop app use this. |
| `app/` | The Wails desktop app's Go backend — see `app/README.md` |
| `app/frontend/` | The desktop app's React/TS UI — see `app/frontend/README.md` |
| `docs/` | Design docs and architecture notes |

## Setup

```sh
go build ./...             # CLI + app backend both compile
cd app/frontend && npm install   # also installs the repo-root git hooks
```

## Everyday commands

Run `make` targets from the repo root — see the [Makefile](Makefile) — or the
per-project `npm run <script>` from `app/frontend/`.

## Commits

Conventional commits (`feat:`, `fix:`, `chore:`, `perf:`, ...) — enforced by
a commit-msg hook, matches this repo's existing `git log`.

## Before opening a PR

- `go build ./...` and `go vet ./...` are clean
- `cd app/frontend && npm run typecheck && npm run lint && npm run build` are clean
- Existing CLI behavior is unaffected unless that's the point of the PR
