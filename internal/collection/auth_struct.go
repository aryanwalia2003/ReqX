package collection

// Auth defines authentication configuration for a request or collection.
// Type controls which fields are used:
//   - "bearer"  → Token
//   - "basic"   → Username, Password
//   - "apikey"  → Key, Value, In ("header"|"query")
//   - "cookie"  → Cookies map
//   - "none"    → no auth applied
type Auth struct {
	Type     string            `json:"type"`
	Token    string            `json:"token,omitempty"`
	Username string            `json:"username,omitempty"`
	Password string            `json:"password,omitempty"`
	Key      string            `json:"key,omitempty"`
	Value    string            `json:"value,omitempty"`
	In       string            `json:"in,omitempty"`
	Cookies  map[string]string `json:"cookies,omitempty"`
}
