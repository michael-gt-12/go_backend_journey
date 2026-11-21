package Repositories

import (
	"context"
	"time"

	"task_manager/Domain"
)

// TaskFilter describes a paginated, filtered list query.
type TaskFilter struct {
	OwnerID string // empty means all owners (admin)
	Status  string // empty means any status
	Page    int
	Limit   int
}

// Normalized returns the filter with defaults applied (page 1, limit 10, max 100).
func (f TaskFilter) Normalized() TaskFilter {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 {
		f.Limit = 10
	}
	if f.Limit > 100 {
		f.Limit = 100
	}
	return f
}

// Offset returns the number of rows to skip for the current page.
func (f TaskFilter) Offset() int {
	f = f.Normalized()
	return (f.Page - 1) * f.Limit
}

// TaskRepository is the persistence contract for tasks.
// Implementations must return Domain.ErrNotFound when no task matches.
type TaskRepository interface {
	GetAll(ctx context.Context) ([]Domain.Task, error)
	// List returns one page of tasks ordered by due_date ASC, id ASC, plus the total match count.
	List(ctx context.Context, f TaskFilter) ([]Domain.Task, int, error)
	GetByID(ctx context.Context, id string) (Domain.Task, error)
	Create(ctx context.Context, task Domain.Task) (Domain.Task, error)
	Update(ctx context.Context, id string, task Domain.Task) (Domain.Task, error)
	Delete(ctx context.Context, id string) error

	// ClaimDueTasks atomically moves up to limit tasks that are due on or before
	// dueBefore from pending to queued and returns only the tasks it moved.
	ClaimDueTasks(ctx context.Context, dueBefore time.Time, limit int) ([]Domain.Task, error)
	// SetReminderStatus sets reminder_status; returns Domain.ErrNotFound if no task matches.
	SetReminderStatus(ctx context.Context, id string, status string) error
}
