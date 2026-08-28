// Hand-written stand-in for what `wails generate module` would generate —
// see wailsjs/go/services/RequestService.ts for why this file exists and
// when to delete it.
import { callWailsMethod } from '@wails/go/wailsRuntime'

import type {
  Collection,
  Environment,
  RunCollectionInput,
  RunCollectionOutput,
} from '@/features/collection-runner/types'

export function PickFile(title: string): Promise<string> {
  return callWailsMethod(window.go?.services?.CollectionService?.PickFile(title))
}

export function Open(path: string): Promise<Collection> {
  return callWailsMethod(window.go?.services?.CollectionService?.Open(path))
}

export function OpenEnvironment(path: string): Promise<Environment> {
  return callWailsMethod(window.go?.services?.CollectionService?.OpenEnvironment(path))
}

export function Run(input: RunCollectionInput): Promise<RunCollectionOutput> {
  return callWailsMethod(window.go?.services?.CollectionService?.Run(input))
}
