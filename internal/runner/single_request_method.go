package runner

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"reqx/internal/collection"
	"reqx/internal/http_executor"
)

// SingleResult is the full outcome of one HTTP request, body and headers
// included. RequestMetric (the Run/runLinear path) deliberately drops the
// body to keep memory flat across a load test's thousands of iterations —
// ExecuteSingle exists for the opposite case: one request, where the caller
// actually wants to see what came back.
type SingleResult struct {
	StatusCode    int
	Status        string
	Headers       map[string]string
	Body          string
	Duration      time.Duration
	BytesSent     int64
	BytesReceived int64
}

// ExecuteSingle runs one HTTP request — var substitution, auth, execution —
// and returns the full response instead of a load-test RequestMetric. It
// reuses the same executor, auth resolution, and var substitution as the
// Run/runLinear path, so behavior (cookie jar, pooled transport, {{var}}
// substitution, auth types) stays identical between the CLI and this path.
func (cr *CollectionRunner) ExecuteSingle(req collection.Request, collAuth *collection.Auth, ctx *RuntimeContext) (*SingleResult, error) {
	urlStr := cr.replaceVars(req.URL, ctx)
	reqBody := cr.replaceVars(req.Body, ctx)

	var bodyReader io.Reader
	if reqBody != "" {
		bodyReader = bytes.NewBufferString(reqBody)
	}

	httpReq, err := http.NewRequest(strings.ToUpper(req.Method), urlStr, bodyReader)
	if err != nil {
		return nil, err
	}
	for k, v := range req.Headers {
		httpReq.Header.Set(k, cr.replaceVars(v, ctx))
	}
	http_executor.ApplyAuth(httpReq, cr.resolveAuth(req.Auth, collAuth, ctx))

	start := time.Now()
	resp, err := cr.executor.Execute(httpReq)
	duration := time.Since(start)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf := acquireBodyBuf()
	bytesReceived, err := io.Copy(buf, resp.Body)
	if err != nil {
		releaseBodyBuf(buf)
		return nil, err
	}
	bodyString := buf.String()
	releaseBodyBuf(buf)

	headers := make(map[string]string, len(resp.Header))
	for k, v := range resp.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	return &SingleResult{
		StatusCode:    resp.StatusCode,
		Status:        resp.Status,
		Headers:       headers,
		Body:          bodyString,
		Duration:      duration,
		BytesSent:     int64(len(reqBody)),
		BytesReceived: bytesReceived,
	}, nil
}
