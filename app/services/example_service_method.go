package services

import "reqx/internal/errs"

// PingRequest/PingResponse show the request/response struct convention —
// every bound method takes and returns plain structs so Wails can generate
// matching TypeScript types for the frontend automatically. Don't bind
// methods with bare primitive params/returns; wrap them.
type PingRequest struct {
	Name string `json:"name"`
}

type PingResponse struct {
	Message string `json:"message"`
}

// Ping is a template bound method. Wails exposes Go's PascalCase method names
// to the frontend as camelCase — ExampleService.Ping becomes
// window.go.services.ExampleService.ping(...) in the generated JS bindings.
//
// Errors always go through internal/errs — the same error-kind convention the
// CLI already uses (see errs.Wrap in cmd/run_cmd_ctor.go). Never return a bare
// fmt.Errorf from a bound method; the frontend's error handling (src/lib/result.ts)
// expects a consistent shape.
func (s *ExampleService) Ping(req PingRequest) (PingResponse, error) {
	if req.Name == "" {
		return PingResponse{}, errs.InvalidInput("name is required")
	}
	return PingResponse{Message: "pong, " + req.Name}, nil
}
