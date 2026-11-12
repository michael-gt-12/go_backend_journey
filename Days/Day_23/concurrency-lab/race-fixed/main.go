package main

import (
	"fmt"
	"sync"
)

func main() {
	counter := 0
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++ // only one goroutine can be here at a time
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Println("expected: 1000, got:", counter)
}
