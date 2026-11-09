package Domain

import (
	"errors"
	"time"
)

// Roles a user can have.
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// Task is the core task entity.
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
}

// User is the core user entity. Password holds the bcrypt hash only.
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Role     string `json:"role"`
}

// Domain errors. Controllers map these to HTTP status codes.
var (
	ErrNotFound           = errors.New("not found")
	ErrInvalidInput       = errors.New("invalid input")
	ErrUsernameTaken      = errors.New("username already taken")
	ErrInvalidCredentials = errors.New("invalid username or password")
)
