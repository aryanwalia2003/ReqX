package runner

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"postman-cli/internal/collection"
	"postman-cli/internal/http_executor"
	"postman-cli/internal/scripting"
)

// SetClearCookiesPerRequest controls whether the cookie jar is cleared before each request.
func (cr *CollectionRunner) SetClearCookiesPerRequest(v bool) {
	cr.clearCookiesPerRequest = v
}

// Run executes all requests within a collection sequentially.
func (cr *CollectionRunner) Run(coll *collection.Collection, ctx *RuntimeContext) error {
	for _, req := range coll.Requests {
		fmt.Printf("Running request: %s\n", req.Name)

		// 1. Pre-request Scripts
		cr.runScripts("prerequest", req.Scripts, ctx, nil)

		// 2. Variable replacement
		urlStr := cr.replaceVars(req.URL, ctx)

		// Check Protocol
		if strings.ToUpper(req.Protocol) == "SOCKETIO" {
			headers := make(map[string]string)
			for k, v := range req.Headers {
				headers[k] = cr.replaceVars(v, ctx)
			}

			err := cr.sioExecutor.Execute(urlStr, headers, req.Events)
			if err != nil {
				fmt.Printf("Socket.IO Request %s failed: %v\n", req.Name, err)
			} else {
				fmt.Printf("Socket.IO Execution %s completed.\n", req.Name)
			}

			// 5. Test Scripts
			cr.runScripts("test", req.Scripts, ctx, nil)
			continue
		}

		// 3. Build HTTP request
		var bodyReader io.Reader
		if req.Body != "" {
			bodyBytes := []byte(cr.replaceVars(req.Body, ctx))
			bodyReader = bytes.NewBuffer(bodyBytes)
		}

		httpReq, err := http.NewRequest(strings.ToUpper(req.Method), urlStr, bodyReader)
		if err != nil {
			fmt.Printf("Failed to create request %s: %v\n", req.Name, err)
			continue
		}

		for k, v := range req.Headers {
			httpReq.Header.Set(k, cr.replaceVars(v, ctx))
		}

		// 3a. Apply auth — request-level overrides collection-level
		effectiveAuth := cr.resolveAuth(req.Auth, coll.Auth, ctx)
		http_executor.ApplyAuth(httpReq, effectiveAuth)

		// 3b. Optionally clear cookies before each request
		if cr.clearCookiesPerRequest {
			cr.executor.ClearCookies()
		}

		// 4. Execute HTTP request
		resp, err := cr.executor.Execute(httpReq)
		if err != nil {
			fmt.Printf("Request %s failed: %v\n", req.Name, err)
			continue
		}

		fmt.Printf("Status: %s\n", resp.Status)

		// Capture body for script access
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		// Prepare scripting response object
		scriptResp := &scripting.ResponseAPI{
			BodyString: string(bodyBytes),
			HeadersMap: make(map[string]string),
			Headers:    &scripting.ResponseHeaders{Headers: make(map[string]string)},
		}

		for k, v := range resp.Header {
			if len(v) > 0 {
				scriptResp.HeadersMap[k] = v[0]
				scriptResp.Headers.Headers[k] = v[0]
			}
		}

		// 5. Test Scripts
		cr.runScripts("test", req.Scripts, ctx, scriptResp)
	}

	return nil
}

// resolveAuth returns the effective auth, applying variable replacement.
// Precedence: request.Auth > collection.Auth > nil.
func (cr *CollectionRunner) resolveAuth(reqAuth, collAuth *collection.Auth, ctx *RuntimeContext) *collection.Auth {
	src := reqAuth
	if src == nil {
		src = collAuth
	}
	if src == nil {
		return nil
	}
	// Deep-copy with vars replaced so original struct is not mutated.
	resolved := &collection.Auth{
		Type:     src.Type,
		Token:    cr.replaceVars(src.Token, ctx),
		Username: cr.replaceVars(src.Username, ctx),
		Password: cr.replaceVars(src.Password, ctx),
		Key:      cr.replaceVars(src.Key, ctx),
		Value:    cr.replaceVars(src.Value, ctx),
		In:       src.In,
	}
	if src.Cookies != nil {
		resolved.Cookies = make(map[string]string, len(src.Cookies))
		for k, v := range src.Cookies {
			resolved.Cookies[k] = cr.replaceVars(v, ctx)
		}
	}
	return resolved
}

func (cr *CollectionRunner) runScripts(scriptType string, scripts []collection.Script, ctx *RuntimeContext, resp *scripting.ResponseAPI) {
	for _, s := range scripts {
		if s.Type == scriptType {
			err := cr.scriptRunner.Execute(&s, ctx.Environment, resp)
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
	out := input
	for k, v := range ctx.Environment.Variables {
		out = strings.ReplaceAll(out, "{{"+k+"}}", v)
	}
	return out
}


