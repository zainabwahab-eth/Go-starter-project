# In-Memory Key-Value Store

An HTTP-accessible key-value store with TTL expiry and automatic background eviction. A miniature Redis — set a key with an expiry time and it disappears automatically when the TTL runs out.

---

## How It Works

```
POST /set  → stores key with value and TTL
GET /get   → returns value if key exists and hasn't expired
DELETE /delete → removes key immediately
GET /keys  → lists all active (non-expired) keys

background goroutine wakes every second
        ↓
checks every key's ExpiresAt time
        ↓
deletes expired keys automatically
```

---

## Usage

```bash
go run .
```

Server starts on `http://localhost:8181`.

---

## Endpoints

### Set a key
```bash
curl -X POST http://localhost:8181/set \
  -H "Content-Type: application/json" \
  -d '{"key": "session:abc123", "value": "zainab", "ttl": 30}'
```
```json
{"message":"Key set successfully"}
```

### Get a key
```bash
curl "http://localhost:8181/get?key=session:abc123"
```
```json
{"value":"zainab"}
```

### Get after expiry
```bash
curl "http://localhost:8181/get?key=session:abc123"
```
```json
{"message":"Key not found or expired"}
```

### Delete a key
```bash
curl -X DELETE "http://localhost:8181/delete?key=session:abc123"
```
```json
{"message":"Key deleted successfully"}
```

### List all active keys
```bash
curl http://localhost:8181/keys
```
```json
{"keys":["session:abc123","user:456"]}
```

---

## Project Structure

```
kv-store/
└── main.go    # Store, eviction, handlers, server
```

---

## Concepts Demonstrated

- **sync.RWMutex** — protects the shared map from race conditions. Multiple goroutines can hold a read lock simultaneously; writes get an exclusive lock. This is how concurrent access to shared data is handled safely in Go.
- **TTL eviction** — a background goroutine wakes every second and deletes keys whose `ExpiresAt` time has passed. This is exactly how Redis handles key expiry internally.
- **Race condition awareness** — without mutex protection, concurrent HTTP requests reading and writing the same map simultaneously would corrupt data or crash the program. Go's race detector (`go run -race .`) catches this.
- **Closures** — handlers receive the shared store via closure rather than a global variable, keeping the code testable and modular.
- **Value vs pointer receivers** — all store methods use pointer receivers so mutations affect the original store, not a copy.

---

## Why This Matters

Key-value stores with TTL are used everywhere in backend engineering — session management, caching, rate limiting, OTP storage, temporary state in multi-step flows. Understanding how TTL eviction works under the hood makes you a better user of Redis and a better systems thinker.

---

## What I Learned

The most interesting part of this project was understanding why a regular mutex wasn't enough. Using `sync.RWMutex` instead of `sync.Mutex` means concurrent GET requests never block each other only writes cause a full lock. For a store that will have far more reads than writes, that distinction matters significantly at scale.

The `defer mu.Unlock()` pattern also clicked here pairing lock and unlock immediately means you can never accidentally forget to release a lock, even if the function returns early due to an error.

---

## Author

**Zainab Wahab** — Backend Engineer (Node.js / Go)

[GitHub](https://github.com/zainabwahab-eth) · [LinkedIn](https://linkedin.com/in/zainab-wahab-8280ba326)
