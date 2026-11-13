package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// producer enqueues one job per second until ctx is cancelled.
// It is the only sender, so it is the one that closes the channel.
func producer(ctx context.Context, jobs chan<- int) {
	defer close(jobs)
	tk := time.NewTicker(1 * time.Second)
	defer tk.Stop() // always stop the ticker

	id := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Println("producer: ctx done, ticker stopped")
			return
		case <-tk.C:
			id++
			select {
			case jobs <- id:
				fmt.Println("producer: enqueued job", id)
			case <-ctx.Done():
				return
			}
		}
	}
}

func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("  worker %d processing job %d\n", id, job)
		time.Sleep(300 * time.Millisecond)
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5500*time.Millisecond)
	defer cancel()

	jobs := make(chan int, 10)

	var wg sync.WaitGroup
	for w := 1; w <= 2; w++ {
		wg.Add(1)
		go worker(w, jobs, &wg)
	}

	go producer(ctx, jobs)

	wg.Wait() // workers exit once the producer closes jobs
	fmt.Println("all workers finished")
}
