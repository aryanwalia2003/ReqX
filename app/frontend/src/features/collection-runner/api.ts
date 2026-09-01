import { toResult } from '@/lib/result'
import {
  Open,
  OpenEnvironment,
  OpenInEditor,
  OpenPersonas,
  PickFile,
  PickSaveFile,
  Run,
  Save,
} from '@wails/go/services/CollectionService'

import type { Collection, RunCollectionInput } from './types'

/** Native OS file picker kholta, *.<extension> filtered — user cancel kare to empty string milta, error nahi. */
export function pickFile(title: string, extension: string) {
  return toResult(PickFile(title, extension))
}

/** Native OS save dialog — Export results (--export) ke liye. */
export function pickSaveFile(title: string, defaultFilename: string) {
  return toResult(PickSaveFile(title, defaultFilename))
}

/** Collection file parse karke uske requests laata — CLI ke `reqx run` jaisa hi read path. */
export function openCollection(path: string) {
  return toResult(Open(path))
}

/** Collection ko JSON file me likhta — CLI ke `saveCollection` jaisa hi write path. */
export function saveCollection(collection: Collection, path: string) {
  return toResult(Save(collection, path))
}

/** File ko VS Code/Vim/OS default text editor me kholta. */
export function openInEditor(kind: string, path: string) {
  return toResult(OpenInEditor(kind, path))
}

/** Environment file parse karta — CLI ke `-e env.json` flag jaisa hi read path. */
export function openEnvironment(path: string) {
  return toResult(OpenEnvironment(path))
}

/** Personas CSV parse karta — CLI ke `--personas` flag jaisa hi read path. */
export function openPersonas(path: string) {
  return toResult(OpenPersonas(path))
}

/** Poori collection run karta — same CollectionRunner/metrics.Analyze pipeline jo CLI use karta. */
export function runCollection(input: RunCollectionInput) {
  return toResult(Run(input))
}
