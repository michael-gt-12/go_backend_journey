package main

import (
	"fmt"
	"sync"
)

func main() {
	counter := 0
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++ // READ, ADD, WRITE: three steps, not atomic
		}()
	}

	wg.Wait()
	fmt.Println("expected: 1000, got:", counter)
}
