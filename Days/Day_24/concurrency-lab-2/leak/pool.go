package leak

import (
	"context"
	"sync"
)

// StartWorkers launches n workers that stop when ctx is cancelled.
func StartWorkers(ctx context.Context, n int, jobs <-chan int) *sync.WaitGroup {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case _, ok := <-jobs:
					if !ok {
						return
					}
				}
			}
		}()
	}
	return &wg
}

// Leaky starts a goroutine that blocks until stop is closed.
// Passing a nil channel makes it block forever: a real leak.
func Leaky(stop <-chan struct{}) {
	go func() { <-stop }()
}
