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
