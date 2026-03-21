# Architecture: Local History & Embedded UI

```mermaid
flowchart TD
    subgraph RUN ["reqx run (Write Path)"]
        A[Load Test Ends] --> B["metrics.AnalyzeSharded()"]
        B --> C["metrics.PrintReport()"]
        C --> D["history.Open()"]
        D --> E["db.SaveRun() — 1 transaction"]
        E --> F[("~/.reqx/history.db\nWAL mode")]
    end

    subgraph UI ["reqx ui (Read Path)"]
        G["NewUICmd()"] --> H["history.Open()"]
        H --> I["ui.NewServer(db, 8090)"]
        I --> J["server.Start()"]
        J --> K["/ → embedded index.html"]
        J --> L["/api/history → db.ListRuns(50)"]
        L --> F
    end

    subgraph BROWSER ["Browser"]
        M["localhost:8090"] --> K
        M -->|fetch /api/history| L
        M --> N["Chart.js: P95 Trend Line"]
        M --> O["Table: Recent Runs"]
    end

    %% Guarantee: zero concurrent writes
    NOTE["⚡ No contention: Write happens once,\nafter test, on the main goroutine.\nSQLite WAL allows concurrent reads."]

    style RUN fill:#1b4d3e,stroke:#2e7d32,color:#ffffff
    style UI fill:#1565c0,stroke:#0d47a1,color:#ffffff
    style BROWSER fill:#37474f,stroke:#263238,color:#ffffff
    style F fill:#311b92,stroke:#1a237e,color:#ffffff
    style NOTE fill:#1a1d27,stroke:#6366f1,color:#a5b4fc
```

## Package Structure

| Package | Files | Purpose |
| :--- | :--- | :--- |
| `internal/history` | `db_struct.go`, `db_ctor.go`, `db_method.go`, `db_query_method.go` | SQLite open/migrate, write (SaveRun), read (ListRuns) |
| `internal/ui` | `embed.go`, `server_struct.go`, `server_ctor.go`, `server_method.go`, `assets/index.html` | HTTP server, embedded assets, browser launch |
| `cmd` | `ui_cmd_ctor.go` | `reqx ui` Cobra command |

## No-Contention Guarantee

The `history.SaveRun()` is called **once per test**, synchronously, on the main goroutine, **after** `PrintReport()` completes. SQLite WAL mode means the `reqx ui` server can read the DB while a test is running without any locking contention.
