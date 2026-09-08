# Go Starter Projects

A collection of backend projects built in Go while transitioning from Node.js. Each project is a miniature version of a real production tool built from scratch to understand how things actually work under the hood.

---

## Projects

### 01. [CLI Load Tester & Rate Limiter](./cli-load-tester-and-rate-limiter)
A command-line tool that hammers any HTTP endpoint with configurable concurrency and rate limiting, then outputs a latency report with p50/p95/p99 percentiles. A baby version of k6 or Apache Bench.

[→ README](./cli-load-tester-and-rate-limiter/Readme.md)

---

### 02. [Real-time Log Aggregator](./log-aggregator)
Tails a log file and streams live entries to a browser dashboard over WebSockets, colour coded by severity. A miniature Datadog.

[→ README](./log-aggregator/Readme.md)

---

### 03. [In-Memory Key-Value Store](./kv-store)
An HTTP-accessible key-value store with TTL expiry and automatic background eviction. Set a key with an expiry time and it disappears automatically when the TTL runs out. A miniature Redis.

[→ README](./kv-store/Readme.md)

---

### 04. [Priority Job Queue](./job-queue)
A job queue where tasks have priority levels — higher priority jobs are processed first regardless of submission order. Supports retry with exponential backoff and a dead letter queue for permanently failed jobs. A miniature Sidekiq.

[→ README](./job-queue/Readme.md)

---

### 05. [URL Shortener with Analytics Pipeline](./url-shortener)
A URL shortener with a real-time async analytics pipeline. The redirect is instant — click events are captured and enriched (device, browser, OS) in the background without adding any latency. Includes IP-based rate limiting and a timeline endpoint showing clicks grouped by hour.

[→ README](./url-shortener/Readme.md)

---

## Why Go

Coming from a Node.js/Express background, I built these projects to understand Go's concurrency model by actually using it. The biggest shift was goroutines — unlike Node's single-threaded event loop, Go achieves real parallelism, which changes how you think about backend architecture at scale.

---

## Author

**Zainab Wahab** — Backend Engineer (Node.js / Go)

[GitHub](https://github.com/zainabwahab-eth) · [LinkedIn](https://linkedin.com/in/zainab-wahab-8280ba326)
