# Concurrency Design: Background Reminders (Task 8)

The API runs a reminder system inside the same process as the HTTP server. A scheduler finds tasks that are due soon, hands them to a pool of workers through a channel, and the workers deliver a reminder through a `Notifier` (for now, a log line).

## 1. Flow

```
                 every REMINDER_INTERVAL
   +-----------+   ClaimDue(window)    +-----------------+
   | Scheduler | --------------------> | ReminderUsecase |--> TaskRepository.ClaimDueTasks
   | (ticker)  |                       +-----------------+      (atomic pending -> queued)
   +-----+-----+
         | sends claimed tasks
         v
   +---------------------------+
   |  jobs channel (buffered)  |   size = REMINDER_QUEUE_SIZE
   +------+--------+-----------+
          |        |        |
          v        v        v
     +--------+ +--------+ +--------+
     | Worker | | Worker | | Worker |   N = REMINDER_WORKERS
     +---+----+ +---+----+ +---+----+
         |          |          |
         +----------+----------+
                    v
        ReminderUsecase.Deliver(task)
         -> Notifier.Notify (up to 3 attempts, 2s apart)
         -> status = sent   (or failed after the last attempt)
```

| Piece            | File                              | Responsibility                                                                                            |
| ---------------- | --------------------------------- | --------------------------------------------------------------------------------------------------------- |
| Scheduler        | `Delivery/workers/scheduler.go`   | Ticks, calls `ClaimDue`, sends tasks to the channel. The **only sender**, so the only one that closes it. |
| Worker pool      | `Delivery/workers/worker_pool.go` | Starts N goroutines that `range` over the channel and call `Deliver`.                                     |
| Reminder usecase | `Usecases/reminder_usecases.go`   | Business logic: claiming, retries, status updates.                                                        |
| Notifier         | `Infrastructure/log_notifier.go`  | Delivers the reminder (logs ID, title, due date).                                                         |
| Repository       | `Repositories/task_repository.go` | `ClaimDueTasks` (atomic) and `SetReminderStatus`.                                                         |

The scheduler and the pool only coordinate. All business rules stay in the usecase, so the Clean Architecture import rules still hold.

## 2. Reminder status lifecycle

```
pending --(claimed by scheduler)--> queued --(Notify ok)--> sent
                                       |
                                       +--(all attempts failed)--> failed
```

- New tasks start as `pending`.
- Changing a tasks due date to a time in the future through `PUT /tasks/:id` resets it to `pending`.
- `ClaimDueTasks` moves tasks from `pending` to `queued` atomically (`FindOneAndUpdate` in a loop up to `limit`), with the filter `reminder_status = pending AND due_date <= dueBefore AND status != done`. Two claims running at the same time never return the same task, so a reminder is never sent twice.
- Overdue tasks match the filter too, so nothing is missed if the app was down.

## 3. Shutdown order

On SIGINT or SIGTERM (`signal.NotifyContext`), `main.go` does these steps in order and logs each one:

1. **Signal received.**
2. **`http.Server.Shutdown`** with a 10 second timeout: stop taking new requests and let running requests finish.
3. **Cancel the scheduler context.** The scheduler stops its ticker, returns, and closes the jobs channel (via `defer close`).
4. **`WorkerPool.Wait()`**: workers finish the jobs already in the channel and exit when it is closed and drained.
5. **Disconnect MongoDB**, then exit.

Why this order:

- HTTP stops first so no new work arrives while we are shutting down.
- The scheduler stops before the workers so nothing new is queued. Only the sender may close the channel, and it closes it exactly once.
- Workers stop last so reminders already claimed are not lost.
- MongoDB closes last because the workers still need it to save the `sent` or `failed` status.

Two details that make this work:

- The scheduler runs on its **own** context, not the signal context. That lets `main.go` choose when it stops.
- Each job gets a context created from `context.Background()` with a timeout (5 seconds), **not** from the shutdown context. A job already in progress can therefore finish while the app shuts down. If it were derived from the signal context, Ctrl+C would cancel the running job.

## 4. Configuration

| Variable              | Default | Meaning                                                                                 |
| --------------------- | ------- | --------------------------------------------------------------------------------------- |
| `REMINDER_INTERVAL`   | `30s`   | Time between scheduler ticks (one tick also runs immediately on start)                  |
| `REMINDER_WINDOW`     | `10m`   | A task is "due soon" if its due date is within this window (overdue tasks are included) |
| `REMINDER_WORKERS`    | `3`     | Number of worker goroutines                                                             |
| `REMINDER_QUEUE_SIZE` | `100`   | Buffer size of the jobs channel                                                         |

Invalid values stop the app at startup with a clear message. Delivery uses up to 3 attempts with 2 seconds between them, and the wait stops immediately if the context is cancelled.

## 5. Manual verification

1. Run with `REMINDER_INTERVAL=5s` and `REMINDER_WINDOW=3m`.
2. Create two tasks due in about 2 minutes.
3. Watch the logs and MongoDB Compass: each task goes `pending`, `queued`, `sent`, and each reminder is logged exactly once, by one of the workers.
4. Press Ctrl+C while a reminder is being delivered. The logs show `signal received`, `http server stopped`, `scheduler stopped, jobs channel closed`, the workers exiting, `workers finished`, and `bye`.

## 6. Tests

| Test                                                      | What it proves                                                                                    |
| --------------------------------------------------------- | ------------------------------------------------------------------------------------------------- |
| `TestWorkerPool_DeliversEveryJobExactlyOnce`              | 20 jobs sent, channel closed, `Wait()` returns, every job delivered once                          |
| `TestScheduler_StopsAndClosesChannelOnCancel`             | `Run` returns after cancel and the channel is closed                                              |
| `TestSchedulerAndPool_EachReminderExactlyOnceAcrossTicks` | Across several ticks each reminder is delivered exactly once                                      |
| `TestScheduler_BlockedSendDoesNotLeakOnCancel`            | A scheduler blocked on a full channel still stops on cancel                                       |
| `ReminderUsecase` suite                                   | Success on first or second attempt, failure after all attempts, context cancelled during the wait |

- `goleak.VerifyTestMain` fails the package if any goroutine is left running.
- Synchronisation uses channels, not `time.Sleep` assertions.
- `go test -race` does not work on Windows without cgo, so the race detector runs in Docker (`docker run --rm -v "${PWD}:/app" -w /app golang:latest go test -race -count=20 ./Delivery/workers/...`) and in GitHub Actions.

## 7. Known limitation

If the app crashes after a task is claimed (`queued`) but before it is delivered, that task stays `queued` and is never reminded. The same happens to tasks the scheduler had claimed but could not enqueue before shutdown (it logs how many).

Production systems solve this with a **visibility timeout** (a claim expires after some time and the task becomes claimable again) or a **recovery step** that resets tasks stuck in `queued` back to `pending` at startup. Task 9 revisits this when claiming moves to Postgres with `FOR UPDATE SKIP LOCKED`.
