import { toResult } from '@/lib/result'
import { Connect, Disconnect, Emit, IsConnected, Send } from '@wails/go/services/SocketService'

import type { ConnectSocketInput } from './types'

/** WS/Socket.IO connection kholta — CLI ke `reqx ws`/`reqx sio` jaisa hi dial path. */
export function connectSocket(input: ConnectSocketInput) {
  return toResult(Connect(input))
}

/** Raw text frame bhejta (plain WebSocket send). */
export function sendSocketMessage(text: string) {
  return toResult(Send(text))
}

/** Socket.IO event emit karta — `42[name, payload]` frame. */
export function emitSocketEvent(eventName: string, payload: string) {
  return toResult(Emit(eventName, payload))
}

export function disconnectSocket() {
  return toResult(Disconnect())
}

export function isSocketConnected() {
  return toResult(IsConnected())
}
