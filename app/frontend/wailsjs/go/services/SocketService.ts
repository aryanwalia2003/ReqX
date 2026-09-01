// Hand-written stand-in for what `wails generate module` would generate —
// see wailsjs/go/services/RequestService.ts for why this file exists and
// when to delete it.
import { callWailsMethod } from '@wails/go/wailsRuntime'

import type { ConnectSocketInput } from '@/features/socket-debugger/types'

export function Connect(input: ConnectSocketInput): Promise<void> {
  return callWailsMethod(window.go?.services?.SocketService?.Connect(input))
}

export function Send(text: string): Promise<void> {
  return callWailsMethod(window.go?.services?.SocketService?.Send(text))
}

export function Emit(eventName: string, payload: string): Promise<void> {
  return callWailsMethod(window.go?.services?.SocketService?.Emit(eventName, payload))
}

export function Disconnect(): Promise<void> {
  return callWailsMethod(window.go?.services?.SocketService?.Disconnect())
}

export function IsConnected(): Promise<boolean> {
  return callWailsMethod(window.go?.services?.SocketService?.IsConnected())
}
