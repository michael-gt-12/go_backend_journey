package Repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"task_manager/Domain"
)

const (
	pgUniqueViolation = "23505" // unique_violation
	pgInvalidTextRepr = "22P02" // invalid_text_representation (e.g. malformed UUID)
)

type postgresUserRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresUserRepository returns a PostgreSQL-backed UserRepository.
func NewPostgresUserRepository(pool *pgxpool.Pool) UserRepository {
	return &postgresUserRepository{pool: pool}
}

// isPgCode reports whether err is a PostgreSQL error with the given SQLSTATE code.
func isPgCode(err error, code string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == code
}

func (r *postgresUserRepository) Create(ctx context.Context, user Domain.User) (Domain.User, error) {
	const q = `
		INSERT INTO users (username, password_hash, role)
		VALUES ($1, $2, $3)
		RETURNING id::text`

	err := r.pool.QueryRow(ctx, q, user.Username, user.Password, user.Role).Scan(&user.ID)
	if err != nil {
		if isPgCode(err, pgUniqueViolation) {
			return Domain.User{}, Domain.ErrUsernameTaken
		}
		return Domain.User{}, err
	}
	return user, nil
}

func (r *postgresUserRepository) GetByUsername(ctx context.Context, username string) (Domain.User, error) {
	const q = `SELECT id::text, username, password_hash, role FROM users WHERE username = $1`

	var u Domain.User
	err := r.pool.QueryRow(ctx, q, username).Scan(&u.ID, &u.Username, &u.Password, &u.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Domain.User{}, Domain.ErrNotFound
		}
		return Domain.User{}, err
	}
	return u, nil
}

func (r *postgresUserRepository) GetByID(ctx context.Context, id string) (Domain.User, error) {
	const q = `SELECT id::text, username, password_hash, role FROM users WHERE id = $1::uuid`

	var u Domain.User
	err := r.pool.QueryRow(ctx, q, id).Scan(&u.ID, &u.Username, &u.Password, &u.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || isPgCode(err, pgInvalidTextRepr) {
			return Domain.User{}, Domain.ErrNotFound
		}
		return Domain.User{}, err
	}
	return u, nil
}

func (r *postgresUserRepository) Count(ctx context.Context) (int64, error) {
	var n int64
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&n); err != nil {
		return 0, err
	}
	return n, nil
}

func (r *postgresUserRepository) UpdateRole(ctx context.Context, id string, role string) error {
	tag, err := r.pool.Exec(ctx, `UPDATE users SET role = $1 WHERE id = $2::uuid`, role, id)
	if err != nil {
		if isPgCode(err, pgInvalidTextRepr) {
			return Domain.ErrNotFound
		}
		return err
	}
	if tag.RowsAffected() == 0 {
		return Domain.ErrNotFound
	}
	return nil
}
