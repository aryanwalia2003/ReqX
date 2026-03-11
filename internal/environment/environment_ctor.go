package environment

// NewEnvironment creates an empty Environment. (Constructor hai)
func NewEnvironment(name string) *Environment {
	return &Environment{
		Name:      name,
		Variables: make(map[string]string),
	}
}
