package runner

import (
	"strings"
	"sync"
)

// builderPool recycles strings.Builder instances used during variable
// substitution. This avoids allocating a new builder on every call.
var builderPool = sync.Pool{
    New: func() any { return new(strings.Builder) },
}

// replaceVarsFast substitutes all {{key}} placeholders in template with
// values from vars in a single linear scan — O(1) allocations regardless
// of how many variables exist.
//
// It replaces the old loop-of-ReplaceAll pattern which allocated one new
// string per variable.
func replaceVarsFast(template string, vars map[string]string) string {
    // Fast path: nothing to substitute
    if !strings.Contains(template, "{{") {
        return template
    }

    sb := builderPool.Get().(*strings.Builder)
    defer func() {
        sb.Reset()
        builderPool.Put(sb)
    }()

    i := 0
    for i < len(template) {
        // Find the next opening brace
        open := strings.Index(template[i:], "{{")
        if open == -1 {
            // No more placeholders — write the rest and stop
            sb.WriteString(template[i:])
            break
        }

        // Write everything before the {{
        sb.WriteString(template[i : i+open])

        // Find the closing brace from the {{ position
        rest := template[i+open:]
        close := strings.Index(rest, "}}")
        if close == -1 {
            // Unclosed {{ — treat as literal, write and stop
            sb.WriteString(rest)
            break
        }

        key := rest[2:close] // text between {{ and }}

        if val, ok := vars[key]; ok {
            sb.WriteString(val)
        } else {
            // Key not found — preserve the original {{key}}
            sb.WriteString("{{")
            sb.WriteString(key)
            sb.WriteString("}}")
        }

        i = i + open + close + 2 // advance past }}
    }

    return sb.String()
}