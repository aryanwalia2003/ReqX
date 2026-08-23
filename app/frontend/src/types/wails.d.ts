/**
 * Placeholder ambient module for the "@wails/*" alias (see vite.config.ts /
 * tsconfig.app.json), which points at ./wailsjs — the bindings Wails
 * generates from app/services/*.go once `wails dev` or
 * `wails generate module` has actually run.
 *
 * DELETE THIS FILE once wailsjs/ exists for real: the generated bindings
 * ship their own accurate .d.ts files, and this fallback would otherwise
 * silently mask type errors in them.
 */
declare module '@wails/*' {
  const value: unknown
  export default value
}
