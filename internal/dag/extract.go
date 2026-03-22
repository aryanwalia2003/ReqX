package dag

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// ExtractAll evaluates every JSONPath expression in paths against body and
// returns a map of variable-name → extracted string value.
//
// body must be a valid JSON byte slice. Paths that fail to resolve are
// silently skipped and returned in the errors slice so callers can log them
// without aborting the extraction of other keys.
//
// No external dependencies are used. The body is decoded once into a generic
// interface{} tree; individual paths are then walked in-memory.
func ExtractAll(body []byte, paths map[string]string) (results map[string]string, errs []error) {
	if len(paths) == 0 || len(body) == 0 {
		return nil, nil
	}

	// Decode once; reuse the tree for every path.
	var root interface{}
	if err := json.Unmarshal(body, &root); err != nil {
		return nil, []error{fmt.Errorf("extract: failed to parse response body as JSON: %w", err)}
	}

	results = make(map[string]string, len(paths))
	for varName, path := range paths {
		val, err := evalPath(root, path)
		if err != nil {
			errs = append(errs, fmt.Errorf("extract: %q → %q: %w", varName, path, err))
			continue
		}
		results[varName] = val
	}
	return results, errs
}

// evalPath walks root according to a simplified JSONPath expression and returns
// the resolved value as a string.
//
// Grammar supported:
//
//	path   = "$" { "." key | "[" index "]" }
//	key    = identifier
//	index  = non-negative integer
//
// Examples:
//
//	"$.token"           → root["token"]
//	"$.data.user.id"    → root["data"]["user"]["id"]
//	"$.items[0].name"   → root["items"][0]["name"]
//	"$.results[2]"      → root["results"][2]
func evalPath(root interface{}, path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("empty path")
	}
	if !strings.HasPrefix(path, "$") {
		return "", fmt.Errorf("path must start with '$', got %q", path)
	}

	// Strip leading "$" and tokenize the remainder.
	tokens, err := tokenize(path[1:])
	if err != nil {
		return "", err
	}

	cur := root
	for _, tok := range tokens {
		cur, err = step(cur, tok)
		if err != nil {
			return "", err
		}
	}

	return stringify(cur)
}

// token is a single navigation step: either a map key or an array index.
type token struct {
	key   string
	index int
	isIdx bool
}

// tokenize splits a path suffix (the part after "$") into navigation tokens.
// ".foo"     → {key:"foo"}
// "[2]"      → {isIdx:true, index:2}
// ".foo[1]"  → {key:"foo"}, {isIdx:true, index:1}
func tokenize(path string) ([]token, error) {
	var tokens []token

	for len(path) > 0 {
		switch path[0] {
		case '.':
			// key segment: consume until '.' or '[' or end
			path = path[1:]
			end := strings.IndexAny(path, ".[")
			if end == -1 {
				end = len(path)
			}
			if end == 0 {
				return nil, fmt.Errorf("empty key after '.'")
			}
			tokens = append(tokens, token{key: path[:end]})
			path = path[end:]

		case '[':
			// index segment: consume digits until ']'
			close := strings.IndexByte(path, ']')
			if close == -1 {
				return nil, fmt.Errorf("unclosed '[' in path")
			}
			idxStr := path[1:close]
			idx, err := strconv.Atoi(idxStr)
			if err != nil || idx < 0 {
				return nil, fmt.Errorf("invalid array index %q", idxStr)
			}
			tokens = append(tokens, token{isIdx: true, index: idx})
			path = path[close+1:]

		default:
			return nil, fmt.Errorf("unexpected character %q in path", string(path[0]))
		}
	}

	return tokens, nil
}

// step advances the cursor by one token.
func step(cur interface{}, tok token) (interface{}, error) {
	if tok.isIdx {
		arr, ok := cur.([]interface{})
		if !ok {
			return nil, fmt.Errorf("expected array, got %T", cur)
		}
		if tok.index >= len(arr) {
			return nil, fmt.Errorf("index %d out of range (array len %d)", tok.index, len(arr))
		}
		return arr[tok.index], nil
	}

	// key step
	obj, ok := cur.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected object, got %T", cur)
	}
	val, exists := obj[tok.key]
	if !exists {
		return nil, fmt.Errorf("key %q not found", tok.key)
	}
	return val, nil
}

// stringify converts any JSON-decoded value to its string representation.
// Strings pass through directly. Numbers, booleans, and null are formatted.
// Objects and arrays are JSON-encoded.
func stringify(v interface{}) (string, error) {
	switch t := v.(type) {
	case string:
		return t, nil
	case float64:
		// json.Unmarshal decodes all numbers as float64.
		// If it is a whole number, omit the decimal point for cleaner env vars.
		if t == float64(int64(t)) {
			return strconv.FormatInt(int64(t), 10), nil
		}
		return strconv.FormatFloat(t, 'f', -1, 64), nil
	case bool:
		return strconv.FormatBool(t), nil
	case nil:
		return "", nil
	default:
		// Arrays and nested objects — encode back to JSON string.
		b, err := json.Marshal(v)
		if err != nil {
			return "", fmt.Errorf("could not stringify value of type %T: %w", v, err)
		}
		return string(b), nil
	}
}