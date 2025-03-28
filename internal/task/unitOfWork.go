package task

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type UnitOfWork interface {
	Begin(ctx context.Context) (Repository, error)
	Save(ctx context.Context, repository Repository) error
	Rollback(ctx context.Context, repository Repository) error
}

type postgresUnitOfWork struct {
	pool *pgxpool.Pool
}

func NewPostgresUnitOfWork(pool *pgxpool.Pool) (*postgresUnitOfWork, error) {
	return &postgresUnitOfWork{pool}, nil
}

func (uow *postgresUnitOfWork) Begin(ctx context.Context) (*PosgresRepository, error) {
	transaction, err := uow.pool.Begin(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin PostgreSQL transaction")
	}

	return newPosgresRepository(transaction), nil
}

func (_ *postgresUnitOfWork) Save(ctx context.Context, repository *PosgresRepository) error {
	return repository.transaction.Commit(ctx)
}

func (uow *postgresUnitOfWork) Rollback(ctx context.Context, repository *PosgresRepository) error {
	return repository.transaction.Rollback(ctx)
}
