# Day 23 Notes: Concurrency Lab 1

## 1. Channels vs mutexes

- Use a **channel** when goroutines need to pass data (or ownership of data) to each other: results from workers, jobs to workers, "I'm done" signals.
- Use a **mutex** when several goroutines need to read/write the same piece of shared state, like the counter in the race demo, and passing it around would be awkward.
- Rule of thumb: channels for communication and coordination, mutexes for protecting state.

## 2. What is a data race?

- A data race happens when two goroutines access the same variable at the same time and at least one of them writes, with no synchronization.
- `counter++` looks like one step but is three: read the value, add 1, write it back. If two goroutines both read 57, both write 58, and one increment is lost.
- Proof from my runs: 1000 goroutines incrementing gave 982, 972, 973, 995 and 991 instead of 1000.
- `go run -race` printed `WARNING: DATA RACE` twice (read/write and write/write conflicts), both pointing at `counter++` on line 16, and ended with `Found 2 data race(s)` and `exit status 66`.
- The race detector only finds races in code that actually runs, and a program can print the right answer while still having a race.
- Fix: wrap `counter++` in `mu.Lock()` / `mu.Unlock()`. Goroutines take turns, no increment is overwritten, and the result is always 1000 with no race warning.

## 3. What is a goroutine leak?

- A goroutine that never exits because it is blocked forever (for example, sending on a channel nobody will ever receive from).
- Example: in the `select` timeout demo, if the work channel is unbuffered and the timeout wins, `waitFor` returns and nobody reads the channel. The slow goroutine wakes up later, tries to send, and blocks forever.
- In a long-running server this means memory grows with every timed-out request.
- Fix: give the channel a buffer of 1 so the send always succeeds, or use a context so the goroutine can be cancelled.
- Check: `runtime.NumGoroutine()` shows an extra goroutine in the leaky version. Tomorrow I'll automate this with `goleak`.
- Rule: every goroutine I start must have a clear way to stop.

## 4. Buffered vs unbuffered channels

- **Unbuffered** (`make(chan int)`): a send blocks until a receiver takes the value, so sender and receiver meet at every handoff. In my producer/consumer run the output strictly alternated (produced, consumed, produced, consumed) because the fast producer had to wait for the slow consumer.
- **Buffered** (`make(chan int, 3)`): a send only blocks when the buffer is full. The producer ran ahead by up to 3 items, which is why `produced 2` and `produced 3` appeared before `consumed 2`.
- When the buffer is full, the producer blocks again. That is **backpressure**, and it is how the jobs queue will behave in the reminder scheduler (Task 8).

## 5. Who closes a channel, and why?

- Only the **sender** closes a channel, because only the sender knows when no more values are coming.
- Sending on a closed channel panics, so a receiver closing it could crash a sender that is still sending.
- Receiving from a closed channel is safe: the receiver drains what is left, then `for range` ends.
- If nobody closes it, `for range` blocks forever waiting for more values (Go reports `all goroutines are asleep - deadlock!` when nothing else can run).
- With several senders, close the channel from one place only: a separate goroutine that calls `wg.Wait()` and then `close(ch)`, so it runs while `main` is already receiving.

## 6. WaitGroup in one paragraph

- `wg.Add(1)` before starting each goroutine, `defer wg.Done()` inside it, and `wg.Wait()` to block until the counter reaches 0.
- Without `wg.Wait()`, `main` exits before the goroutines finish and the program ends early.
- Pass the loop variable as an argument to the goroutine so each one gets its own copy.

## 7. Timing result

- 5 URLs fetched sequentially: 4.494s
- The same 5 fetched concurrently with a WaitGroup: 1.164s (channel version: 1.163s)
- Concurrent time is about the slowest single request, while sequential time is the sum of all requests.

## 8. Windows note

- The race detector needs cgo and a C compiler, so `go run -race` doesn't work in my local Windows setup. I ran it in a Docker golang container instead:
  `docker run --rm -v "${PWD}:/app" -w /app golang:latest go run -race ./race`
