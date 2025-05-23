package repository

import (
	crud_errors "CRUD-HOME-APPLIANCE-STORE/internal/errors"
	"CRUD-HOME-APPLIANCE-STORE/internal/uow"
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"context"

	"github.com/jackc/pgx/v5"
)

type transaction struct {
	tx    pgx.Tx
	repos map[uow.RepositoryName]uow.RepositoryGenerator
}

func NewTransaction(tx pgx.Tx, repos map[uow.RepositoryName]uow.RepositoryGenerator) *transaction {
	return &transaction{
		tx:    tx,
		repos: repos,
	}
}

func (tx *transaction) Get(name uow.RepositoryName) (uow.Repository, error) {
	if repo, ok := tx.repos[name]; ok {
		return repo, nil
	}

	return nil, crud_errors.ErrRepoIsNotExitst
}

func (tx *transaction) GetTX() pgx.Tx {
	return tx.tx
}

type unitOfWork struct {
	db           *pgx.Conn
	logger       *logger.Logger
	repositories map[uow.RepositoryName]uow.RepositoryGenerator
}

func NewUnitOfWork(conn *pgx.Conn, logger *logger.Logger) *unitOfWork {
	return &unitOfWork{
		db:           conn,
		logger:       logger,
		repositories: make(map[uow.RepositoryName]uow.RepositoryGenerator),
	}
}

func (unit *unitOfWork) Register(name uow.RepositoryName, gen uow.RepositoryGenerator) error {
	if _, ok := unit.repositories[name]; ok {
		return crud_errors.ErrRepoIsExist
	}

	unit.repositories[name] = gen

	return nil
}

func (unit *unitOfWork) Remove(name uow.RepositoryName) error {
	if _, ok := unit.repositories[name]; !ok {
		return crud_errors.ErrRepoIsNotExitst
	}

	delete(unit.repositories, name)
	return nil
}

func (unit *unitOfWork) Clear() {
	unit.repositories = make(map[uow.RepositoryName]uow.RepositoryGenerator)
}

func (unit *unitOfWork) Do(ctx context.Context, fn func(ctx context.Context, tx uow.Transaction) error) error {
	tx, err := unit.db.Begin(ctx)
	if err != nil {
		return err
	}

	if err := fn(ctx, NewTransaction(tx, unit.repositories)); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return err
		}

		return err
	}

	return tx.Commit(ctx)
}
