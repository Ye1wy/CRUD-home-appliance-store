package services

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/repositories/postgres"
	"CRUD-HOME-APPLIANCE-STORE/internal/uow"
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var productRepoName = "product"

type productReader interface {
	GetAll(ctx context.Context, limit, offset int) ([]domain.Product, error)
	GetById(ctx context.Context, id uuid.UUID) (*domain.Product, error)
}

type productService struct {
	uow    uow.UOW
	reader productReader
	logger *logger.Logger
}

func NewProductService(reader productReader, logger *logger.Logger) *productService {
	logger.Debug("Product service is created")
	return &productService{
		reader: reader,
		logger: logger,
	}
}

func (s *productService) Create(ctx context.Context, product domain.Product) error {
	op := "services.productService.Create"
	err := s.uow.Do(ctx, func(ctx context.Context, tx uow.Transaction) error {
		repo, err := tx.Get(uow.RepositoryName(productRepoName))
		if err != nil {
			s.logger.Debug("Create product transaction problem on creating", logger.Err(err), "op", op)
			return err
		}

		productRepo := repo.(postgres.ProductRepo)

		return productRepo.Create(ctx, product)
	})

	if err != nil {
		s.logger.Debug("Something wrong with UOW creating", logger.Err(err), "op", op)
		return fmt.Errorf("Product service: unit of work creating problem: %v", err)
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
	if errors.Is(err, postgres.ErrProductNotFound) {
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
	op := "services.productService.GetById"
	product, err := s.reader.GetById(ctx, id)
	if errors.Is(err, postgres.ErrProductNotFound) {
		s.logger.Debug("Product not found", "op", op)
		return nil, postgres.ErrProductNotFound
	}

	if err != nil {
		s.logger.Debug("Exctract data is failed", logger.Err(err), "op", op)
		return nil, err
	}

	s.logger.Debug("Product retrived", "id", id, "op", op)
	return product, nil
}

func (s *productService) Update(ctx context.Context, id uuid.UUID, decrease int) error {
	op := "services.productsService.Update"
	if decrease <= 0 {
		s.logger.Debug("Not valid input value", "value", decrease, "op", op)
		return fmt.Errorf("Product Service: Decrease value must be greater that 0")
	}

	err := s.uow.Do(ctx, func(ctx context.Context, tx uow.Transaction) error {
		repo, err := tx.Get(uow.RepositoryName(productRepoName))
		if err != nil {
			s.logger.Debug("Product transaction problem on updating", logger.Err(err), "op", op)
			return err
		}

		productRepo := repo.(postgres.ProductRepo)

		return productRepo.Update(ctx, id, decrease)
	})

	if err != nil {
		s.logger.Debug("Somthing wrong with UOW updating", logger.Err(err), "op", op)
		return fmt.Errorf("Product service: unit of work update problem: %v", err)
	}

	s.logger.Debug("Data is updated", "op", op)
	return nil
}

func (s *productService) Delete(ctx context.Context, id uuid.UUID) error {
	op := "serice.productService.Delete"

	err := s.uow.Do(ctx, func(ctx context.Context, tx uow.Transaction) error {
		repo, err := tx.Get(uow.RepositoryName(productRepoName))
		if err != nil {
			s.logger.Debug("Get transaction problem", logger.Err(err), "op", op)
			return err
		}

		userRepo := repo.(postgres.ClientRepo)
		return userRepo.Delete(ctx, id)
	})

	if err != nil {
		s.logger.Debug("Somthin wrong with UOW deleting", logger.Err(err), "op", op)
		return fmt.Errorf("Product service: unit of work delete problem: %v", err)
	}

	s.logger.Debug("Product successfully deleted", "op", op)
	return nil
}
