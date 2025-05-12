package repository

import (
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

var (
	ErrRepoIsExist     = errors.New("repository is exist")
	ErrRepoIsNotExitst = errors.New("repository is not exitst")
)

type RepositoryName string
type Repository any
type RepositoryGenerator func(tx *pgx.Tx, log *logger.Logger) Repository

type Transaction interface {
	Get(name RepositoryName) (Repository, error)
}

type transaction struct {
	tx    pgx.Tx
	repos map[RepositoryName]RepositoryGenerator
}

func NewTransaction(tx pgx.Tx, repos map[RepositoryName]RepositoryGenerator) *transaction {
	return &transaction{
		tx:    tx,
		repos: repos,
	}
}

func (tx *transaction) Get(name RepositoryName) (Repository, error) {
	if repo, ok := tx.repos[name]; ok {
		return repo, nil
	}

	return nil, ErrRepoIsNotExitst
}

type unitOfWork struct {
	db           pgx.Tx
	logger       *logger.Logger
	repositories map[RepositoryName]RepositoryGenerator
}

func NewUnitOfWork(tx pgx.Tx, logger *logger.Logger) *unitOfWork {
	return &unitOfWork{
		db:           tx,
		logger:       logger,
		repositories: make(map[RepositoryName]RepositoryGenerator),
	}
}

func (uow *unitOfWork) Register(name RepositoryName, gen RepositoryGenerator) error {
	if _, ok := uow.repositories[name]; ok {
		return ErrRepoIsExist
	}

	uow.repositories[name] = gen

	return nil
}

func (uow *unitOfWork) Remove(name RepositoryName) error {
	if _, ok := uow.repositories[name]; !ok {
		return ErrRepoIsNotExitst
	}

	delete(uow.repositories, name)
	return nil
}

func (uow *unitOfWork) Clear() {
	uow.repositories = make(map[RepositoryName]RepositoryGenerator)
}

func (uow *unitOfWork) Do(ctx context.Context, fn func(ctx context.Context, tx Transaction) error) error {
	tx, err := uow.db.Begin(ctx)
	if err != nil {
		return err
	}

	if err := fn(ctx, NewTransaction(tx, uow.repositories)); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return err
		}

		return err
	}

	return tx.Commit(ctx)
}
