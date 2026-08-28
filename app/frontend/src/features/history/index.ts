// Public API for the "history" feature. Only export what other features and
// src/app are allowed to import — everything else in this folder is private.
export { HistoryPanel } from './components/HistoryPanel'
export type { RunRow, StatRow } from './types'
