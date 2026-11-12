package main

import (
	"fmt"
	"time"
)

// producer sends numbers into out, then closes it (sender closes).
func producer(out chan<- int) {
	for i := 1; i <= 5; i++ {
		fmt.Println("produced", i)
		out <- i
		time.Sleep(100 * time.Millisecond)
	}
	close(out)
}

// consumer reads until the channel is closed, then signals done.
func consumer(in <-chan int, done chan<- struct{}) {
	for v := range in {
		fmt.Println("    consumed", v)
		time.Sleep(300 * time.Millisecond) // deliberately slower than the producer
	}
	done <- struct{}{}
}

func main() {
	ch := make(chan int, 3) // unbuffered
	done := make(chan struct{})

	go producer(ch)
	go consumer(ch, done)

	<-done // wait for the consumer to finish
}
