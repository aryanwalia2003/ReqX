package dag

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

// ExtractAll evaluates every JSONPath expression in paths against body and
// returns a map of variable-name → extracted string value.
//
// Internal optimization: Uses github.com/tidwall/gjson to parse values 
// directly from the byte slice without unmarshalling into maps.
func ExtractAll(body []byte, paths map[string]string) (results map[string]string, errs []error) {
	if len(paths) == 0 || len(body) == 0 {
		return nil, nil
	}

	results = make(map[string]string, len(paths))
	for varName, path := range paths {
		cleanPath := sanitizePath(path)
		
		res := gjson.GetBytes(body, cleanPath)
		if !res.Exists() {
			errs = append(errs, fmt.Errorf("extract: %q → %q: path not found", varName, path))
			continue
		}

		// gjson.Result.String() returns the raw string for strings, 
		// "true"/"false" for bools, and the raw JSON for objects/arrays.
		results[varName] = res.String()
	}
	return results, errs
}

// sanitizePath converts standard JSONPath ($.foo.bar[0]) to gjson syntax (foo.bar.0).
func sanitizePath(path string) string {
	path = strings.TrimSpace(path)
	if path == "$" {
		return ""
	}
	
	// Remove leading $
	if strings.HasPrefix(path, "$.") {
		path = path[2:]
	} else if strings.HasPrefix(path, "$") {
		path = path[1:]
	}

	// Convert [index] to .index
	// Note: Simple implementation for common cases. 
	// For example: items[0].name -> items.0.name
	path = strings.ReplaceAll(path, "[", ".")
	path = strings.ReplaceAll(path, "]", "")

	return path
}