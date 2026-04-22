package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type LogEntry struct {
	Level     string `json:"level"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func watchFile(filename string, lines chan<- string) {
	//os.Open keeps the file open so it can be constantly read. It is different from
	//os.Readfile that read te file once, save it in the memory and close the file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	reader := bufio.NewReader(file)

	for {

		line, err := reader.ReadString('\n')

		if err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		lines <- strings.TrimSuffix(line, "\n")
	}

}

func parseLogLine(line string) (LogEntry, error) {
	s := strings.SplitN(line, " ", 4)

	if len(s) != 4 {
		return LogEntry{}, fmt.Errorf("invalid log line: %s", line)
	}

	return LogEntry{
		Level:     s[2],
		Message:   s[3],
		Timestamp: s[0] + " " + s[1],
	}, nil
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func createHandler(entries chan LogEntry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			fmt.Println("Upgrade Error", err)
			return
		}
		defer conn.Close()

		for entry := range entries {
			err := conn.WriteJSON(entry)

			if err != nil {
				fmt.Println("Write Error:", err)
				break
			}
		}
	}
}

func main() {
	lines := make(chan string)
	entries := make(chan LogEntry)
	go watchFile("operation/app.log", lines)
	// fmt.Println("HHHiiii")

	go func() {
		for line := range lines {
			entry, err := parseLogLine(line)
			if err != nil {
				fmt.Println(err)
			}
			entries <- entry
			// fmt.Println(entry)
		}
	}()

	http.HandleFunc("/ws", createHandler(entries))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.ListenAndServe(":8181", nil)

}
