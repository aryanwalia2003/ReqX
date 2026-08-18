package runner

import "github.com/tidwall/gjson"

// detectGraphQLError reports a non-empty message when body carries a
// top-level GraphQL "errors" array — a 200 OK does not mean success for
// GraphQL, so status-code-only pass/fail (the previous behaviour) silently
// counted these as passing requests.
func detectGraphQLError(body []byte) string {
	if len(body) == 0 {
		return ""
	}
	errs := gjson.GetBytes(body, "errors")
	if !errs.IsArray() {
		return ""
	}
	arr := errs.Array()
	if len(arr) == 0 {
		return ""
	}
	msg := arr[0].Get("message").String()
	if msg == "" {
		return "graphql error (no message)"
	}
	if len(arr) > 1 {
		return msg + " (+ more)"
	}
	return msg
}
