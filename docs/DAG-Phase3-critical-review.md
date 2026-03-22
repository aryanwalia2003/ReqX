# DAG Phase 3 Critical Review: Data Extraction (`extract`)

## 🔴 CRITICAL: Potential Memory Corruption & Lifetime Fragility

In `internal/runner/collection_runner_method.go`:
```go
179:  if len(req.Extract) > 0 && resp.StatusCode < 400 {
180:      applyExtracts(req.Extract, bodyBytes, ctx, req.Name, cr.verbosity)
181:  }
...
247:  releaseBodyBuf(buf) // Buffer goes back to pool
```

### The Problem
`bodyBytes` is a direct slice (`buf.Bytes()`) of a pooled buffer. Standard `json.Unmarshal` in Go (used by `ExtractAll`) into an `interface{}` tree creates a representation where strings *may* hold direct references to the source byte slice depending on internal implementation or future optimizations (like using `json.Unmarshaler` on types). 

Currently, `releaseBodyBuf` is after `applyExtracts`, which is safe. **BUT**, if a future refactor moves `releaseBodyBuf` even one line up, the memory in `bodyBytes` will be overwritten by a new request. Any `interface{}` tree (the `root` variable) that hasn't finished its `evalPath` traversal will silently read garbage data.

### Recommendation
**Hard Order Constraint**: Add a critical comment in `collection_runner_method.go` or make a hard copy of `bodyBytes` if the overhead is acceptable to decouple the lifecycle.

---

## 🟡 MEDIUM: GC Pressure on Large Response Bodies

`json.Unmarshal` into an `interface{}` tree is memory-intensive. For a 1MB JSON response, it can allocate 2–4MB of short-lived nodes on the heap. While the `Wait()` prevents a data race, it doesn't prevent **memory spikes** if 50 parallel DAG nodes all receive 1MB responses at once.

### The Problem
There is no size guard in `applyExtracts`. If a backend accidentally sends a 10MB JSON dump, the CLI might OOM (Out Of Memory) trying to parse it for a single token.

### Recommendation
Set a reasonable threshold (e.g., 1MB) beyond which `extract` is skipped with a warning, or use a streaming decoder:
```go
if len(bodyBytes) > (1024 * 1024) {
    color.Yellow("⚠ [EXTRACT] %s: body too large for fast extraction (>1MB). Use scripts for large bodies.\n", reqName)
    return
}
```

---

## 🟡 MEDIUM: Integer Precision & Overflow in `stringify`

In `internal/dag/extract.go`:
```go
case float64:
    if t == float64(int64(t)) {
        return strconv.FormatInt(int64(t), 10), nil
    }
```

### The Problem
`float64` to `int64` conversion is unsafe for extremely large integers (over $2^{53}$). While common IDs fit within this limit, IDs exceeding it will suffer precision loss. Additionally, converting a float that is larger than `math.MaxInt64` to `int64` is undefined behavior in Go.

### Recommendation
Add bounds checking to the integer conversion:
```go
if t >= math.MinInt64 && t <= math.MaxInt64 && t == float64(int64(t)) {
    return strconv.FormatInt(int64(t), 10), nil
}
```

---

## 🟢 MINOR: Shadowing Built-in `close`

In `internal/dag/extract.go`'s `tokenize` function:
```go
close := strings.IndexByte(path, ']')
```

### The Problem
`close` is a built-in keyword in Go for closing channels. While shadowing is allowed, it’s a code smell that makes the code harder to maintain and can prevent usage of the real `close()` function in that scope.

### Recommendation
Rename the variable to `closeIdx` or `bracketEnd`.

---

## 🟢 MINOR: Map Allocation Optimization

In `internal/dag/extract.go`'s `ExtractAll`:
```go
results = make(map[string]string, len(paths))
```

### The Problem
If all paths fail to resolve, an empty map is returned. Since `ExtractAll` is called trillions of times in load tests, returning `nil` on a total failure is more GC-friendly.

### Recommendation
Only initialize the map when the first successful extraction occurs or return `nil` if `len(results) == 0`.
