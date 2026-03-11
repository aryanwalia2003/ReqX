package storage

import (
	"os"
)

// ReadJSONFile reads a file and returns its raw byte content.
func ReadJSONFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteJSONFile writes raw bytes to a file.
func WriteJSONFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}
