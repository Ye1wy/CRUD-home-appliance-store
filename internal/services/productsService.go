package services

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/dto"
	"CRUD-HOME-APPLIANCE-STORE/internal/logger"
	"CRUD-HOME-APPLIANCE-STORE/internal/mapper"
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
	"CRUD-HOME-APPLIANCE-STORE/internal/repositories"
	"context"
	"fmt"
)

type ProductServiceInterface interface {
	CrudServiceInterface[model.Product, dto.ProductDTO]
	DecreaseStock(ctx context.Context, id string, decrease int) error
}

type productServiceImpl struct {
	*CrudService[model.Product, dto.ProductDTO]
	repo repositories.ProductRepositoryInterface
}

func NewProductService(repo repositories.ProductRepositoryInterface, logger *logger.Logger) *productServiceImpl {
	service := NewCrudService(repo, mapper.ProductToDTO, mapper.ProductToModel, logger)
	logger.Debug("Product service is created")
	return &productServiceImpl{
		CrudService: service,
		repo:        repo,
	}
}

func (ps *productServiceImpl) DecreaseStock(ctx context.Context, id string, decrease int) error {
	op := "services.productsService.DecreaseStock"
	if decrease <= 0 {
		ps.Logger.Debug("Not valid input value", "value", decrease, "op", op)
		return fmt.Errorf("Product Service: Decrease value must be greater that 0")
	}

	ps.Logger.Debug("Data is updated", "op", op)
	return ps.repo.DecreaseParameter(ctx, id, decrease)
}
