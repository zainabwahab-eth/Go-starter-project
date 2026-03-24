# CLI Load Tester & Rate Limiter

A command-line tool that load tests any HTTP endpoint with configurable concurrency and rate limiting, then outputs a detailed latency report.

Built in Go from scratch — no third-party HTTP or stats libraries.

---

## How It Works

The tool spawns a pool of worker goroutines that process requests concurrently. A rate limiter controls how many requests are dispatched per second using Go's `time.Ticker`. Results are collected and analysed to produce latency percentiles.

```
main goroutine
     │
     ├── spawns N worker goroutines (concurrency)
     ├── feeds jobs into a channel at X req/s (rate limiter)
     └── collects results and calculates stats
```

---

## Usage

```bash
go run main.go -url <url> -requests <n> -concurrency <n> -rate <n>
```

### Flags

| Flag | Description | Default |
|---|---|---|
| `-url` | Target URL to test | `https://google.com` |
| `-requests` | Total number of requests to send | `10` |
| `-concurrency` | Number of concurrent workers | `3` |
| `-rate` | Requests per second (rate limit) | `3` |

### Example

```bash
go run main.go -url https://api.example.com/users -requests 100 -concurrency 10 -rate 5
```

---

## Sample Output

```
Total Requests: 100
Successful:     98   (98%)
Failed:         2    (2%)

Latency
  p50:  742ms
  p95:  1.231s
  p99:  1.812s

Total Time: 21.4s
```

---

## Concepts Demonstrated

- **Goroutines** — lightweight concurrent workers, one per concurrency slot
- **Channels** — typed pipes coordinating jobs in and results out between goroutines
- **Worker pool pattern** — N workers sharing a fixed job queue, each picking tasks independently
- **Token bucket rate limiting** — `time.Ticker` controls dispatch pace without a third-party library
- **Latency percentiles** — p50/p95/p99 calculated from sorted duration slices

---

## Project Structure

```
cli-load-tester/
└── main.go       # All logic — worker pool, rate limiter, stats, CLI
```

---

## What I Learned

This was my first Go project coming from a Node.js background. The biggest shift was understanding that goroutines achieve real parallelism — unlike Node's event loop which is single-threaded. Writing `Promise.all` in Node simulates concurrency; Go's worker pool actually runs requests simultaneously across CPU cores.

The `<-channel` syntax for blocking receives replaced what I'd normally do with `async/await` — and turned out to be more explicit and easier to reason about once it clicked.