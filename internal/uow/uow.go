package uow

import (
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"context"

	"github.com/jackc/pgx/v5"
)

type RepositoryName string
type Repository any
type RepositoryGenerator func(tx pgx.Tx, log *logger.Logger) Repository

type Transaction interface {
	Get(name RepositoryName) (Repository, error)
}

type UOW interface {
	Register(name RepositoryName, gen RepositoryGenerator) error
	Remove(name RepositoryName) error
	Clear()
	Do(ctx context.Context, fn func(ctx context.Context, tx Transaction) error) error
}
