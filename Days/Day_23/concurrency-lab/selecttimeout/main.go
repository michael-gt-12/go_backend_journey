package main

import (
	"fmt"
	"runtime"
	"time"
)

// slowWork pretends to do work that takes d, then sends a result.
func slowWork(d time.Duration) <-chan string {
	ch := make(chan string, 1) // buffered so the goroutine can always finish
	go func() {
		time.Sleep(d)
		ch <- fmt.Sprintf("finished after %v", d)
	}()
	return ch
}

// waitFor waits for the work to finish, but gives up after timeout.
func waitFor(name string, work <-chan string, timeout time.Duration) {
	select {
	case res := <-work:
		fmt.Printf("%s: %s\n", name, res)
	case <-time.After(timeout):
		fmt.Printf("%s: timed out after %v\n", name, timeout)
	}
}

func main() {
	waitFor("fast job", slowWork(200*time.Millisecond), 1*time.Second)
	waitFor("slow job", slowWork(2*time.Second), 1*time.Second)

	time.Sleep(2 * time.Second) // give the slow goroutine time to reach its send
	fmt.Println("goroutines still alive:", runtime.NumGoroutine())
}
