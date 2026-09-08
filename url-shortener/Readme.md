# URL Shortener with Analytics Pipeline

A URL shortener with a real-time async analytics pipeline. The interesting part isn't the shortening — it's what happens every time someone clicks a link: device, browser, and OS data is captured and stored asynchronously without slowing down the redirect.

---

## How It Works

```
POST /shorten → generate cryptographically random short code → store in PostgreSQL

GET /{code}  → look up original URL → redirect instantly
                        ↓
              click event dropped into buffered channel (non-blocking)
                        ↓
              background goroutine picks it up
                        ↓
              parses user agent → saves enriched click to PostgreSQL

GET /stats/{code}          → click count, device, browser, OS breakdown
GET /stats/{code}/timeline → clicks grouped by hour
```

The redirect happens in milliseconds. The analytics processing happens in the background without adding any latency to the user experience.

---

## Setup

**1. Start PostgreSQL with Docker:**
```bash
docker run --name url-shortener-db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=urlshortener \
  -p 5432:5432 \
  -d postgres
```

**2. Create the tables:**
```bash
docker exec -i url-shortener-db psql -U postgres -d urlshortener < schema.sql
```

**3. Set up environment variables:**
```bash
cp .env.example .env
# edit .env with your database connection string
```

**4. Run the server:**
```bash
go run .
```

Server starts on `http://localhost:8181`.

---

## Endpoints

### Shorten a URL
```bash
curl -X POST http://localhost:8181/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://google.com"}'
```
```json
{
  "message": "Key set successfully",
  "url": {
    "id": 1,
    "short_code": "4n6bzE",
    "original_url": "https://google.com",
    "created_at": "2026-09-01T09:43:12.844114Z"
  }
}
```

### Redirect
```bash
curl -L http://localhost:8181/4n6bzE
# redirects to https://google.com
# click event recorded asynchronously in the background
```

### Click stats
```bash
curl http://localhost:8181/stats/4n6bzE
```
```json
{
  "url": {
    "id": 1,
    "short_code": "4n6bzE",
    "original_url": "https://google.com",
    "created_at": "2026-09-01T09:43:12.844114Z"
  },
  "clicks": [
    {
      "id": 1,
      "short_code": "4n6bzE",
      "clicked_at": "2026-09-01T11:55:29.170416Z",
      "device": "Desktop",
      "browser": "Chrome",
      "os": "Linux x86_64",
      "referrer": ""
    }
  ]
}
```

### Click timeline
```bash
curl http://localhost:8181/stats/4n6bzE/timeline
```
```json
{
  "url": { "short_code": "4n6bzE", "original_url": "https://google.com" },
  "timelines": [
    { "hour": "2026-09-01 11:00", "clicks": 3 },
    { "hour": "2026-09-01 12:00", "clicks": 7 }
  ]
}
```

---

## Rate Limiting

All endpoints are protected by an IP-based token bucket rate limiter — 10 requests burst, 5 requests per second sustained. Exceeding the limit returns a `429 Too Many Requests`.

```bash
# Demonstrates rate limiting
for i in {1..15}; do curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8181/stats/4n6bzE; done
# 200 200 200 200 200 200 200 200 200 200 429 429 429 429 429
```

---

## Project Structure

```
url-shortener/
├── main.go              # Server, DB connection, analytics pipeline
├── schema.sql           # PostgreSQL table definitions
├── .env.example         # Environment variable template
├── handlers/
│   └── handlers.go      # HTTP handlers
├── middleware/
│   └── middleware.go    # Rate limiting middleware
├── repository/
│   ├── repository.go    # Database queries
│   └── shortcode.go     # Short code generator
└── utils/
    └── utils.go         # Shared response types and helpers
```

---

## Concepts Demonstrated

- **Async analytics pipeline** — click events are sent into a buffered channel and processed by a background goroutine. The redirect is never blocked by analytics writes.
- **crypto/rand short codes** — uses the OS entropy pool instead of `math/rand` so short codes are cryptographically unpredictable and not enumerable.
- **Token bucket rate limiting** — `golang.org/x/time/rate` implements a token bucket per IP address. Burst of 10, sustained at 5/sec. Built on the same algorithm as the CLI load tester in this repo.
- **User agent parsing** — device, browser, and OS extracted from the `User-Agent` header on every click using `mssola/useragent`.
- **PostgreSQL with foreign keys** — the `clicks` table references `urls(short_code)`, so the database enforces referential integrity without application-level checks.
- **Package structure** — `repository`, `handlers`, `middleware`, and `utils` are separate packages to avoid circular dependencies and keep concerns separated.
- **Environment variables** — database connection string loaded from `.env` using `godotenv`, never hardcoded.

---

## Why This Matters

URL shorteners are a common backend interview topic — but the analytics pipeline is what makes this one interesting. In production systems like Bit.ly, the redirect path is kept as fast as possible while analytics are processed asynchronously. Building this pattern from scratch — buffered channel, background goroutine, enrichment step — gives you a concrete understanding of how async event pipelines work, which is the same architecture behind systems like Kafka consumers and webhook processors.

---

## What I Learned

The most surprising thing was how simple the async pipeline turned out to be in Go. In Node.js, you'd reach for a message queue like BullMQ or Redis Pub/Sub to decouple the redirect from the analytics write. In Go, a buffered channel and a goroutine replaces all of that for a single-instance service. The channel acts as the queue, the goroutine is the consumer, and the buffer means the producer (redirect handler) never blocks. It's not a replacement for a real message queue at scale — but understanding why taught me exactly when you do and don't need one.

---

## Author

**Zainab Wahab** — Backend Engineer (Node.js / Go)

[GitHub](https://github.com/zainabwahab-eth) · [LinkedIn](https://linkedin.com/in/zainab-wahab-8280ba326)
