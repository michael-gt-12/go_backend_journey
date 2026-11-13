package main

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID   int
	Data string
}

type Result struct {
	JobID    int
	WorkerID int
	Value    int
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		time.Sleep(50 * time.Millisecond) // simulate work so jobs spread across workers
		results <- Result{JobID: job.ID, WorkerID: id, Value: len(job.Data) * job.ID}
	}
}

func main() {
	jobs := make(chan Job, 20)
	results := make(chan Result, 20)

	var wg sync.WaitGroup
	for w := 1; w <= 5; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	for j := 1; j <= 20; j++ {
		jobs <- Job{ID: j, Data: "task"}
	}
	close(jobs)

	// only this goroutine knows when all workers are done, so it closes results
	go func() {
		wg.Wait()
		close(results)
	}()

	count := 0
	for r := range results {
		count++
		fmt.Printf("worker %d handled job %2d -> %d\n", r.WorkerID, r.JobID, r.Value)
	}
	fmt.Println("results received:", count)
}
