package Usecases

import (
	"context"
	"fmt"
	"time"

	"task_manager/Domain"
	"task_manager/Repositories"
)

// Notifier delivers one reminder. Implementations live in Infrastructure.
type Notifier interface {
	Notify(ctx context.Context, task Domain.Task) error
}

// ReminderUsecase is what the scheduler and the worker pool call.
type ReminderUsecase interface {
	ClaimDue(ctx context.Context, window time.Duration) ([]Domain.Task, error)
	Deliver(ctx context.Context, task Domain.Task) error
}

const (
	// DefaultMaxAttempts is used when the constructor gets a value below 1.
	DefaultMaxAttempts = 3
	// claimBatchSize is how many tasks one ClaimDue call may claim.
	claimBatchSize = 50
)

type reminderUsecase struct {
	repo        Repositories.TaskRepository
	notifier    Notifier
	maxAttempts int
	retryDelay  time.Duration
}

// NewReminderUsecase builds the usecase. maxAttempts and retryDelay are
// parameters so tests can use a few attempts and a few milliseconds.
func NewReminderUsecase(repo Repositories.TaskRepository, notifier Notifier, maxAttempts int, retryDelay time.Duration) ReminderUsecase {
	if maxAttempts < 1 {
		maxAttempts = DefaultMaxAttempts
	}
	return &reminderUsecase{repo: repo, notifier: notifier, maxAttempts: maxAttempts, retryDelay: retryDelay}
}

// ClaimDue atomically moves tasks due within the window (and overdue ones)
// from pending to queued and returns them. If the repository fails part-way,
// the tasks already claimed are still returned together with the error so the
// caller can deliver them instead of leaving them stuck in "queued".
func (u *reminderUsecase) ClaimDue(ctx context.Context, window time.Duration) ([]Domain.Task, error) {
	return u.repo.ClaimDueTasks(ctx, time.Now().Add(window), claimBatchSize)
}

// Deliver tries the notifier up to maxAttempts times, waiting retryDelay
// between attempts, then records sent or failed. The wait stops at once when
// ctx is cancelled; in that case the status is left untouched (the context is
// already dead, so the write would fail) and ctx's error is returned.
func (u *reminderUsecase) Deliver(ctx context.Context, task Domain.Task) error {
	var notifyErr error
	for attempt := 1; attempt <= u.maxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("deliver task %s: %w", task.ID, err)
		}
		if notifyErr = u.notifier.Notify(ctx, task); notifyErr == nil {
			if err := u.repo.SetReminderStatus(ctx, task.ID, Domain.ReminderSent); err != nil {
				return fmt.Errorf("mark task %s sent: %w", task.ID, err)
			}
			return nil
		}
		if attempt == u.maxAttempts {
			break
		}
		select {
		case <-time.After(u.retryDelay):
		case <-ctx.Done():
			return fmt.Errorf("deliver task %s: %w", task.ID, ctx.Err())
		}
	}

	if err := u.repo.SetReminderStatus(ctx, task.ID, Domain.ReminderFailed); err != nil {
		return fmt.Errorf("mark task %s failed: %w (last notify error: %v)", task.ID, err, notifyErr)
	}
	return fmt.Errorf("deliver task %s: gave up after %d attempts: %w", task.ID, u.maxAttempts, notifyErr)
}
