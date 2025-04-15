package services

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	psgrep "CRUD-HOME-APPLIANCE-STORE/internal/repositories/postgres"
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// type ProductRepository interface {
// 	Create(ctx context.Context, product domain.Product) error
// 	GetAll(ctx context.Context, limit, offset int) ([]domain.Product, error)
// 	GetById(ctx context.Context, id uuid.UUID) (domain.Product, error)
// 	Update(ctx context.Context, id uuid.UUID, decrease int) error
// 	Delete(ctx context.Context, id uuid.UUID) error
// }

type productWriter interface {
	Create(ctx context.Context, product domain.Product) error
	Update(ctx context.Context, id uuid.UUID, decrease int) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type productReader interface {
	GetAll(ctx context.Context, limit, offset int) ([]domain.Product, error)
	GetById(ctx context.Context, id uuid.UUID) (*domain.Product, error)
}

type productService struct {
	logger *logger.Logger
	// repo   ProductRepository
	reader productReader
	writer productWriter
}

func NewProductService( /*repo ProductRepository,*/ reader productReader, writer productWriter, logger *logger.Logger) *productService {
	logger.Debug("Product service is created")
	return &productService{
		logger: logger,
		// repo:   repo,
		reader: reader,
		writer: writer,
	}
}

func (s *productService) Create(ctx context.Context, product *domain.Product) error {
	op := "services.productService.Create"

	if err := s.writer.Create(ctx, *product); err != nil {
		s.logger.Debug("Repo Error", logger.Err(err), "op", op)
		return fmt.Errorf("Why are you gay: %v", err)
	}

	s.logger.Debug("Product created", "op", op)
	return nil
}

func (s *productService) GetAll(ctx context.Context, limit, offset int) ([]domain.Product, error) {
	op := "services.productService.GetAll"

	if limit <= 0 || offset < 0 {
		s.logger.Debug("Error: limit <= 0 or offset < 0", "op", op)
		return nil, fmt.Errorf("Product service: limit or offset invalid payload")
	}

	products, err := s.reader.GetAll(ctx, limit, offset)
	if errors.Is(err, psgrep.ErrClientNotFound) {
		s.logger.Debug("Product's not found", "op", op)
		return nil, err
	}

	if err != nil {
		s.logger.Debug("Repo error", logger.Err(err), "op", op)
		return nil, fmt.Errorf("Product Service: %v", err)
	}

	return products, nil
}

func (s *productService) GetById(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	// op := "services.productService.GetById"

	return nil, nil
}

func (s *productService) Update(ctx context.Context, id string, decrease int) error {
	op := "services.productsService.Update"
	if decrease <= 0 {
		s.logger.Debug("Not valid input value", "value", decrease, "op", op)
		return fmt.Errorf("Product Service: Decrease value must be greater that 0")
	}

	s.logger.Debug("Data is updated", "op", op)
	return nil
}
