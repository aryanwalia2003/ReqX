# Shared UI components

Presentational components with no feature-specific logic and no direct
`@wails/*` calls — a `Button`, a `Tooltip`, a `Table`, the kind of thing two
or more features would otherwise duplicate. Run `npm run gen:component` to
scaffold a new one with the right naming/export shape already in place.

If a component only makes sense inside one feature, it belongs in that
feature's own `components/` folder instead.

## Usage

Everything here is re-exported from `@/components/ui` — import from there,
not from the individual file:

```tsx
import { Button, Field, Input, useToast } from '@/components/ui'
```

## Theme

Colors are design tokens in `src/styles/globals.css` (`@theme`), strictly
grayscale — `bg`, `surface`, `surface-2`, `surface-3`, `border`,
`border-strong`, `fg`, `fg-muted`, `fg-subtle`, `ring`. No other hue exists
in the app. Semantic state (success/warning/danger) is signalled by icon
shape and border style (`solid` vs `dashed`), never by color — keep that
convention when adding components.

Every interactive component takes a plain `className` prop (merged via
`cn()`, so Tailwind conflicts resolve correctly) and forwards a `ref` where
the underlying DOM node is focusable/measurable — no separate "ref" API to
learn.

## What's here

- **Form**: `Button`, `Input`, `Textarea`, `Select`, `Checkbox`, `Radio`,
  `Switch`, `Label`, `Field` (wraps a control with label/hint/error and
  wires up the `id`/`aria-*` for you)
- **Layout**: `Card` (+ `CardHeader`/`CardTitle`/`CardDescription`/
  `CardContent`/`CardFooter`), `Separator`, `EmptyState`
- **Feedback**: `Alert`, `Badge`, `Spinner`, `Skeleton`, `Toast` (wrap the
  app once in `<ToastProvider>`, then call `useToast()` anywhere)
- **Overlay**: `Dialog` (+ `DialogHeader`/`DialogTitle`/
  `DialogDescription`/`DialogContent`/`DialogFooter`), `Tooltip`
- **Navigation**: `Tabs` (+ `TabsList`/`TabsTrigger`/`TabsContent`)
- **Misc**: `Kbd`, `icons` (inline SVGs — no icon package installed)
