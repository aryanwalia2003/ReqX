# 🚀 Postman CLI

A high-performance, scriptable API client for the terminal. Built in Go for speed, designed for developer productivity.

This tool acts as a lightweight, command-line alternative to UI-heavy applications like Postman. It's perfect for automating complex API test flows, debugging real-time services, and managing collections directly from your terminal.

## ✨ Features

- **HTTP & Socket.IO v4:** First-class support for both standard REST APIs and real-time Socket.IO v4 connections.
- **Stateful Collections:** Run sequential requests where variables (like auth tokens) from one request are used in the next.
- **Async Socket Testing:** Test complex, event-driven flows by running background socket listeners alongside HTTP requests.
- **Powerful Scripting:** Use Postman-style JavaScript (`pm.env.set`, `pm.response.json`) to write tests and handle dynamic data.
- **Advanced Flow Control:** Run collections in multiple iterations (`-n`), filter specific requests (`-f`), and temporarily inject new requests on the fly.
- **Full-featured CLI:** Includes a curl-like `req` command, an interactive `sio` REPL, and a `collection` manager to edit files without a GUI.

---

## 💻 Installation (Windows)

1.  Make sure you have `postman-cli.exe` and `install.ps1` in the same directory.
2.  Right-click on your **PowerShell** icon and select **"Run as Administrator"**.
3.  Navigate to the directory containing the script.
4.  Run the installer:
    ```powershell
    Set-ExecutionPolicy Unrestricted -Scope Process
    .\install.ps1
    ```
5.  **Restart any open terminals** (this is critical!).
6.  Verify the installation by typing `postman-cli --help` in your new terminal window.

---

## 📚 Core Commands

### `run`: Execute a Collection
The main command for running a full test suite.

```bash
# Basic run with an environment file
postman-cli run vuc-collection.json -e vuc-env.json

# Run 10 times with verbose output and no cookies
postman-cli run vuc-collection.json -e vuc-env.json -n 10 -v --no-cookies

# Run only requests with "Login" in their name
postman-cli run vuc-collection.json -f "Login"
```

### `req`: Send a Single Request
A powerful, curl-like command for quick, ad-hoc API calls.

```bash
# Simple GET request
postman-cli req https://httpbin.org/get

# POST request with a body and variables from an env file
postman-cli req "{{base_url}}/auth/login" -e vuc-env.json -X POST -d '{"user":"test"}'
```

### `sio`: Interactive Socket.IO REPL
Connect to a Socket.IO v4 server and debug events in real-time.

```bash
# Connect with an auth cookie
postman-cli sio http://localhost:7879 -H "Cookie: auth-token={{token}}"

# Once inside the REPL:
> listen NEW_DISPATCH
> emit APPOINTMENT_TAKEN {"id": 123}
```

### `collection`: Manage Your JSON Files
Edit your collections without opening a text editor.

```bash
# List all requests with their index numbers
postman-cli collection list vuc-collection.json

# Add a new request to the end of the file
postman-cli collection add vuc-collection.json -n "Health Check" -u "{{base_url}}/health"

# Move request #5 to the #2 position
postman-cli collection move vuc-collection.json 5 2
```

---

## 🧬 Advanced: Async Socket Testing

To test flows where a user (e.g., a Doctor) listens for events triggered by another user (e.g., a Receptionist), you can start a socket connection in the background using `"async": true`.

**Example `collection.json`:**
```json
"requests": [
  {
    "name": "Doctor Connects and Listens (Async)",
    "protocol": "SOCKETIO",
    "async": true,
    "url": "{{socket_base_url}}",
    "headers": { "Cookie": "auth-token={{doctor_token}}" },
    "events": [{ "type": "listen", "name": "NEW_DISPATCH" }]
  },
  {
    "name": "Receptionist Broadcasts (HTTP)",
    "method": "POST",
    "url": "{{base_url}}/appointments/broadcast",
    "auth": { "type": "bearer", "token": "{{receptionist_token}}" }
  }
]
```
When you run this, the Doctor socket will connect and wait. The CLI will *immediately* proceed to the next request, allowing the Receptionist to trigger the `NEW_DISPATCH` event, which will be caught by the background socket.