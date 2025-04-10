package services

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"context"

	"github.com/google/uuid"
)

type ProductRepository interface {
	Create(ctx context.Context, product domain.Product) error
	GetAll(ctx context.Context, limit, offset int) ([]domain.Product, error)
	GetById(ctx context.Context, id uuid.UUID) (domain.Product, error)
	Update(ctx context.Context, id uuid.UUID, decrease int) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// type productService struct {
// 	repo ProductRepository
// }

// func NewProductService(repo ProductRepositoryInterface, logger *logger.Logger) *productService {
// 	service := NewCrudService(repo, mapper.ProductToDTO, mapper.ProductToModel, logger)
// 	logger.Debug("Product service is created")
// 	return &productService{
// 		CrudService: service,
// 		repo:        repo,
// 	}
// }

// func (ps *productService) DecreaseStock(ctx context.Context, id string, decrease int) error {
// 	op := "services.productsService.DecreaseStock"
// 	if decrease <= 0 {
// 		ps.Logger.Debug("Not valid input value", "value", decrease, "op", op)
// 		return fmt.Errorf("Product Service: Decrease value must be greater that 0")
// 	}

// 	ps.Logger.Debug("Data is updated", "op", op)
// 	return ps.repo.DecreaseParameter(ctx, id, decrease)
// }
