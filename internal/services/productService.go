package services

import (
	crud_errors "CRUD-HOME-APPLIANCE-STORE/internal/errors"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/repositories/postgres"
	"CRUD-HOME-APPLIANCE-STORE/internal/uow"
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type productReader interface {
	GetAll(ctx context.Context, limit, offset int) ([]domain.Product, error)
	GetById(ctx context.Context, id uuid.UUID) (*domain.Product, error)
}

type productService struct {
	uow    uow.UOW
	reader productReader
	logger *logger.Logger
}

func NewProductService(reader productReader, uow uow.UOW, logger *logger.Logger) *productService {
	logger.Debug("Product service is created")
	return &productService{
		uow:    uow,
		reader: reader,
		logger: logger,
	}
}

func (s *productService) Create(ctx context.Context, product domain.Product) error {
	op := "services.productService.Create"
	err := s.uow.Do(ctx, func(ctx context.Context, tx uow.Transaction) error {
		uowOp := op + ".uow"
		productRepoGen, err := getReposiotry(tx, uow.ProductRepoName, s.logger)
		if err != nil {
			s.logger.Debug("get product repository generator is unable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: get product repository generator is unable: %v", uowOp, err)
		}

		productRepo := productRepoGen.(*postgres.ProductRepo)

		if err := productRepo.Create(ctx, product); err != nil {
			s.logger.Debug("failed to create product", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: failed to create product: %v", uowOp, err)
		}

		return nil
	})

	if err != nil {
		s.logger.Debug("Something wrong with UOW creating", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unit of work creating problem: %v", op, err)
	}

	return nil
}

func (s *productService) GetAll(ctx context.Context, limit, offset int) ([]domain.Product, error) {
	op := "services.productService.GetAll"

	if limit <= 0 || offset < 0 {
		s.logger.Debug("Error: limit <= 0 or offset < 0", "op", op)
		return nil, fmt.Errorf("%s: %w", op, crud_errors.ErrInvalidParam)
	}

	products, err := s.reader.GetAll(ctx, limit, offset)
	if err != nil {
		s.logger.Debug("Repo error", logger.Err(err), "op", op)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return products, nil
}

func (s *productService) GetById(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	op := "services.productService.GetById"
	product, err := s.reader.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return product, nil
}

func (s *productService) Update(ctx context.Context, id uuid.UUID, decrease int) error {
	op := "services.productsService.Update"
	if decrease <= 0 {
		s.logger.Debug("Not valid input value", "value", decrease, "op", op)
		return fmt.Errorf("%s: %w", op, crud_errors.ErrInvalidParam)
	}

	err := s.uow.Do(ctx, func(ctx context.Context, tx uow.Transaction) error {
		uowOp := op + ".uow"
		productRepoGen, err := getReposiotry(tx, uow.ProductRepoName, s.logger)
		if err != nil {
			s.logger.Debug("get product repository generator is unable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: get product repository generator is unable: %v", uowOp, err)
		}

		productRepo := productRepoGen.(*postgres.ProductRepo)
		if err := productRepo.Update(ctx, id, decrease); err != nil {
			s.logger.Debug("failed to update stock with product", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: failed to update stock with product: %v", uowOp, err)
		}

		return nil
	})

	if err != nil {
		s.logger.Debug("Somthing wrong with UOW updating", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unit of work update problem: %v", op, err)
	}

	return nil
}

func (s *productService) Delete(ctx context.Context, id uuid.UUID) error {
	op := "services.productService.Delete"

	err := s.uow.Do(ctx, func(ctx context.Context, tx uow.Transaction) error {
		uowOp := op + ".uow"
		productRepoGen, err := getReposiotry(tx, uow.ProductRepoName, s.logger)
		if err != nil {
			s.logger.Debug("get product repository generator is unable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: get product repository generator is unable: %v", uowOp, err)
		}

		productRepo := productRepoGen.(*postgres.ProductRepo)

		savepoint := `sp_delete_address`
		err = safeDeleteAddress(ctx, tx.GetTX(), id, productRepo.Delete, s.logger, uowOp, savepoint)
		if err != nil {
			s.logger.Debug("unable to safe delete product", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: unable to safe delete product: %v", uowOp, err)
		}

		return nil
	})

	if err != nil {
		s.logger.Debug("Somthin wrong with UOW deleting", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unit of work delete problem: %v", op, err)
	}

	return nil
}
