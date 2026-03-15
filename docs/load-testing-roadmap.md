# ReqX: Load Testing Roadmap (k6-Lite)

This roadmap outlines the transformation of **ReqX** from a sequential API runner into a high-performance, lightweight load testing tool. The objective is to achieve **"Performance of k6, simplicity of Postman collections."**

---

## Phase 1: The Concurrency Engine (Foundation)
*Goal: Run multiple copies of the same collection in parallel.*

- **Introduce Workers (Virtual Users):**
    - Add `-c` or `--workers` flag to the `run` command.
    - Implement a Worker Pool pattern using goroutines and `sync.WaitGroup`.
- **Isolate Worker State:**
    - Each worker must have an independent, deep-copied `Environment`.
    - Ensure separate `CollectionRunner` instances per worker (unique cookie jars, etc.).
- **Job Distribution:**
    - Use a Go channel as a job queue to distribute total iterations (`-n`) across workers (`-c`).

---

## Phase 2: High-Fidelity Metrics (The Science)
*Goal: Provide deep performance insights beyond simple averages.*

- **Centralized Metrics Collector:**
    - Use a thread-safe channel to aggregate `RequestMetric` data from all workers.
- **Advanced Statistics:**
    - **Percentiles:** Calculate P90, P95, and P99 response times.
    - **Throughput:** Calculate Requests Per Second (RPS) and data transfer rates (MB/s).
    - **Error Rate:** Track success/failure ratios across the entire load test.
- **Thresholds & Assertions:**
    - Implement flags like `--threshold-p95=500ms` or `--threshold-error-rate=1%`.
    - Exit with non-zero code if thresholds are breached (CI/CD integration).

---

## Phase 3: Realistic Load Simulation (The Art)
*Goal: Simulate real-world traffic patterns.*

- **Duration-based Tests:**
    - Add `-d` or `--duration` flag (e.g., `5m`) to run continuously for a fixed time.
- **Rate Limiting (Constant RPS):**
    - Add `--rps` flag to control the exact injection rate using a `time.Ticker`.
- **Load Ramping (Stages):**
    - Implement "stages" to simulate ramp-up and ramp-down periods.

---

## Phase 4: Data-Driven Testing (The Realism)
*Goal: Dynamically inject unique data into each Virtual User (VU).*

- **External Data Feeders:**
    - Add `--data-file` support for CSV/JSON files.
- **Atomic Iteration Handling:**
    - Use atomic counters to ensure each worker "checks out" a unique row of data per iteration.
    - Map data columns directly to environment variables (e.g., `{{username}}`).

---

## Immediate Next Step
**Start with Phase 1.** Implementing the `--workers` flag is the core architectural shift. Once ReqX can run 10 parallel goroutines with isolated state, the foundation for a true load testing engine will be complete.
