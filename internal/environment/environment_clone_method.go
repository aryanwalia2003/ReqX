package environment

// Clone creates a deep copy of the environment to ensure virtual users do not share state.
func (e *Environment) Clone() *Environment {
	if e == nil {
		return nil
	}

	clone := &Environment{
		Name:      e.Name,
		Variables: make(map[string]string, len(e.Variables)),
	}

	for k, v := range e.Variables {
		clone.Variables[k] = v
	}

	return clone
}