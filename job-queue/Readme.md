# Priority Queue
A Priority queue using Go's container/heap that let you set jobs (Tasks) with priority level and automatically process those tasks in the background based on their priority level.

---

## How It Works

```
POST /jobs  → add a new job to the queue
GET /jobs   → get all jobs in the queue
GET /jobs/{id} → get a job based on ID
GET /jobs/failed  → get all failed jobs

background goroutine worker checks the queue everytime
        ↓
Remove the greater priority job from queue and process it -> If process is successful, mark job status as succeessful
        ↓
If processing fails increase attempt and add job back to queue to be retried 
```

---

## Usage

```bash
go run .
```

Server starts on `http://localhost:8181`.

---

## Endpoints

### Add a job
```bash
curl -X POST http://localhost:8181/jobs \
  -H "Content-Type: application/json" \
  -d '{"type": "invoice", "payload": "send invoice `#123`", "priority": 2}'
```

```json
{"message":"Job add successfully","job":{"id":"373805a7-3807-4a82-a4e5-05f298017131","type":"password_reset","payload":"kiki@gmail.com","priority":3,"status":"running","attempts":0,"maxAttempts":3,"createdAt":"2026-08-17T14:31:25.951136023+01:00"}}
```

### Get all jobs
```bash
curl "http://localhost:8181/jobs"
```
```json
{"message":"All Jobs Found","jobs":[{"id":"d257e827-6f6c-4f7f-8f2d-7d5b6252cd9e","type":"invoice","payload":"generate invoice #123","priority":2,"status":"completed","attempts":0,"maxAttempts":3,"createdAt":"2026-08-17T14:21:30.599481465+01:00"},{"id":"f16aa791-d15f-4a56-9109-665d3bf2286f","type":"newsletter","payload":"send to 10000 users","priority":1,"status":"completed","attempts":0,"maxAttempts":3,"createdAt":"2026-08-17T14:21:30.557822335+01:00"},{"id":"b9b2bd44-a662-4633-b2da-efde89851816","type":"password_reset","payload":"zainab@gmail.com","priority":3,"status":"completed","attempts":0,"maxAttempts":3,"createdAt":"2026-08-17T14:21:30.579932723+01:00"}]}
```

### Get a job
```bash
curl "http://localhost:8181/jobs/d257e827-6f6c-4f7f-8f2d-7d5b6252cd9e"
```
```json
{"message":"Successful","job":{"id":"d257e827-6f6c-4f7f-8f2d-7d5b6252cd9e","type":"invoice","payload":"generate invoice #123","priority":2,"status":"completed","attempts":0,"maxAttempts":3,"createdAt":"2026-08-17T14:21:30.599481465+01:00"}}
```

### Get all failed jobs
```bash
curl http://localhost:8181/jobs/failed
```
```json
{"message":"Failed jobs found","jobs":[{"id":"373805a7-3807-4a82-a4e5-05f298017131","type":"password_reset","payload":"kiki@gmail.com","priority":3,"status":"failed","attempts":3,"maxAttempts":3,"createdAt":"2026-08-17T14:31:25.951136023+01:00","error":"always fails"}]}
```

---

## Project Structure

```
job-queue/
└── main.go    # Priority queue, worker, retry logic, handlers, server
```

---

## Concepts Demonstrated

- **sync.RWMutex** — protects the shared map from race conditions. Multiple goroutines can hold a read lock simultaneously; writes get an exclusive lock. This is how concurrent access to shared data is handled safely in Go.
- **container/heap** — a Golang package that implements a heap data structure. It keeps the highest priority on top so that can grab it in 0(1) time.
- **Exponential backoff retry logic** — when a job fails we retry exponentially to give a temporary problem (db overload, network timeout) time to resolve itself. We check this retry attempts by a max-retry number.
- **Closures** — handlers receive the shared store via closure rather than a global variable, keeping the code testable and modular.
- **UUID job IDs** — google UUID is used to generate unique id for each job.

---

## Why This Matters

User experience is one of the most important things in engineering today. When a process takes time to complete, it is better to handle it in the background asynchronously and let users continue their activity rather than have them staring at a loading screen. And since not all background jobs are equal — a password reset email is more urgent than a newsletter — priority scoring ensures the most critical work always gets done first.

---

## What I Learned

The most interesting part of this project was understanding how container/heap works in Go. You don't use a pre-built priority queue — you implement an interface with five methods (Len, Less, Swap, Push, Pop) and Go's heap package handles all the sorting logic for you. The Less method is where the priority logic lives — returning pq[i].Priority > pq[j].Priority means higher numbers always bubble to the top. Building it from scratch made me understand why priority queues are more efficient than sorting a slice on every insertion.

---

## Author

**Zainab Wahab** — Backend Engineer (Node.js / Go)

[GitHub](https://github.com/zainabwahab-eth) · [LinkedIn](https://linkedin.com/in/zainab-wahab-8280ba326)
