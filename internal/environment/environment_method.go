package environment

// Set adds or updates a variable.
func (e *Environment) Set(key string, value string) { //enviroment ka method hai, key-values string leta hai
	if e.Variables == nil {
		e.Variables = make(map[string]string)
	}
	e.Variables[key] = value
}

// Get retrieves a variable by key.
func (e *Environment) Get(key string) (string, bool) {
	if e.Variables == nil {
		return "", false
	}
	val, ok := e.Variables[key]
	return val, ok
}

// Merge copies variables from another Environment into this one.
func (e *Environment) Merge(other *Environment) {
	if other == nil || other.Variables == nil {
		return
	}
	for k, v := range other.Variables {
		e.Set(k, v)
	}
}
