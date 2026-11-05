package Repositories

import (
	"context"

	"task_manager/Domain"
)

// TaskRepository is the persistence contract for tasks.
// Implementations must return Domain.ErrNotFound when no task matches.
type TaskRepository interface {
	GetAll(ctx context.Context) ([]Domain.Task, error)
	GetByID(ctx context.Context, id string) (Domain.Task, error)
	Create(ctx context.Context, task Domain.Task) (Domain.Task, error)
	Update(ctx context.Context, id string, task Domain.Task) (Domain.Task, error)
	Delete(ctx context.Context, id string) error
}

// Day 18: add the Mongo implementation below this line.
