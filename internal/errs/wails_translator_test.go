package errs

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestFormatForWails(t *testing.T) {
	cases := []struct {
		name     string
		err      error
		wantKind string
		wantMsg  string
	}{
		{
			name:     "app error passes through kind and message",
			err:      InvalidInput("name is required"),
			wantKind: string(KindInvalidInput),
			wantMsg:  "name is required",
		},
		{
			name:     "bare error becomes internal with a safe message",
			err:      errors.New("sql: no rows in result set"),
			wantKind: string(KindInternal),
			wantMsg:  "An unexpected error occurred.",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := FormatForWails(tc.err).(string)
			if !ok {
				t.Fatalf("FormatForWails did not return a string: %v", got)
			}

			var resp WailsErrorResponse
			if err := json.Unmarshal([]byte(got), &resp); err != nil {
				t.Fatalf("FormatForWails output isn't valid JSON: %v", err)
			}
			if resp.Kind != tc.wantKind {
				t.Errorf("kind = %q, want %q", resp.Kind, tc.wantKind)
			}
			if resp.Message != tc.wantMsg {
				t.Errorf("message = %q, want %q", resp.Message, tc.wantMsg)
			}
		})
	}
}
