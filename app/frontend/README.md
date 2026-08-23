# ReqX desktop app — frontend

React + TypeScript + Vite, running inside a [Wails](https://wails.io) native
webview (see `../` for the Go side). This is scaffolding only right now —
no feature is wired up yet.

## Getting started

```sh
npm install       # also sets up the repo-root git hooks (see "prepare" script)
npm run dev        # Vite dev server, standalone (no Go backend yet)
```

Once the Go side is wired up, the app runs as a whole via `wails dev` from
`app/` — see `../README.md`.

## Conventions

- **Never hand-write boilerplate.** `npm run gen:feature`, `gen:component`,
  `gen:hook` scaffold new code with the right shape already in place. If
  you're copy-pasting an existing file as a starting point, there should
  probably be a generator for that instead — add one.
- **Import via the `@/` alias**, never a deep relative `../../../` path —
  ESLint enforces this. Wails-generated bindings (once they exist) live
  under `./wailsjs` and are reached via the `@wails/` alias the same way.
- **One barrel per feature.** Each `src/features/<name>/index.ts` is the
  only file other code may import from that feature. Nothing else in the
  folder is a stable import target.
- **Naming**: components `PascalCase.tsx`, hooks `useThing.ts`, everything
  else `kebab-case.ts`. One component per file.
- **Errors are values, not exceptions.** See `src/lib/result.ts` — a bound
  Wails call becomes `toResult(SomeGoMethod(...))`, not a `try/catch`.
- **Every folder that isn't obvious has a `README.md`** explaining what
  belongs there — read the nearest one before guessing.

## Layout

```
src/
  app/          # shell: providers, root layout — see src/app/README.md
  features/     # one folder per feature slice — see src/features/README.md
  components/ui/ # shared presentational components
  hooks/         # shared hooks
  lib/           # framework-agnostic utilities (result.ts, cn.ts)
  types/         # shared/global types + the @wails ambient placeholder
  styles/        # Tailwind entry + design tokens
plop-templates/  # source templates for the gen:* generators
```

## Tooling

| Command                                              | What it does                               |
| ---------------------------------------------------- | ------------------------------------------ |
| `npm run dev`                                        | Vite dev server                            |
| `npm run build`                                      | typecheck + production build               |
| `npm run typecheck`                                  | `tsc` only, no emit                        |
| `npm run lint` / `lint:fix`                          | ESLint                                     |
| `npm run format` / `format:check`                    | Prettier (auto-sorts Tailwind classes too) |
| `npm run test` / `test:watch`                        | Vitest + React Testing Library             |
| `npm run gen:feature` / `gen:component` / `gen:hook` | Plop generators                            |

A pre-commit hook (repo-root `.husky/`) runs ESLint + Prettier on staged
files automatically — you shouldn't need to run lint/format by hand before
committing.

## Design tokens

Colors in `src/styles/globals.css` are ported from the existing `reqx ui`
dashboard (`internal/ui/assets/index.html`) so the new app doesn't invent a
second, inconsistent theme.
