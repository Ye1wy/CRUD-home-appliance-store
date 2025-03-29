package services

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/dto"
	"CRUD-HOME-APPLIANCE-STORE/internal/logger"
	"CRUD-HOME-APPLIANCE-STORE/internal/mapper"
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
	"CRUD-HOME-APPLIANCE-STORE/internal/repositories"
	"context"
	"errors"
)

type ProductServiceInterface interface {
	CrudServiceInterface[model.Product, dto.ProductDTO]
	DecreaseStock(ctx context.Context, id string, decrease int) error
}

type productServiceImpl struct {
	*CrudService[model.Product, dto.ProductDTO]
	repo repositories.ProductRepositoryInterface
}

func NewProductService(rep repositories.ProductRepositoryInterface, logger *logger.Logger) *productServiceImpl {
	service := NewCrudService(rep, mapper.ProductToDTO, mapper.ProductToModel, logger)
	logger.Debug("Product service is created")
	return &productServiceImpl{
		CrudService: service,
		repo:        rep,
	}
}

func (ps *productServiceImpl) DecreaseStock(ctx context.Context, id string, decrease int) error {
	if decrease <= 0 {
		return errors.New("decrease value must be greater that 0")
	}

	return ps.repo.DecreaseParameter(ctx, id, decrease)
}
