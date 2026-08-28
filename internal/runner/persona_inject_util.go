package runner

import (
	"strings"

	"reqx/internal/personas"
)

// ApplyPersona sets ctx.Environment.Variables["persona.<col>"] for every
// column in p, so a request can reference {{persona.<col>}}. Exported so
// callers outside package runner (the CLI's sequential path in
// cmd/run_cmd_ctor.go, app/services.CollectionService) don't each
// reimplement it.
func ApplyPersona(ctx *RuntimeContext, p personas.Persona) {
	if ctx == nil {
		return
	}
	if ctx.Environment == nil {
		return
	}
	if ctx.Environment.Variables == nil {
		ctx.Environment.Variables = make(map[string]string)
	}
	for k, v := range p {
		key := strings.TrimSpace(k)
		if key == "" {
			continue
		}
		ctx.Environment.Variables["persona."+key] = v
	}
}
