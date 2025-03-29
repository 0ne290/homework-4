package task

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"homework-4/internal"
)

type Repository interface {
	Create(ctx context.Context, task *Task) error
	GetAll(ctx context.Context) ([]*Task, error)
	GetByUuid(ctx context.Context, taskUuid uuid.UUID) (*Task, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, taskUuid uuid.UUID) error

	save(ctx context.Context) error
	rollback(ctx context.Context) error
}

type PosgresRepository struct {
	transaction pgx.Tx
}

func newPosgresRepository(transaction pgx.Tx) *PosgresRepository {
	return &PosgresRepository{transaction}
}

func (r *PosgresRepository) save(ctx context.Context) error {
	return r.transaction.Commit(ctx)
}

func (r *PosgresRepository) rollback(ctx context.Context) error {
	return r.transaction.Rollback(ctx)
}

func (r *PosgresRepository) Create(ctx context.Context, task *Task) error {
	const query string = "INSERT INTO tasks VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := r.transaction.Exec(ctx, query, task.Uuid, task.Title, task.Description, task.Status, task.CreatedAt, task.UpdatedAt)

	return err
}

func (r *PosgresRepository) GetAll(ctx context.Context) ([]*Task, error) {
	const query string = "SELECT * FROM tasks"

	rows, err := r.transaction.Query(ctx, query)
	if err != nil {
		return make([]*Task, 0), err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task := &Task{}

		err = rows.Scan(&task.Uuid, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return make([]*Task, 0), err
		}

		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return make([]*Task, 0), err
	}

	return tasks, nil
}

func (r *PosgresRepository) GetByUuid(ctx context.Context, taskUuid uuid.UUID) (*Task, error) {
	const query string = "SELECT * FROM tasks WHERE uuid = $1 FOR UPDATE"

	task := &Task{}

	err := r.transaction.QueryRow(ctx, query, taskUuid).Scan(&task.Uuid, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &internal.InvariantViolationError{Message: "task does not exists"}
		}

		return nil, err
	}

	return task, nil
}

func (r *PosgresRepository) Update(ctx context.Context, task *Task) error {
	const query string = "UPDATE tasks SET status = $2, updated_at = $3 WHERE uuid = $1"

	_, err := r.transaction.Exec(ctx, query, task.Uuid, task.Status, task.UpdatedAt)

	return err
}

func (r *PosgresRepository) Delete(ctx context.Context, taskUuid uuid.UUID) error {
	const query string = "DELETE FROM tasks WHERE uuid = $1"

	res, err := r.transaction.Exec(ctx, query, taskUuid)
	if err != nil {
		return err
	}
	count := res.RowsAffected()
	if count == 0 {
		return &internal.InvariantViolationError{Message: "task does not exists"}
	}

	return nil
}
