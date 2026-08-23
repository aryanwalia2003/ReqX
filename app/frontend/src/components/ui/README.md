# Shared UI components

Presentational components with no feature-specific logic and no direct
`@wails/*` calls — a `Button`, a `Tooltip`, a `Table`, the kind of thing two
or more features would otherwise duplicate. Run `npm run gen:component` to
scaffold a new one with the right naming/export shape already in place.

If a component only makes sense inside one feature, it belongs in that
feature's own `components/` folder instead.
