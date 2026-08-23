# Shared state machines

Machine configs (`createMachine`, see `src/lib/machine.ts`) used by two or
more features. Run `npm run gen:machine` to scaffold one. Feature-specific
machines belong in that feature's own folder instead (e.g.
`src/features/socket-debugger/machines/connection.ts`), hand-written with
the same `createMachine` — there's no generator for those since not every
feature needs one.

Wire a machine into a component with `useMachine` — see `src/hooks/README.md`.
