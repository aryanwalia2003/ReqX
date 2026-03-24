package runner

import "reqx/internal/environment"

// SetGlobalVariable sets a variable in the global context.
func (rc *RuntimeContext) SetGlobalVariable(key string, value interface{}) {
	if rc.GlobalVariables == nil {
		rc.GlobalVariables = make(map[string]interface{})
	}
	rc.GlobalVariables[key] = value
}

// GetVariable attempts to find a variable, checking globals first, then environment.
func (rc *RuntimeContext) GetVariable(key string) (interface{}, bool) {
	// Check global variables
	if rc.GlobalVariables != nil {
		if val, exists := rc.GlobalVariables[key]; exists {
			return val, true
		}
	}

	// Check environment variables
	if rc.Environment != nil {
		if val, exists := rc.Environment.Get(key); exists {
			return val, true
		}
	}

	return nil, false
}

// SetEnvironment configures the active environment for this context.
func (rc *RuntimeContext) SetEnvironment(env *environment.Environment) {
	rc.Environment = env
}

// IsConnected reports whether an async socket URL is already connected for
// this worker. Safe for concurrent use.
func (rc *RuntimeContext) IsConnected(url string) bool {
	rc.connMu.Lock()
	defer rc.connMu.Unlock()
	if rc.connectedURLs == nil {
		return false
	}
	_, ok := rc.connectedURLs[url]
	return ok
}

// MarkConnected records that an async socket URL has been dialled for this
// worker so subsequent iterations can skip the reconnect. Safe for concurrent use.
func (rc *RuntimeContext) MarkConnected(url string) {
	rc.connMu.Lock()
	defer rc.connMu.Unlock()
	if rc.connectedURLs == nil {
		rc.connectedURLs = make(map[string]struct{})
	}
	rc.connectedURLs[url] = struct{}{}
}
