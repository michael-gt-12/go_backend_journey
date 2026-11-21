package Repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"task_manager/Domain"
)

// taskColumns lists the columns in the order scanTask reads them.
// UUIDs are cast to text so they scan straight into Go strings.
const taskColumns = `id::text, user_id::text, title, description, due_date, status, reminder_status`

type postgresTaskRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresTaskRepository returns a PostgreSQL-backed TaskRepository.
func NewPostgresTaskRepository(pool *pgxpool.Pool) TaskRepository {
	return &postgresTaskRepository{pool: pool}
}

// scanTask reads one row (from pgx.Row or pgx.Rows) into a Domain.Task.
func scanTask(row pgx.Row) (Domain.Task, error) {
	var t Domain.Task
	err := row.Scan(&t.ID, &t.OwnerID, &t.Title, &t.Description, &t.DueDate, &t.Status, &t.ReminderStatus)
	return t, err
}

// mapPgErr translates pgx/PostgreSQL errors into domain errors without leaking SQL details.
func mapPgErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return Domain.ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "22P02": // invalid_text_representation: e.g. malformed UUID
			return fmt.Errorf("%w: invalid id", Domain.ErrInvalidInput)
		case "23503": // foreign_key_violation: owner does not exist
			return fmt.Errorf("%w: owner does not exist", Domain.ErrInvalidInput)
		case "23514": // check_violation: bad status value
			return fmt.Errorf("%w: invalid field value", Domain.ErrInvalidInput)
		}
	}
	return err
}

func (r *postgresTaskRepository) GetAll(ctx context.Context) ([]Domain.Task, error) {
	rows, err := r.pool.Query(ctx, `SELECT `+taskColumns+` FROM tasks ORDER BY due_date ASC, id ASC`)
	if err != nil {
		return nil, mapPgErr(err)
	}
	defer rows.Close()
	return collectTasks(rows)
}

func collectTasks(rows pgx.Rows) ([]Domain.Task, error) {
	tasks := []Domain.Task{}
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, mapPgErr(err)
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, mapPgErr(err)
	}
	return tasks, nil
}

func (r *postgresTaskRepository) List(ctx context.Context, f TaskFilter) ([]Domain.Task, int, error) {
	f = f.Normalized()

	// Only fixed SQL fragments are concatenated; all values go through $n placeholders.
	where := " WHERE TRUE"
	args := []any{}
	if f.OwnerID != "" {
		args = append(args, f.OwnerID)
		where += fmt.Sprintf(" AND user_id = $%d", len(args))
	}
	if f.Status != "" {
		args = append(args, f.Status)
		where += fmt.Sprintf(" AND status = $%d", len(args))
	}

	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM tasks`+where, args...).Scan(&total); err != nil {
		return nil, 0, mapPgErr(err)
	}

	args = append(args, f.Limit, f.Offset())
	query := `SELECT ` + taskColumns + ` FROM tasks` + where +
		fmt.Sprintf(` ORDER BY due_date ASC, id ASC LIMIT $%d OFFSET $%d`, len(args)-1, len(args))

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, mapPgErr(err)
	}
	defer rows.Close()
	tasks, err := collectTasks(rows)
	if err != nil {
		return nil, 0, err
	}
	return tasks, total, nil
}

func (r *postgresTaskRepository) GetByID(ctx context.Context, id string) (Domain.Task, error) {
	row := r.pool.QueryRow(ctx, `SELECT `+taskColumns+` FROM tasks WHERE id = $1`, id)
	t, err := scanTask(row)
	if err != nil {
		return Domain.Task{}, mapPgErr(err)
	}
	return t, nil
}

func (r *postgresTaskRepository) Create(ctx context.Context, task Domain.Task) (Domain.Task, error) {
	if task.ReminderStatus == "" {
		task.ReminderStatus = Domain.ReminderPending
	}
	row := r.pool.QueryRow(ctx,
		`INSERT INTO tasks (user_id, title, description, due_date, status, reminder_status)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING `+taskColumns,
		task.OwnerID, task.Title, task.Description, task.DueDate.UTC(), task.Status, task.ReminderStatus)
	t, err := scanTask(row)
	if err != nil {
		return Domain.Task{}, mapPgErr(err)
	}
	return t, nil
}

func (r *postgresTaskRepository) Update(ctx context.Context, id string, task Domain.Task) (Domain.Task, error) {
	row := r.pool.QueryRow(ctx,
		`UPDATE tasks
		    SET title = $2, description = $3, due_date = $4, status = $5,
		        reminder_status = $6, updated_at = now()
		  WHERE id = $1
		  RETURNING `+taskColumns,
		id, task.Title, task.Description, task.DueDate.UTC(), task.Status, task.ReminderStatus)
	t, err := scanTask(row)
	if err != nil {
		return Domain.Task{}, mapPgErr(err)
	}
	return t, nil
}

func (r *postgresTaskRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return mapPgErr(err)
	}
	if tag.RowsAffected() == 0 {
		return Domain.ErrNotFound
	}
	return nil
}

// ClaimDueTasks atomically moves pending tasks to queued. FOR UPDATE SKIP LOCKED
// guarantees two concurrent claimers never receive the same row.
func (r *postgresTaskRepository) ClaimDueTasks(ctx context.Context, dueBefore time.Time, limit int) ([]Domain.Task, error) {
	rows, err := r.pool.Query(ctx,
		`UPDATE tasks
		    SET reminder_status = 'queued', updated_at = now()
		  WHERE id IN (
		        SELECT id FROM tasks
		         WHERE reminder_status = 'pending'
		           AND due_date <= $1
		           AND status <> 'done'
		         ORDER BY due_date
		         LIMIT $2
		         FOR UPDATE SKIP LOCKED)
		  RETURNING `+taskColumns,
		dueBefore.UTC(), limit)
	if err != nil {
		return nil, mapPgErr(err)
	}
	defer rows.Close()
	return collectTasks(rows)
}

func (r *postgresTaskRepository) SetReminderStatus(ctx context.Context, id string, status string) error {
	switch status {
	case Domain.ReminderPending, Domain.ReminderQueued, Domain.ReminderSent, Domain.ReminderFailed:
	default:
		return fmt.Errorf("%w: invalid reminder status %q", Domain.ErrInvalidInput, status)
	}
	tag, err := r.pool.Exec(ctx,
		`UPDATE tasks SET reminder_status = $2, updated_at = now() WHERE id = $1`, id, status)
	if err != nil {
		return mapPgErr(err)
	}
	if tag.RowsAffected() == 0 {
		return Domain.ErrNotFound
	}
	return nil
}
