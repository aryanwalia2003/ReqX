package runner

import (
	"reqx/internal/http_executor"
	"reqx/internal/scripting"
	"reqx/internal/socketio_executor"
	"reqx/internal/websocket_executor"
	"time"
)

type RequestMetric struct {
	Name         string
	Protocol     string
	StatusCode   int
	Duration     time.Duration
	StatusString string
	Error        error
	ErrorMsg     string
	GraphQLError string // non-empty when a 200 OK response body carried a GraphQL "errors" array
	WorkerID     int

	BytesSent     int64
	BytesReceived int64
	TTFB          time.Duration
}

// CollectionRunner handles executing a full collection of requests.
type CollectionRunner struct {
	executor               *http_executor.DefaultExecutor
	sioExecutor            socketio_executor.SocketIOExecutor
	weExecutor             websocket_executor.WebSocketExecutor
	scriptRunner           scripting.ScriptRunner
	clearCookiesPerRequest bool // if true, jar is cleared before each request
	graphqlErrorCheck      bool // --graphql: parse response body for a GraphQL "errors" array
	verbosity              int
}
