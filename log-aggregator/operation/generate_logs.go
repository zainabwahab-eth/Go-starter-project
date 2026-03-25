package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"time"
)

func main() {
	//set ticker for 1 second stop
	ticker := time.NewTicker(time.Second / time.Duration(1))
	defer ticker.Stop()

	//Different levels for log
	levels := []string{"INFO", "WARN", "ERROR"}

	//Different log messages
	messages := []string{
		"Server started on port 3000",
		"Memory usage above 80%",
		"Database connection failed",
	}

	//Create file (if it does not exist) and check for error
	file, createErr := os.Create("app.log")

	if createErr != nil {
		fmt.Println("Error", createErr)
		return
	}
	defer file.Close()

	//Infinity loop
	for {
		<-ticker.C

		//Randomly pick level and message string
		level := levels[rand.IntN(len(levels))]
		message := messages[rand.IntN(len(messages))]

		//Write string line
		line := fmt.Sprintf("%s %s %s", time.Now().Format("2006/01/02 15:04:05"), level, message)

		//write it in file
		file.WriteString(string(line) + "\n")
	}

}
