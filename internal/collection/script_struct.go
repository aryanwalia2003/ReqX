package collection

// Script holds JavaScript code to execute before or after a request.
type Script struct {
	Type string   `json:"type"` // "prerequest" or "test"
	Exec []string `json:"exec"` // Array of strings representing lines of code
}

