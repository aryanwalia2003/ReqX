import { toResult } from '@/lib/result'
import { Open, OpenEnvironment, PickFile, Run } from '@wails/go/services/CollectionService'

import type { RunCollectionInput } from './types'

/** Native OS file picker kholta, .json filtered — user cancel kare to empty string milta, error nahi. */
export function pickFile(title: string) {
  return toResult(PickFile(title))
}

/** Collection file parse karke uske requests laata — CLI ke `reqx run` jaisa hi read path. */
export function openCollection(path: string) {
  return toResult(Open(path))
}

/** Environment file parse karta — CLI ke `-e env.json` flag jaisa hi read path. */
export function openEnvironment(path: string) {
  return toResult(OpenEnvironment(path))
}

/** Poori collection run karta — same CollectionRunner/metrics.Analyze pipeline jo CLI use karta. */
export function runCollection(input: RunCollectionInput) {
  return toResult(Run(input))
}
