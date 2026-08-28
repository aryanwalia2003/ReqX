package services

import (
	"strings"

	"reqx/internal/collection"
	"reqx/internal/errs"
	"reqx/internal/runner"
)

// SendRequestInput is the desktop app's minimal single-request contract —
// deliberately not the full collection.Request (that also carries DAG/script/
// socket fields no ad-hoc send needs). Method/Auth mirror internal/collection
// so the same request, once saved into a collection, needs no reshaping.
type SendRequestInput struct {
	Method       string            `json:"method"`
	URL          string            `json:"url"`
	Headers      map[string]string `json:"headers,omitempty"`
	Body         string            `json:"body,omitempty"`
	Auth         *collection.Auth  `json:"auth,omitempty"`
	EnvVariables map[string]string `json:"envVariables,omitempty"`
}

// SendRequestOutput is the full response — status, headers, body, timing.
type SendRequestOutput struct {
	StatusCode    int               `json:"statusCode"`
	Status        string            `json:"status"`
	Headers       map[string]string `json:"headers"`
	Body          string            `json:"body"`
	DurationMs    int64             `json:"durationMs"`
	BytesSent     int64             `json:"bytesSent"`
	BytesReceived int64             `json:"bytesReceived"`
}

// Send builds and executes one HTTP request through the shared
// CollectionRunner (pooled connections + cookie jar persist across calls),
// mirroring the CLI's `reqx req` path.
func (s *RequestService) Send(input SendRequestInput) (SendRequestOutput, error) {
	if strings.TrimSpace(input.URL) == "" {
		return SendRequestOutput{}, errs.InvalidInput("url is required")
	}
	method := strings.ToUpper(strings.TrimSpace(input.Method))
	if method == "" {
		method = "GET"
	}

	ctx := runner.NewRuntimeContext()
	for k, v := range input.EnvVariables {
		ctx.Environment.Set(k, v)
	}

	req := collection.Request{
		Method:  method,
		URL:     input.URL,
		Headers: input.Headers,
		Body:    input.Body,
		Auth:    input.Auth,
	}

	result, err := s.runner.ExecuteSingle(req, nil, ctx)
	if err != nil {
		return SendRequestOutput{}, errs.Wrap(err, errs.KindExternal, "request failed")
	}

	return SendRequestOutput{
		StatusCode:    result.StatusCode,
		Status:        result.Status,
		Headers:       result.Headers,
		Body:          result.Body,
		DurationMs:    result.Duration.Milliseconds(),
		BytesSent:     result.BytesSent,
		BytesReceived: result.BytesReceived,
	}, nil
}
