# Scripting Engine Requirements

To build a fully functional Postman-compatible CLI, our JavaScript execution engine (powered by Goja) needs to support the following core capabilities within pre-request and test scripts.

## 1. Environment Variable Management (`pm.environment`)
Scripts frequently need to interact with the environment state to chain requests together.
- `pm.environment.get("key")`: Retrieve the string value of a variable.
- `pm.environment.set("key", "value")`: Set or update a variable.
- `pm.environment.unset("key")`: Remove a variable from the context entirely.

## 2. Console Logging Interception
Developers rely on console output to debug their logic. The Goja runtime must intercept native JS console calls and route them beautifully to the CLI terminal.
- `console.log(...)`: Standard information logging.
- `console.warn(...)`: Warning messages (yellow).
- `console.error(...)`: Error messages (red).
- `console.dir(...)`: Pretty-printing object structures.
- `console.table(...)`: Formatting arrays/objects into a readable terminal table.

## 3. Response Handling & Dynamic Tokens (`pm.response`)
To parse authentication tokens deeply nested in JSON or Headers and push them into the environment, scripts need full access to the HTTP transaction results.
- `pm.response.json()`: Parse the raw response body into a JavaScript Object.
- `pm.response.text()`: Access the raw response body as a string.
- `pm.response.headers.get("key")`: Extract specific response headers (e.g., `Authorization` or `Set-Cookie`).

## 4. Testing & Assertions (`pm.test` & `pm.expect`)
- `pm.test("Test Name", function() { ... })`: Define boundaries for test suites so the runner can track Pass/Fail metrics.
- `pm.expect(value).to...`: Lightweight assertions to validate response codes, body tokens, and latencies.