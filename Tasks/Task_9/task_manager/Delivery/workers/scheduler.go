package workers

import (
	"context"
	"log"
	"time"

	"task_manager/Domain"
	"task_manager/Usecases"
)

// Scheduler periodically claims due tasks and feeds them to the jobs channel.
// It is the only sender on the channel, so it is the only one allowed to close it.
type Scheduler struct {
	uc       Usecases.ReminderUsecase
	interval time.Duration
	window   time.Duration
	jobs     chan<- Domain.Task
}

func NewScheduler(uc Usecases.ReminderUsecase, interval, window time.Duration,
	jobs chan<- Domain.Task) *Scheduler {
	return &Scheduler{uc: uc, interval: interval, window: window, jobs: jobs}
}

// Run blocks until ctx is cancelled, then closes jobs and returns.
func (s *Scheduler) Run(ctx context.Context) {
	defer close(s.jobs) // exactly once, runs after the loop exits
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	log.Printf("scheduler: started (interval=%s window=%s)", s.interval, s.window)
	s.tick(ctx) // one tick immediately on start

	for {
		select {
		case <-ctx.Done():
			log.Println("scheduler: stopping, closing jobs channel")
			return
		case <-ticker.C:
			s.tick(ctx)
		}
	}
}

func (s *Scheduler) tick(ctx context.Context) {
	tasks, err := s.uc.ClaimDue(ctx, s.window)
	if err != nil {
		log.Printf("scheduler: claim failed: %v", err)
		return
	}
	for i, t := range tasks {
		select {
		case s.jobs <- t: // waits here if the channel is full
		case <-ctx.Done():
			// Claimed but not enqueued: never drop silently. These stay "queued"

			log.Printf("scheduler: shutdown with %d claimed task(s) not enqueued", len(tasks)-i)
			return
		}
	}
}
