package services

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/logger"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"context"

	"github.com/google/uuid"
)

type SupplierRepository interface {
	Create(ctx context.Context, supplier *domain.Supplier) error
	GetAll(ctx context.Context, limit, offset int) ([]domain.Supplier, error)
	GetById(ctx context.Context, id uuid.UUID) (domain.Supplier, error)
	Update(ctx context.Context, supplier domain.Supplier) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type supplierService struct {
	Repo SupplierRepository
}

func NewSupplierService(repo SupplierRepository, logger *logger.Logger) *supplierService {
	// service := NewCrudService(repo, mapper.SupplierToDTO, mapper.SupplierToModel, logger)
	logger.Debug("Supplier serivce created")
	return nil
}
