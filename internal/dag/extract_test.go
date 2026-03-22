package dag

import (
	"encoding/json"
	"testing"
)

func TestEvalPath(t *testing.T) {
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

	var root interface{}
	if err := json.Unmarshal(raw, &root); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		path    string
		want    string
		wantErr bool
	}{
		{"$.token", "abc123", false},
		{"$.data.user.id", "42", false},
		{"$.data.user.name", "Aryan", false},
		{"$.data.score", "9.5", false},
		{"$.items[0]", "first", false},
		{"$.items[2]", "third", false},
		{"$.nested[1].val", "y", false},
		{"$.nested[0].id", "1", false},
		{"$.flag", "true", false},
		{"$.empty", "", false},
		// error cases
		{"$.missing", "", true},
		{"$.items[99]", "", true},
		{"$.data.user.id.bad", "", true},
		{"no_dollar", "", true},
	}

	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			got, err := evalPath(root, tc.path)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error for %q, got %q", tc.path, got)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error for %q: %v", tc.path, err)
				return
			}
			if got != tc.want {
				t.Errorf("%q: want %q got %q", tc.path, tc.want, got)
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
	if len(errs) != 1 {
		t.Errorf("expected 1 parse error, got %d", len(errs))
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