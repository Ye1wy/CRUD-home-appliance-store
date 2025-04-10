package psgrep

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/logger"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type supplierRepo struct {
	*basePostgresRepository
}

func NewSupplierRepository(conn *pgx.Conn, logger *logger.Logger) *supplierRepo {
	repo := newBasePostgresRepository(conn, logger)
	logger.Debug("Postgres Supplier Repository is created")
	return &supplierRepo{
		repo,
	}
}

func (r *supplierRepo) Create(ctx context.Context, supplier domain.Supplier) error {
	return nil
}

func (r *supplierRepo) GetAll(ctx context.Context, limit, offset int) ([]domain.Supplier, error) {
	return nil, nil
}

func (r *supplierRepo) GetById(ctx context.Context, id uuid.UUID) (*domain.Supplier, error) {
	return nil, nil
}

func (r *supplierRepo) Update(ctx context.Context, supplier domain.Supplier) error {
	return nil
}

func (r *supplierRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
