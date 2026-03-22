# 🚀 ReqX: High-Performance, Scriptable API Client

**ReqX** is a lightweight, terminal-centric API execution engine built in Go. Designed for developers who value speed, automation, and a clean CLI experience, ReqX allows you to run Postman-style collections, debug real-time Socket.IO streams, and automate complex load tests—all with the power of an embedded SQLite history engine and a sleek web dashboard.

---

## ⚡ Quick Install (Windows)

### Option 1: PowerShell (Recommended)
Open PowerShell as **Administrator** and run:
```powershell
iwr -useb https://raw.githubusercontent.com/aryanwalia2003/reqx/main/install.ps1 | iex
```

### Option 2: Command Prompt (CMD)
Open CMD as **Administrator** and run:
```cmd
powershell -ExecutionPolicy ByPass -Command "iwr -useb https://raw.githubusercontent.com/aryanwalia2003/reqx/main/install.ps1 | iex"
```

---

## ✨ Features at a Glance

- **🚀 Blazing Fast**: Engineered in Go with `sync.Pool` optimizations and pre-compiled JS bytecode for near-zero overhead.
- **📊 Local History (v2.1)**: Every test run is automatically indexed in an embedded SQLite database (WAL mode).
- **📡 Protocol Support**: Full HTTP/HTTPS support and an interactive Socket.IO v4 / WebSocket REPL.
- **🔐 Stateful Flows**: Automatically handles environment variables, cookie persistence, and auth inheritance.
- **📜 JS Scripting**: Use familiar Postman-style JavaScript (`pm.env.set`, `pm.test`) for advanced test logic.
- **🔄 Multi-Iteration**: Run massive load tests with high concurrency and view aggregated HDR Histogram metrics.
- **🛠 Zero Dependencies**: A single, portable binary that works right out of the box.

---

## 📊 Local History & Dashboard (v2.1)

ReqX now includes a built-in historical tracking system. Never lose a test run again.

### 🖥️ The Web UI
Launch the embedded dashboard to visualize performance trends, latency heatmaps, and per-request breakdowns:
```bash
reqx ui
```
*Dashboards are served locally via an embedded web server, keeping your data private and offline.*

### 🔍 Drilldown Analysis
The UI allows you to click into any historical run to see exactly which requests failed, their P95 latency, and throughput—helping you spot regressions instantly.

---

## 🚀 Performance Architecture

ReqX is built for scale. Our v2.1 core features:
- **Goja VM Pooling**: Reuses JavaScript runtimes to eliminate GC pressure during heavy load.
- **Bytecode Caching**: Pre-compiles JS scripts in the planning phase to save CPU cycles per iteration.
- **HDR Histograms**: Provides high-fidelity latency percentiles (P95, P99) with O(1) performance.
- **Zero-Write Contention**: SQLite writes are batch-processed on the main thread after runs to ensure 0% impact on test throughput.

---

## 📚 Core Command Guide

### 1. `run`: The Load Engine
Execute full collections with variables, cookies, and iterations.
```bash
# Basic run with environment
reqx run collection.json -e dev.json

# Performance Test: 500 iterations with 50 workers
reqx run collection.json -n 500 -w 50

# Target specific requests by name
reqx run api.json -f "Login" -f "Update Session"
```

### 2. `ui`: The Dashboard
Visualize your history and performance trends.
```bash
reqx ui  # Opens http://localhost:8090 automatically
```

### 3. `sio` & `ws`: The Real-time REPL
Debug streams interactively.
```bash
# Connect with a session cookie
reqx sio http://localhost:7879 -H "Cookie: auth={{token}}"

# Inside the REPL:
> listen NEW_MESSAGE
> emit send_chat {"text": "hello"}
```

### 4. `collection`: The CLI Editor
Modify your test suites without leaving the terminal.
```bash
# Add a health check request
reqx collection add api.json -n "Health" -u "{{base_url}}/health"
```

---

## 🤝 Contributing & Architecture

ReqX follows a strict **Interface-Driven Design**. If you are contributing, please refer to our architecture guides:
- [Local History UI Architecture](docs/guides/local_history_ui_architecture.md)
- [Performance Optimization Guide](docs/guides/hdr-histogram-optimization.md)
- [System Diagram](docs/diagrams/local_history_ui_diagram.md)

*Developed by Aryan Walia | 2026*