package main

import (
	"fmt"
	"net/http"
	"time"
)

type Result struct {
	StatusCode int
	Duration   time.Duration
	Error      error
}

func makeRequest(url string) Result {
	//Record start time
	start := time.Now()

	//Make get request
	resp, err := http.Get(url)
	if err != nil {
		return Result{Error: err, Duration: time.Since(start)}
	}

	//close url
	defer resp.Body.Close()

	//Record end time
	end := time.Since(start)

	//Return result
	return Result{
		StatusCode: resp.StatusCode,
		Duration:   end,
		Error:      nil,
	}
}
func runWorkerPool(url string, totalRequests int, concurrency int) []Result {
	jobs := make(chan string, totalRequests)
	results := make(chan Result)
	finalRes := []Result{}

	N := concurrency
	// t := totalRequests
	//Loop concurrency times
	for i := 1; i <= N; i++ {
		go func() {
			for item := range jobs {
				results <- makeRequest(item)
			}
		}()
	}

	for j := 1; j <= totalRequests; j++ {
		jobs <- url
	}
	close(jobs)

	for k := 1; k <= totalRequests; k++ {
		finalRes = append(finalRes, <-results)
	}
	return finalRes
}

func main() {
	url := "https://google.com"
	result := runWorkerPool(url, 10, 3)
	fmt.Println(result)
	// reresults <- makeRequest(url)sult := makeRequest(url)

	// if result.Error != nil {
	// 	fmt.Println("Request Error", result.Error)
	// 	return
	// }

	// fmt.Printf("status: %d | duration: %v\n", result.StatusCode, result.Duration)
}
