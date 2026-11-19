package workers

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"

	"task_manager/Domain"
)

func TestMain(m *testing.M) { goleak.VerifyTestMain(m) }

// fakeUsecase implements Usecases.ReminderUsecase without any DB.
type fakeUsecase struct {
	mu        sync.Mutex
	pending   []Domain.Task  // returned once, then emptied (simulates the atomic claim)
	delivered map[string]int // title -> times delivered
	claimCall int
	twoTicks  chan struct{}
	once      sync.Once
}

func newFake(tasks []Domain.Task) *fakeUsecase {
	return &fakeUsecase{
		pending:   tasks,
		delivered: map[string]int{},
		twoTicks:  make(chan struct{}),
	}
}

func (f *fakeUsecase) ClaimDue(ctx context.Context, window time.Duration) ([]Domain.Task, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.claimCall++
	if f.claimCall >= 2 {
		f.once.Do(func() { close(f.twoTicks) })
	}
	out := f.pending
	f.pending = nil
	return out, nil
}

func (f *fakeUsecase) Deliver(ctx context.Context, t Domain.Task) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.delivered[t.Title]++
	return nil
}

func makeTasks(n int) []Domain.Task {
	tasks := make([]Domain.Task, n)
	for i := range tasks {
		tasks[i] = Domain.Task{Title: fmt.Sprintf("task-%d", i)}
	}
	return tasks
}

func TestWorkerPool_DeliversEveryJobExactlyOnce(t *testing.T) {
	fake := newFake(nil)
	pool := NewWorkerPool(4, fake, time.Second)

	jobs := make(chan Domain.Task, 20)
	pool.Start(jobs)
	for _, task := range makeTasks(20) {
		jobs <- task
	}
	close(jobs)
	pool.Wait()

	assert.Len(t, fake.delivered, 20)
	for title, n := range fake.delivered {
		assert.Equal(t, 1, n, title)
	}
}

func TestScheduler_StopsAndClosesChannelOnCancel(t *testing.T) {
	fake := newFake(nil)
	jobs := make(chan Domain.Task, 10)
	sched := NewScheduler(fake, 5*time.Millisecond, time.Minute, jobs)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { defer close(done); sched.Run(ctx) }()

	<-fake.twoTicks // proves the ticker is firing
	cancel()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return after cancel")
	}
	_, open := <-jobs
	assert.False(t, open, "jobs channel must be closed")
}

func TestSchedulerAndPool_EachReminderExactlyOnceAcrossTicks(t *testing.T) {
	fake := newFake(makeTasks(5))
	jobs := make(chan Domain.Task, 10)
	sched := NewScheduler(fake, 5*time.Millisecond, time.Minute, jobs)
	pool := NewWorkerPool(3, fake, time.Second)

	ctx, cancel := context.WithCancel(context.Background())
	pool.Start(jobs)
	done := make(chan struct{})
	go func() { defer close(done); sched.Run(ctx) }()

	<-fake.twoTicks // at least two ticks have happened
	cancel()
	<-done      // scheduler returned, channel closed
	pool.Wait() // workers drained everything

	assert.Len(t, fake.delivered, 5)
	for title, n := range fake.delivered {
		assert.Equal(t, 1, n, title)
	}
}

func TestScheduler_BlockedSendDoesNotLeakOnCancel(t *testing.T) {
	fake := newFake(makeTasks(5))
	jobs := make(chan Domain.Task) // unbuffered, nobody reading: send blocks
	sched := NewScheduler(fake, time.Hour, time.Minute, jobs)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { defer close(done); sched.Run(ctx) }()

	time.Sleep(20 * time.Millisecond) // let it block on the first send
	cancel()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("scheduler stuck on a blocked send after cancel")
	}
}
