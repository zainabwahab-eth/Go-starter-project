package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
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

func main() {
	lines := make(chan string)
	go watchFile("operation/app.log", lines)

	for line := range lines {
		fmt.Println(line)
	}
}
