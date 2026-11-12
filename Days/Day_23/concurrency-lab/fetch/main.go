package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var urls = []string{
	"https://go.dev",
	"https://github.com",
	"https://www.wikipedia.org",
	"https://httpbin.org/delay/1",
	"https://example.com",
}

func fetch(url string) (int, error) {
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func runSequential() {
	start := time.Now()
	for _, u := range urls {
		t := time.Now()
		status, err := fetch(u)
		if err != nil {
			fmt.Printf("%-30s error: %v\n", u, err)
			continue
		}
		fmt.Printf("%-30s %d (%v)\n", u, status, time.Since(t).Round(time.Millisecond))
	}
	fmt.Println("sequential total:", time.Since(start).Round(time.Millisecond))
}

func runConcurrent() {
	start := time.Now()
	var wg sync.WaitGroup
	for _, u := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			t := time.Now()
			status, err := fetch(u)
			if err != nil {
				fmt.Printf("%-30s error: %v\n", u, err)
				return
			}
			fmt.Printf("%-30s %d (%v)\n", url, status, time.Since(t).Round(time.Millisecond))
		}(u)

	}
	wg.Wait()
	fmt.Println("concurrent total:", time.Since(start).Round(time.Millisecond))
}

// result carries the outcome of one fetch between goroutines.
type result struct {
	url      string
	status   int
	err      error
	duration time.Duration
}

func runWithChannel() {
	start := time.Now()
	results := make(chan result, len(urls)) // buffered: senders never block
	var wg sync.WaitGroup

	for _, u := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			t := time.Now()
			status, err := fetch(url)
			results <- result{url, status, err, time.Since(t)} // send
		}(u)
	}

	// Closer goroutine: close the channel once all senders are done.
	go func() {
		wg.Wait()
		close(results)
	}()

	// Receive until the channel is closed and drained.
	for r := range results {
		if r.err != nil {
			fmt.Printf("%-30s error: %v\n", r.url, r.err)
			continue
		}
		fmt.Printf("%-30s %d (%v)\n", r.url, r.status, r.duration.Round(time.Millisecond))
	}

	fmt.Println("channel total:", time.Since(start).Round(time.Millisecond))
}

func main() {
	fmt.Println("--- sequential ---")
	runSequential()

	fmt.Println("--- concurrent (WaitGroup) ---")
	runConcurrent()

	fmt.Println("--- concurrent (channel) ---")
	runWithChannel()
}
