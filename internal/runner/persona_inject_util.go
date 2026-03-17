package runner

import (
	"strings"

	"reqx/internal/personas"
)

func applyPersona(ctx *RuntimeContext, p personas.Persona) {
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

