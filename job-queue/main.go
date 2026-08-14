package main

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	StatusPending   = "pending"
	StatusRunning   = "running"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
)

type Job struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Payload     string    `json:"payload"`
	Priority    int       `json:"priority"`
	Status      string    `json:"status"`
	Attempts    int       `json:"attempts"`
	MaxAttempts int       `json:"maxAttempts"`
	CreatedAt   time.Time `json:"createdAt"`
	Error       string    `json:"error,omitempty"`
}

type PriorityQueue []*Job

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority > pq[j].Priority
}

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	*pq = append(*pq, x.(*Job))
}

func (pq *PriorityQueue) Pop() any {
	var prev = *pq
	var prevLen = len(prev)
	var last = prev[prevLen-1]
	*pq = prev[0 : prevLen-1]
	return last
}

type Queue struct {
	mu   sync.RWMutex
	pq   PriorityQueue
	jobs map[string]*Job
}

func (q *Queue) Enqueue(job *Job) {
	q.mu.Lock()
	defer q.mu.Unlock()
	heap.Push(&q.pq, job)
	q.jobs[job.ID] = job
}

func (q *Queue) Dequeue() *Job {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.pq.Len() == 0 {
		return nil
	}

	return heap.Pop(&q.pq).(*Job)
}

func (q *Queue) GetJob(id string) *Job {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.jobs[id]
}

func (q *Queue) GetFailedJobs() []*Job {
	jobs := []*Job{}
	q.mu.RLock()
	for _, job := range q.jobs {
		if job.Status == StatusFailed {
			jobs = append(jobs, job)
		}
	}
	q.mu.RUnlock()
	return jobs
}

func (q *Queue) GetAllJobs() []*Job {
	jobs := []*Job{}
	q.mu.RLock()
	for _, job := range q.jobs {
		jobs = append(jobs, job)
	}
	q.mu.Unlock()
	return jobs
}

func NewQueue() *Queue {
	q := &Queue{
		jobs: make(map[string]*Job),
		pq:   PriorityQueue{},
	}

	heap.Init(&q.pq)
	return q
}

func processJob(job *Job) error {
	fmt.Printf("Processing job %s of type %s\n", job.ID, job.Type)
	time.Sleep(2 * time.Second)
	if rand.IntN(2) == 0 {
		return fmt.Errorf("random failure")
	}
	return nil
}

func Worker(q *Queue) {
	for {
		var job *Job
		if job = q.Dequeue(); job == nil {
			time.Sleep(500 * time.Microsecond)
			continue
		}

		job.Status = StatusRunning

		var err error
		if err = processJob(job); err == nil {
			job.Status = StatusCompleted
		} else {
			job.Attempts++

			if job.Attempts < job.MaxAttempts {
				time.Sleep(time.Duration(math.Pow(2, float64(job.Attempts))) * time.Second)
				q.Enqueue(job)
			} else {
				job.Status = StatusFailed
				job.Error = err.Error()
			}
		}
	}

}

type Response struct {
	Message string `json:"message,omitempty"`
	Jobs    []*Job `json:"jobs,omitempty"`
	Job     *Job   `json:"job,omitempty"`
}

func writeResponse(status int, res *Response, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(res)
}

func addNewJob(q *Queue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var j JobRequest

		if err := json.NewDecoder(r.Body).Decode(&j); err != nil {
			writeResponse(http.StatusBadRequest, &Response{Message: "something went wrong"}, w)
			return
		}

		if j.Type == "" || j.Payload == "" {
			writeResponse(http.StatusBadRequest, &Response{Message: "type and payload cannot be empty"}, w)
			return
		}

		if j.Priority > 3 || j.Priority < 1 {
			writeResponse(http.StatusBadRequest, &Response{Message: "Priority cannot be greater than 3 or less than 1"}, w)
			return
		}

		job := Job{
			ID:          uuid.New().String(),
			Type:        j.Type,
			Payload:     j.Payload,
			Priority:    j.Priority,
			Status:      StatusPending,
			CreatedAt:   time.Now(),
			MaxAttempts: 3,
		}

		q.Enqueue(&job)
		writeResponse(http.StatusCreated, &Response{Message: "Job add successfully", Job: &job}, w)
	}

}

func getAJob(q *Queue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		if id == "" {
			writeResponse(http.StatusBadRequest, &Response{Message: "Invalid ID"}, w)
			return
		}

		job := q.GetJob(id)

		if job == nil {
			writeResponse(http.StatusBadRequest, &Response{Message: "Job not found"}, w)
			return
		}

		writeResponse(http.StatusOK, &Response{Message: "Successful", Job: job}, w)
	}
}

func getFailedJobs(q *Queue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobs := q.GetFailedJobs()

		if len(jobs) == 0 {
			writeResponse(http.StatusOK, &Response{Message: "No failed jobs"}, w)
			return
		}

		writeResponse(http.StatusOK, &Response{Message: "Failed jobs found", Jobs: jobs}, w)
	}
}

func getAllJobs(q *Queue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobs := q.GetAllJobs()

		if len(jobs) == 0 {
			writeResponse(http.StatusOK, &Response{Message: "No failed jobs"}, w)
			return
		}

		writeResponse(http.StatusOK, &Response{Message: "All Jobs Found", Jobs: jobs}, w)
	}
}

type JobRequest struct {
	Type     string `json:"type"`
	Payload  string `json:"payload"`
	Priority int    `json:"priority"`
}

func main() {
	queue := NewQueue()

	go Worker(queue)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /jobs", addNewJob(queue))
	mux.HandleFunc("GET /jobs", getAllJobs(queue))
	mux.HandleFunc("GET /jobs/{id}", getAJob(queue))
	mux.HandleFunc("GET /jobs/failed", getFailedJobs(queue))

	server := &http.Server{
		Handler: mux,
		Addr:    ":8181",
	}

	fmt.Println("Server starting at 8181")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
