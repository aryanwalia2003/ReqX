export type SocketProtocol = 'ws' | 'sio'

export interface ConnectSocketInput {
  url: string
  protocol: SocketProtocol
  headers?: Record<string, string>
}

export interface SocketMessage {
  direction: 'in' | 'out' | 'system'
  data: string
  eventName?: string
  timestamp: number
}
