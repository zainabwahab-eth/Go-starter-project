# Real-time Log Aggregator

A tool that tails log files and streams live entries to a browser dashboard over WebSockets. A miniature version of what Datadog does under the hood.

---

## How It Works

```
log generator writes to app.log
        ↓
watchFile detects new lines (tail -f behaviour)
        ↓
parseLogLine converts raw string → LogEntry struct
        ↓
WebSocket server broadcasts JSON to connected browsers
        ↓
dashboard renders entries live, colour coded by severity
```

Two programs run simultaneously — the log generator simulates a real application writing logs, and the aggregator server watches the file and streams updates to the browser in real time.

---

## Usage

**Terminal 1 — start the log generator:**
```bash
cd operation
go run generate_logs.go
```

**Terminal 2 — start the aggregator server:**
```bash
go run .
```

Open `http://localhost:8181` in your browser. Log entries appear live as they are written to `app.log`.

---

## Sample Output

```
[2026/04/22 16:01:12] ERROR: Database connection failed
[2026/04/22 16:01:13] WARN: Memory usage above 80%
[2026/04/22 16:01:14] INFO: Server started on port 3000
[2026/04/22 16:01:15] ERROR: Database connection failed
```

Entries are colour coded — green for INFO, orange for WARN, red for ERROR.

---

## Project Structure

```
log-aggregator/
├── main.go              # File watcher, log parser, WebSocket server
├── index.html           # Browser dashboard
└── operation/
    └── generate_logs.go # Simulates an application writing logs
```

---

## Concepts Demonstrated

- **File tailing** — `bufio.NewReader` holds the file open and picks up new lines as they arrive, replicating `tail -f` behaviour
- **Goroutines + channels** — file watcher, parser, and HTTP server run concurrently, passing data through typed channels
- **WebSockets** — `gorilla/websocket` upgrades HTTP connections and streams JSON to the browser in real time
- **Struct tags** — `LogEntry` fields serialise to lowercase JSON keys consumed directly by the frontend
- **Closures** — WebSocket handler receives the entries channel via closure rather than a global variable

---

## What I Learned

Coming from Node.js where WebSocket servers require juggling event emitters and callbacks, Go's approach felt more explicit — a goroutine ranges over a channel, and whatever arrives gets written to the connection immediately. The trickiest part was understanding that `bufio.Scanner` stops permanently at EOF while `bufio.NewReader` retries from the same position, which is what makes live tailing possible.

---

## Author

**Zainab Wahab** — Backend Engineer (Node.js / Go)

[GitHub](https://github.com/zainabwahab-eth) · [LinkedIn](https://linkedin.com/in/zainab-wahab-8280ba326)