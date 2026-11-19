package Usecases

import (
	"context"
	"fmt"
	"strings"
	"time"

	"task_manager/Domain"
	"task_manager/Repositories"
)

// TaskUsecase defines the business operations on tasks.
type TaskUsecase interface {
	GetAllTasks(ctx context.Context) ([]Domain.Task, error)
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

func (u *taskUsecase) GetAllTasks(ctx context.Context) ([]Domain.Task, error) {
	return u.repo.GetAll(ctx)
}

func (u *taskUsecase) GetTaskByID(ctx context.Context, id string) (Domain.Task, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *taskUsecase) CreateTask(ctx context.Context, task Domain.Task) (Domain.Task, error) {
	if err := validateTask(task); err != nil {
		return Domain.Task{}, err
	}
	task.ReminderStatus = Domain.ReminderPending // new tasks always start pending
	return u.repo.Create(ctx, task)
}

func (u *taskUsecase) UpdateTask(ctx context.Context, id string, task Domain.Task) (Domain.Task, error) {
	if err := validateTask(task); err != nil {
		return Domain.Task{}, err
	}
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return Domain.Task{}, err
	}
	// The reminder status is server-owned: keep the stored value, but reset it
	// to pending when the due date moves to a future time.
	task.ReminderStatus = existing.ReminderStatus
	if !task.DueDate.Equal(existing.DueDate) && task.DueDate.After(time.Now()) {
		task.ReminderStatus = Domain.ReminderPending
	}
	return u.repo.Update(ctx, id, task)
}

func (u *taskUsecase) DeleteTask(ctx context.Context, id string) error {
	if _, err := u.repo.GetByID(ctx, id); err != nil {
		return err
	}
	return u.repo.Delete(ctx, id)
}
