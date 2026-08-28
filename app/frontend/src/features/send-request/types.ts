// Types local to the "send-request" feature.

/** Mirrors internal/collection.Auth (Go) — see app/services/request_service_method.go. */
export interface AuthConfig {
  type: 'bearer' | 'basic' | 'apikey' | 'cookie' | 'none'
  token?: string
  username?: string
  password?: string
  key?: string
  value?: string
  in?: 'header' | 'query'
  cookies?: Record<string, string>
}

/** Mirrors app/services.SendRequestInput (Go). */
export interface SendRequestInput {
  method: string
  url: string
  headers?: Record<string, string>
  body?: string
  auth?: AuthConfig
  envVariables?: Record<string, string>
}

/** Mirrors app/services.SendRequestOutput (Go). */
export interface SendRequestOutput {
  statusCode: number
  status: string
  headers: Record<string, string>
  body: string
  durationMs: number
  bytesSent: number
  bytesReceived: number
}
