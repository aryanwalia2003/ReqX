package dag

import (
	"testing"
)

func TestExtractAll(t *testing.T) {
	raw := []byte(`{
		"token": "abc123",
		"data": {
			"user": {"id": 42, "name": "Aryan"},
			"score": 9.5
		},
		"items": ["first", "second", "third"],
		"nested": [{"id": 1, "val": "x"}, {"id": 2, "val": "y"}],
		"flag": true,
		"empty": null
	}`)

	cases := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{"Top level string", "$.token", "abc123", false},
		{"Nested object int", "$.data.user.id", "42", false},
		{"Nested object string", "$.data.user.name", "Aryan", false},
		{"Nested object float", "$.data.score", "9.5", false},
		{"Array index 0", "$.items[0]", "first", false},
		{"Array index 2", "$.items[2]", "third", false},
		{"Nested array object", "$.nested[1].val", "y", false},
		{"Nested array index", "$.nested[0].id", "1", false},
		{"Boolean true", "$.flag", "true", false},
		{"Null value", "$.empty", "", false},
		// error cases
		{"Missing key", "$.missing", "", true},
		{"Out of bounds", "$.items[99]", "", true},
		{"Path through leaf", "$.data.user.id.bad", "", true},
		{"No dollar prefix", "no_dollar", "", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			results, errs := ExtractAll(raw, map[string]string{"var": tc.path})
			
			if tc.wantErr {
				if len(errs) == 0 {
					t.Errorf("expected error for path %q, but got none", tc.path)
				}
				return
			}

			if len(errs) > 0 {
				t.Errorf("unexpected errors for path %q: %v", tc.path, errs)
				return
			}

			if results["var"] != tc.want {
				t.Errorf("path %q: want %q, got %q", tc.path, tc.want, results["var"])
			}
		})
	}
}

func TestExtractAll_HappyPath(t *testing.T) {
	body := []byte(`{"auth":{"token":"xyz"},"uid":7}`)
	paths := map[string]string{"my_token": "$.auth.token", "my_uid": "$.uid"}
	results, errs := ExtractAll(body, paths)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if results["my_token"] != "xyz" {
		t.Errorf("my_token: want xyz got %q", results["my_token"])
	}
	if results["my_uid"] != "7" {
		t.Errorf("my_uid: want 7 got %q", results["my_uid"])
	}
}

func TestExtractAll_PartialErrors(t *testing.T) {
	body := []byte(`{"auth":{"token":"xyz"},"uid":7}`)
	paths := map[string]string{"my_token": "$.auth.token", "bad": "$.does_not_exist"}
	results, errs := ExtractAll(body, paths)
	if results["my_token"] != "xyz" {
		t.Errorf("my_token: want xyz got %q", results["my_token"])
	}
	if _, ok := results["bad"]; ok {
		t.Error("bad path should not appear in results")
	}
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d: %v", len(errs), errs)
	}
}

func TestExtractAll_InvalidJSON(t *testing.T) {
	_, errs := ExtractAll([]byte(`not json`), map[string]string{"k": "$.k"})
	if len(errs) == 0 {
		t.Error("expected parse error for invalid JSON, got none")
	}
}

func TestExtractAll_EmptyInputs(t *testing.T) {
	r1, e1 := ExtractAll(nil, map[string]string{"k": "$.k"})
	if r1 != nil || e1 != nil {
		t.Error("nil body should return nil,nil")
	}
	r2, e2 := ExtractAll([]byte(`{"k":"v"}`), nil)
	if r2 != nil || e2 != nil {
		t.Error("nil paths should return nil,nil")
	}
}