package Usecases

import (
	"context"
	"fmt"
	"strings"
	"time"

	"task_manager/Domain"
	"task_manager/Repositories"
)

// Caller is the authenticated user making the request.
type Caller struct {
	ID   string
	Role string
}

// IsAdmin reports whether the caller has the admin role.
func (c Caller) IsAdmin() bool { return c.Role == "admin" }

type callerKey struct{}

// WithCaller stores the authenticated caller in the context (called by the controller).
func WithCaller(ctx context.Context, c Caller) context.Context {
	return context.WithValue(ctx, callerKey{}, c)
}

// callerFrom reads the caller back; a missing caller is rejected, never treated as admin.
func callerFrom(ctx context.Context) (Caller, error) {
	c, ok := ctx.Value(callerKey{}).(Caller)
	if !ok || c.ID == "" {
		return Caller{}, fmt.Errorf("%w: missing authenticated user", Domain.ErrInvalidInput)
	}
	return c, nil
}

// TaskUsecase defines the business operations on tasks.
type TaskUsecase interface {
	GetAllTasks(ctx context.Context) ([]Domain.Task, error)
	ListTasks(ctx context.Context, f Repositories.TaskFilter) ([]Domain.Task, int, error)
	GetTaskByID(ctx context.Context, id string) (Domain.Task, error)
	CreateTask(ctx context.Context, task Domain.Task) (Domain.Task, error)
	UpdateTask(ctx context.Context, id string, task Domain.Task) (Domain.Task, error)
	DeleteTask(ctx context.Context, id string) error
}

type taskUsecase struct {
	repo Repositories.TaskRepository
}

// NewTaskUsecase builds a TaskUsecase that depends only on the repository interface.
func NewTaskUsecase(repo Repositories.TaskRepository) TaskUsecase {
	return &taskUsecase{repo: repo}
}

func validStatus(s string) bool {
	return s == "pending" || s == "in_progress" || s == "done"
}

// validateTask enforces the business rules for task fields.
func validateTask(t Domain.Task) error {
	if strings.TrimSpace(t.Title) == "" {
		return fmt.Errorf("%w: title is required", Domain.ErrInvalidInput)
	}
	if t.DueDate.IsZero() {
		return fmt.Errorf("%w: due date is required", Domain.ErrInvalidInput)
	}
	if !validStatus(t.Status) {
		return fmt.Errorf("%w: status must be pending, in_progress or done", Domain.ErrInvalidInput)
	}
	return nil
}

// loadOwned fetches a task and hides it (ErrNotFound) from users who do not own it.
// Admins can access every task. Returning NotFound avoids revealing that it exists.
func (u *taskUsecase) loadOwned(ctx context.Context, caller Caller, id string) (Domain.Task, error) {
	t, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return Domain.Task{}, err
	}
	if !caller.IsAdmin() && t.OwnerID != caller.ID {
		return Domain.Task{}, Domain.ErrNotFound
	}
	return t, nil
}

// GetAllTasks returns every task (used internally / by admin tools); the API uses ListTasks.
func (u *taskUsecase) GetAllTasks(ctx context.Context) ([]Domain.Task, error) {
	return u.repo.GetAll(ctx)
}

// ListTasks returns one page of the caller's tasks (admins see all owners).
func (u *taskUsecase) ListTasks(ctx context.Context, f Repositories.TaskFilter) ([]Domain.Task, int, error) {
	caller, err := callerFrom(ctx)
	if err != nil {
		return nil, 0, err
	}
	if f.Status != "" && !validStatus(f.Status) {
		return nil, 0, fmt.Errorf("%w: status must be pending, in_progress or done", Domain.ErrInvalidInput)
	}
	if !caller.IsAdmin() {
		f.OwnerID = caller.ID // normal users are always restricted to their own tasks
	}
	return u.repo.List(ctx, f.Normalized())
}

func (u *taskUsecase) GetTaskByID(ctx context.Context, id string) (Domain.Task, error) {
	caller, err := callerFrom(ctx)
	if err != nil {
		return Domain.Task{}, err
	}
	return u.loadOwned(ctx, caller, id)
}

func (u *taskUsecase) CreateTask(ctx context.Context, task Domain.Task) (Domain.Task, error) {
	caller, err := callerFrom(ctx)
	if err != nil {
		return Domain.Task{}, err
	}
	if err := validateTask(task); err != nil {
		return Domain.Task{}, err
	}
	task.OwnerID = caller.ID                     // owner comes from the token, never from the body
	task.ReminderStatus = Domain.ReminderPending // new tasks always start pending
	return u.repo.Create(ctx, task)
}

func (u *taskUsecase) UpdateTask(ctx context.Context, id string, task Domain.Task) (Domain.Task, error) {
	caller, err := callerFrom(ctx)
	if err != nil {
		return Domain.Task{}, err
	}
	if err := validateTask(task); err != nil {
		return Domain.Task{}, err
	}
	existing, err := u.loadOwned(ctx, caller, id)
	if err != nil {
		return Domain.Task{}, err
	}
	task.OwnerID = existing.OwnerID // ownership cannot be changed through update
	// The reminder status is server-owned: keep the stored value, but reset it
	// to pending when the due date moves to a future time.
	task.ReminderStatus = existing.ReminderStatus
	if !task.DueDate.Equal(existing.DueDate) && task.DueDate.After(time.Now()) {
		task.ReminderStatus = Domain.ReminderPending
	}
	return u.repo.Update(ctx, id, task)
}

func (u *taskUsecase) DeleteTask(ctx context.Context, id string) error {
	caller, err := callerFrom(ctx)
	if err != nil {
		return err
	}
	if _, err := u.loadOwned(ctx, caller, id); err != nil {
		return err
	}
	return u.repo.Delete(ctx, id)
}
