package main

import (
	"fmt"
	"net/http"
	"sort"
	"time"
)

type Result struct {
	StatusCode int
	Duration   time.Duration
	Error      error
}

type Stats struct {
	TotalRequests int
	Successful    int
	Failed        int
	p50           time.Duration
	p95           time.Duration
	p99           time.Duration
	TotalTime     time.Duration
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

func runWorkerPool(url string, totalRequests int, concurrency int, ratePerSecond int) []Result {
	//Create 2 channels; jobs carries URLs (string) from main and workers consume it
	//results carries Results out from workers and main collect it
	jobs := make(chan string, totalRequests)
	results := make(chan Result)
	finalRes := []Result{}

	N := concurrency
	for i := 1; i <= N; i++ {

		//Loop concurrency times to create workers. We have 3 workers (based on N)
		//and 10 jobs item (based on totalrequests) the workers pick on the job one by
		//one till the 10 is done. So in the end we get 10 results (depending on totalRequest)
		//Once a job (url) is picked by one worker the other workers does not know about it
		//and the job is done. Each of the workers pick different jobs
		go func() {
			for item := range jobs {
				results <- makeRequest(item)
			}
		}()
	}

	//Create for ticker rateLimiting
	ticker := time.NewTicker(time.Second / time.Duration(ratePerSecond))
	defer ticker.Stop()

	// totalrequest is the amount of time each worker call the url
	for j := 1; j <= totalRequests; j++ {
		//Wait until ticker fires, then;
		<-ticker.C
		//send one url to the job channel
		jobs <- url
	}
	close(jobs)

	for k := 1; k <= totalRequests; k++ {
		finalRes = append(finalRes, <-results)
	}
	return finalRes
}

func calcPercentile(percentile int, durations []time.Duration) int {
	return (percentile * len(durations)) / 100
}

func calculateStats(results []Result, totalTime time.Duration) Stats {
	success := 0
	failure := 0
	durations := []time.Duration{}
	for _, r := range results {
		if r.StatusCode == 200 {
			success++
		} else {
			failure++
		}

		durations = append(durations, r.Duration)
	}

	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
	// fmt.Println(durations)

	return Stats{
		TotalRequests: len(results),
		Successful:    success,
		Failed:        failure,
		p50:           durations[calcPercentile(50, durations)],
		p95:           durations[calcPercentile(95, durations)],
		p99:           durations[calcPercentile(99, durations)],
		TotalTime:     totalTime,
	}
}

func main() {
	url := "https://google.com"

	start := time.Now()
	results := runWorkerPool(url, 10, 3, 2)
	totalTime := time.Since(start)
	// var totalTime time.Duration
	// for _, r := range results {
	// 	totalTime += r.Duration
	// }

	stats := calculateStats(results, totalTime)
	fmt.Printf("Total Requests: %d\n", stats.TotalRequests)
	fmt.Printf("Successful: %d\n", stats.Successful)
	fmt.Printf("Failed: %d\n", stats.Failed)
	fmt.Printf("p50: %v\n", stats.p50)
	fmt.Printf("p95: %v\n", stats.p95)
	fmt.Printf("p99: %v\n", stats.p99)
	fmt.Printf("Total Time: %v\n", stats.TotalTime)

	// for i, r := range results {
	// 	fmt.Printf("[%d] Status: %d | Duration: %v\n", i+1, r.StatusCode, r.Duration)
	// }

	// for i, r := range results {
	// 	fmt.Printf("[%d] Status: %d | Duration: %v\n", i+1, r.StatusCode, r.Duration)
	// }
	// reresults <- makeRequest(url)sult := makeRequest(url)

	// if result.Error != nil {
	// 	fmt.Println("Request Error", result.Error)
	// 	return
	// }

	// fmt.Printf("status: %d | duration: %v\n", result.StatusCode, result.Duration)
}
