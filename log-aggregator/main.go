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

	scanner := bufio.NewScanner(file)

	for {

		for scanner.Scan() {
			lines <- string(scanner.Text())
		}
		time.Sleep(500 * time.Millisecond)
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
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	conn.WriteJSON()
}

func main() {
	lines := make(chan string)
	go watchFile("operation/app.log", lines)

	for line := range lines {
		entry, err := parseLogLine(line)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(entry)
	}

}
