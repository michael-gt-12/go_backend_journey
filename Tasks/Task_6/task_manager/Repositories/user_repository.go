package Repositories

import (
	"context"

	"task_manager/Domain"
)

// UserRepository is the persistence contract for users.
// Implementations must return Domain.ErrNotFound when no user matches.
type UserRepository interface {
	Create(ctx context.Context, user Domain.User) (Domain.User, error)
	GetByUsername(ctx context.Context, username string) (Domain.User, error)
	GetByID(ctx context.Context, id string) (Domain.User, error)
	Count(ctx context.Context) (int64, error)
	UpdateRole(ctx context.Context, id string, role string) error
}

// Day 18: add the Mongo implementation below this line.
