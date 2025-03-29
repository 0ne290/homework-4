package task

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UnitOfWork interface {
	Begin(ctx context.Context) (Repository, error)
	Save(ctx context.Context, repository Repository) error
	Rollback(ctx context.Context, repository Repository) error
}

type PostgresUnitOfWork struct {
	pool *pgxpool.Pool
}

func NewPostgresUnitOfWork(pool *pgxpool.Pool) *PostgresUnitOfWork {
	return &PostgresUnitOfWork{pool}
}

func (uow *PostgresUnitOfWork) Begin(ctx context.Context) (Repository, error) {
	transaction, err := uow.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return newPosgresRepository(transaction), nil
}

func (*PostgresUnitOfWork) Save(ctx context.Context, repository Repository) error {
	return repository.save(ctx)
}

func (*PostgresUnitOfWork) Rollback(ctx context.Context, repository Repository) error {
	return repository.rollback(ctx)
}
