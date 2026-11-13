# Day 24 notes: worker pools, context and tickers

## What I built

- `workerpool`: 5 workers read `Job`s from a buffered channel and send `Result`s on a second channel. A `sync.WaitGroup` plus one goroutine closes `results` when every worker is done, so `main` can `range` over it.
- `ctxcancel`: workers take a `context.Context`. Each job gets `context.WithTimeout` (500ms), and the whole run is cancelled after 1 second.
- `ticker`: a producer enqueues one job per second with `time.Ticker` until the context ends. Then it stops the ticker and closes the jobs channel.
- `leak`: a small pool tested with `goleak.VerifyNone`, plus a deliberate leak that goleak catches.

## Worker pool

- A channel is a shared pipe, not a snapshot. Workers started before any job exists simply block on `range jobs` until a value arrives, and `range` ends when the channel is closed and drained.
- Only the sender closes a channel. Workers must never close `jobs`.
- `results` can only be closed after all workers have finished, so a goroutine does `wg.Wait(); close(results)`.
- Which worker gets which job is not predictable. Job order in the output is not guaranteed either.

## Context

- `ctx.Done()` is a channel that closes on cancel or timeout. `ctx.Err()` says why: `context.Canceled` or `context.DeadlineExceeded`.
- A worker has to check `ctx.Done()` itself, in a `select` next to its other channel operations. Cancelling does nothing on its own.
- Always `defer cancel()` after `WithTimeout` or `WithCancel`.
- If two `select` cases are ready, Go picks one at random. After a cancel, a worker may still send a result or may return silently. This is why my local and Docker runs printed different lines after the cancel.
- A `defer` in the worker (for example a log line) runs on every return path, so shutdown logging is reliable.

## Ticker

- `time.NewTicker(d)` sends on `tk.C` every `d`. Always `defer tk.Stop()`, or the ticker keeps running.
- Pattern: `select { case <-ctx.Done(): return; case <-tk.C: ... }`.
- The producer is the only sender, so it closes the jobs channel (`defer close(jobs)`).

## goroutine leaks and goleak

- A goroutine leak is a goroutine that never exits, for example one blocked on a receive with no sender, or on a `nil` channel.
- Every goroutine I start needs a way to stop: ctx cancel, a closed channel, or a WaitGroup that ends it.
- `goleak.VerifyNone(t)` fails the test if goroutines are still alive when it runs. Its output shows the stack of the leaked goroutine.
- `goleak.Find()` returns an error instead of failing the test, which lets a test assert that a leak is detected.

## Channels vs mutexes (from Day 23, still true today)

- Channels: passing work or results between goroutines, and signalling (done, cancel).
- Mutex: protecting shared state that several goroutines read or write.

## Commands

```
go run ./workerpool
go run ./ctxcancel
go run ./ticker
go test ./leak -v
$env:LEAK_DEMO="1"; go test ./leak -run TestDeliberateLeak -v   # fails on purpose
docker run --rm -v "${PWD}:/app" -w /app golang:latest go run -race ./workerpool
docker run --rm -v "${PWD}:/app" -w /app golang:latest go test -race ./leak -v
```

## Why this matters for Task 8

The reminder system is a ticker (scheduler), a channel and a worker pool. It uses a context for shutdown, and goleak to prove nothing leaks.
