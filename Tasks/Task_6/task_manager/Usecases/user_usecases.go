package Usecases

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"task_manager/Domain"
	"task_manager/Repositories"
)

// PasswordService hashes and verifies passwords (implemented in Infrastructure).
type PasswordService interface {
	Hash(password string) (string, error)
	Compare(hash, password string) bool
}

// TokenService issues JWTs for a user (implemented in Infrastructure).
type TokenService interface {
	Generate(user Domain.User) (string, error)
}

// UserUsecase defines the business operations on users.
type UserUsecase interface {
	Register(ctx context.Context, username, password string) (Domain.User, error)
	Login(ctx context.Context, username, password string) (string, error)
	Promote(ctx context.Context, id string) error
}

type userUsecase struct {
	repo      Repositories.UserRepository
	passwords PasswordService
	tokens    TokenService
}

// NewUserUsecase wires the usecase with its dependencies (all interfaces).
func NewUserUsecase(repo Repositories.UserRepository, p PasswordService, t TokenService) UserUsecase {
	return &userUsecase{repo: repo, passwords: p, tokens: t}
}

// Register validates input, rejects duplicates, hashes the password and stores the user.
// The first registered user becomes admin; everyone else is a normal user.
func (u *userUsecase) Register(ctx context.Context, username, password string) (Domain.User, error) {
	username = strings.TrimSpace(username)
	if len(username) < 3 {
		return Domain.User{}, fmt.Errorf("%w: username must be at least 3 characters", Domain.ErrInvalidInput)
	}
	if len(password) < 6 {
		return Domain.User{}, fmt.Errorf("%w: password must be at least 6 characters", Domain.ErrInvalidInput)
	}

	_, err := u.repo.GetByUsername(ctx, username)
	if err == nil {
		return Domain.User{}, Domain.ErrUsernameTaken
	}
	if !errors.Is(err, Domain.ErrNotFound) {
		return Domain.User{}, err
	}

	hash, err := u.passwords.Hash(password)
	if err != nil {
		return Domain.User{}, err
	}

	count, err := u.repo.Count(ctx)
	if err != nil {
		return Domain.User{}, err
	}
	role := Domain.RoleUser
	if count == 0 {
		role = Domain.RoleAdmin
	}

	created, err := u.repo.Create(ctx, Domain.User{Username: username, Password: hash, Role: role})
	if err != nil {
		return Domain.User{}, err
	}
	created.Password = "" // never return the hash
	return created, nil
}

// Login verifies credentials and returns a JWT. The error is the same for a
// wrong password and an unknown user, so attackers cannot probe usernames.
func (u *userUsecase) Login(ctx context.Context, username, password string) (string, error) {
	user, err := u.repo.GetByUsername(ctx, strings.TrimSpace(username))
	if err != nil {
		if errors.Is(err, Domain.ErrNotFound) {
			return "", Domain.ErrInvalidCredentials
		}
		return "", err
	}
	if !u.passwords.Compare(user.Password, password) {
		return "", Domain.ErrInvalidCredentials
	}
	return u.tokens.Generate(user)
}

// Promote makes a user an admin.
func (u *userUsecase) Promote(ctx context.Context, id string) error {
	if _, err := u.repo.GetByID(ctx, id); err != nil {
		return err
	}
	return u.repo.UpdateRole(ctx, id, Domain.RoleAdmin)
}
