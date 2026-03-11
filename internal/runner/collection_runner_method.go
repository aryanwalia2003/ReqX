package runner

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"postman-cli/internal/collection"
)

// Run executes all requests within a collection sequentially.
func (cr *CollectionRunner) Run(coll *collection.Collection, ctx *RuntimeContext) error {
	for _, req := range coll.Requests {
		fmt.Printf("Running request: %s\n", req.Name)

		// 1. Pre-request Scripts
		cr.runScripts("prerequest", req.Scripts, ctx)

		// 2. Variable replacement (simple text replace for now)
		url := cr.replaceVars(req.URL, ctx)

		// 3. Build HTTP request
		var bodyReader io.Reader
		if req.Body != "" {
			bodyBytes := []byte(cr.replaceVars(req.Body, ctx))
			bodyReader = bytes.NewBuffer(bodyBytes)
		}

		httpReq, err := http.NewRequest(strings.ToUpper(req.Method), url, bodyReader)
		if err != nil {
			fmt.Printf("Failed to create request %s: %v\n", req.Name, err)
			continue
		}

		for k, v := range req.Headers {
			httpReq.Header.Set(k, cr.replaceVars(v, ctx))
		}

		// 4. Exec HTTP Request
		resp, err := cr.executor.Execute(httpReq)
		if err != nil {
			fmt.Printf("Request %s failed: %v\n", req.Name, err)
			continue
		}
		
		fmt.Printf("Status: %s\n", resp.Status)
		resp.Body.Close()

		// 5. Test Scripts
		cr.runScripts("test", req.Scripts, ctx)
	}

	return nil
}

func (cr *CollectionRunner) runScripts(scriptType string, scripts []collection.Script, ctx *RuntimeContext) {
	for _, s := range scripts {
		if s.Type == scriptType {
			err := cr.scriptRunner.Execute(&s, ctx.Environment)
			if err != nil {
				fmt.Printf("Warning: script execution failed: %v\n", err)
			}
		}
	}
}

func (cr *CollectionRunner) replaceVars(input string, ctx *RuntimeContext) string {
	if ctx == nil || ctx.Environment == nil {
		return input
	}
	// Very simple {{var}} naive replacement for MVP
	out := input
	for k, v := range ctx.Environment.Variables {
		out = strings.ReplaceAll(out, "{{"+k+"}}", v)
	}
	return out
}
