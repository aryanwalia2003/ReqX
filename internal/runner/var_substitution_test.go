package runner

import (
	"strings"
	"testing"
)

func TestReplaceVarsFast(t *testing.T) {
    vars := map[string]string{
        "base_url": "https://api.example.com",
        "token":    "abc123",
        "user_id":  "42",
    }

    cases := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "single var",
            input:    "{{base_url}}/users",
            expected: "https://api.example.com/users",
        },
        {
            name:     "multiple vars",
            input:    "{{base_url}}/users/{{user_id}}",
            expected: "https://api.example.com/users/42",
        },
        {
            name:     "var in header value",
            input:    "Bearer {{token}}",
            expected: "Bearer abc123",
        },
        {
            name:     "missing key preserved",
            input:    "{{base_url}}/{{unknown}}",
            expected: "https://api.example.com/{{unknown}}",
        },
        {
            name:     "no placeholders — fast path",
            input:    "https://hardcoded.com/path",
            expected: "https://hardcoded.com/path",
        },
        {
            name:     "empty string",
            input:    "",
            expected: "",
        },
        {
            name:     "unclosed brace",
            input:    "{{base_url}}/{{unclosed",
            expected: "https://api.example.com/{{unclosed",
        },
        {
            name:     "var at start and end",
            input:    "{{base_url}}{{token}}",
            expected: "https://api.example.com" + "abc123",
        },
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            got := replaceVarsFast(tc.input, vars)
            if got != tc.expected {
                t.Errorf("\ninput:    %q\nexpected: %q\ngot:      %q", tc.input, tc.expected, got)
            }
        })
    }
}

func BenchmarkReplaceVarsOld(b *testing.B) {
    vars := map[string]string{
        "base_url": "https://api.example.com",
        "token":    "abc123",
        "user_id":  "42",
        "env":      "production",
        "version":  "v2",
    }
    template := "{{base_url}}/api/{{version}}/users/{{user_id}}?env={{env}}&token={{token}}"

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        out := template
        for k, v := range vars {
            out = strings.ReplaceAll(out, "{{"+k+"}}", v)
        }
        _ = out
    }
}

func BenchmarkReplaceVarsFast(b *testing.B) {
    vars := map[string]string{
        "base_url": "https://api.example.com",
        "token":    "abc123",
        "user_id":  "42",
        "env":      "production",
        "version":  "v2",
    }
    template := "{{base_url}}/api/{{version}}/users/{{user_id}}?env={{env}}&token={{token}}"

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = replaceVarsFast(template, vars)
    }
}