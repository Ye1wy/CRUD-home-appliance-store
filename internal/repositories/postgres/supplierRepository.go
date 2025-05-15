package postgres

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type SupplierRepo struct {
	*basePostgresRepository
}

func NewSupplierRepository(tx pgx.Tx, logger *logger.Logger) *SupplierRepo {
	repo := newBasePostgresRepository(tx, logger)
	logger.Debug("Postgres Supplier Repository is created")
	return &SupplierRepo{
		repo,
	}
}

func (r *SupplierRepo) Create(ctx context.Context, supplier domain.Supplier) error {
	return nil
}

func (r *SupplierRepo) GetAll(ctx context.Context, limit, offset int) ([]domain.Supplier, error) {
	return nil, nil
}

func (r *SupplierRepo) GetById(ctx context.Context, id uuid.UUID) (*domain.Supplier, error) {
	return nil, nil
}

func (r *SupplierRepo) Update(ctx context.Context, id, address uuid.UUID) error {
	return nil
}

func (r *SupplierRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
