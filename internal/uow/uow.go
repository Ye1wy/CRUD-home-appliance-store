package uow

import (
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"context"

	"github.com/jackc/pgx/v5"
)

type RepositoryName string
type Repository any
type RepositoryGenerator func(tx pgx.Tx, log *logger.Logger) Repository

const (
	AddressRepoName  = RepositoryName("address")
	ClientRepoName   = RepositoryName("client")
	SupplierRepoName = RepositoryName("supplier")
	ProductRepoName  = RepositoryName("product")
	ImageRepoName    = RepositoryName("image")
)

type CommandTag interface {
	RowsAffected() int64
}

type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close()
	Err() error
}

type Row interface {
	Scan(dest ...any) error
}

type Tx interface {
	Begin(ctx context.Context) (Tx, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	// Exec(ctx context.Context, sql string, arguments ...any) (CommandTag, error)
	// Query(ctx context.Context, sql string, args ...any) (Rows, error)
	// QueryRow(ctx context.Context, sql string, args ...any) Row
}

type Transaction interface {
	Get(name RepositoryName) (Repository, error)
	GetTX() pgx.Tx
}

type UOW interface {
	Register(name RepositoryName, gen RepositoryGenerator) error
	Remove(name RepositoryName) error
	Clear()
	Do(ctx context.Context, fn func(ctx context.Context, tx Transaction) error) error
}
