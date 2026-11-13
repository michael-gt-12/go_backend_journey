package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Job struct{ ID int }

type Result struct {
	JobID    int
	WorkerID int
	Err      error
}

// process takes 200ms, except job 3 which takes 2s and must hit its 500ms timeout.
func process(ctx context.Context, job Job) error {
	d := 200 * time.Millisecond
	if job.ID == 3 {
		d = 2 * time.Second
	}
	jobCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-jobCtx.Done():
		return jobCtx.Err()
	}
}

func worker(ctx context.Context, id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("worker %d stopping: %v\n", id, ctx.Err())
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			err := process(ctx, job)
			select {
			case results <- Result{JobID: job.ID, WorkerID: id, Err: err}:
			case <-ctx.Done():
				return
			}
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make(chan Job, 20)
	results := make(chan Result, 20)

	var wg sync.WaitGroup
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go worker(ctx, w, jobs, results, &wg)
	}

	for j := 1; j <= 20; j++ {
		jobs <- Job{ID: j}
	}
	close(jobs)

	// cancel halfway through the work
	time.AfterFunc(1*time.Second, func() {
		fmt.Println(">>> cancelling context")
		cancel()
	})

	go func() {
		wg.Wait()
		close(results)
	}()

	done := 0
	for r := range results {
		if r.Err != nil {
			fmt.Printf("job %2d (worker %d) failed: %v\n", r.JobID, r.WorkerID, r.Err)
			continue
		}
		done++
		fmt.Printf("job %2d done by worker %d\n", r.JobID, r.WorkerID)
	}
	fmt.Printf("completed %d of 20 jobs\n", done)
}
