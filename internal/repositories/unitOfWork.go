package repository

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/repositories/postgres"
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ClientWriter interface {
	Create(ctx context.Context, client domain.Client) error
	UpdateAddress(ctx context.Context, id, address uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type Repo interface {
	ClientRepo() postgres.ClientWriter
}

type unitOfWork struct {
	db     pgx.Tx
	logger *logger.Logger
	client ClientWriter
	// product productRepo
	// supplie supplierRepo
	// image imageRepo
}

func NewUnitOfWork(tx pgx.Tx, logger *logger.Logger) *unitOfWork {
	client := postgres.NewClientRepository(tx, logger)
	return &unitOfWork{
		db:     tx,
		logger: logger,
		client: client,
	}
}

func (uow *unitOfWork) Transaction(ctx context.Context, fn func(repo Repo) error) error {
	tx, err := uow.db.Begin(ctx)
	if err != nil {
		return err
	}

	repos := postgres.NewRepository(tx, uow.logger)

	if err := fn(repos); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
