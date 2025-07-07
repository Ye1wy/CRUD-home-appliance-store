package services

import (
	crud_errors "CRUD-HOME-APPLIANCE-STORE/internal/errors"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/uow"
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type productReader interface {
	GetAll(ctx context.Context, limit, offset int) ([]domain.Product, error)
	GetById(ctx context.Context, id uuid.UUID) (*domain.Product, error)
}

type productWriter interface {
	Create(ctx context.Context, product *domain.Product) error
	Update(ctx context.Context, id uuid.UUID, decrease int) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type productService struct {
	uow    uow.UOW
	reader productReader
	logger *logger.Logger
}

func NewProductService(reader productReader, uow uow.UOW, logger *logger.Logger) *productService {
	logger.Debug("product service is created")
	return &productService{
		uow:    uow,
		reader: reader,
		logger: logger,
	}
}

func (s *productService) Create(ctx context.Context, product *domain.Product) error {
	op := "services.productService.Create"

	err := s.uow.Do(ctx, func(ctx context.Context, tx uow.Transaction) error {
		uowOp := op + ".uow"
		// supplierRepoGen, err := getReposiotry(tx, uow.SupplierRepoName, s.logger)
		// if err != nil {
		// 	s.logger.Error("get supplier repository generator is unable", logger.Err(err), "op", uowOp)
		// 	return fmt.Errorf("%s: get supplier repository generator is unable: %v", uowOp, err)
		// }

		// supplierRepo, ok := supplierRepoGen.(*postgres.SupplierRepo)
		// if !ok {
		// 	s.logger.Error("Conversion problem, not contained expected convesion", "op", op)
		// 	return fmt.Errorf("%s: %w", uowOp, crud_errors.ErrConversionProblem)
		// }

		// // Check supplier is exist or not. If not exist give error invalid payload on controller
		// if _, err := supplierRepo.GetById(ctx, product.Supplier.Id); err != nil {
		// 	if errors.Is(err, crud_errors.ErrNotFound) {
		// 		s.logger.Warn("supplier not found", "op", uowOp)
		// 		return fmt.Errorf("%s: %w", uowOp, err)
		// 	}

		// 	s.logger.Error("get supplier by id is failed", logger.Err(err), "op", uowOp)
		// 	return fmt.Errorf("%s: supplier get is failed or supplier is not found: %w", uowOp, err)
		// }

		// imageRepoGen, err := getReposiotry(tx, uow.ImageRepoName, s.logger)
		// if err != nil {
		// 	s.logger.Error("get image repository generator is unable", logger.Err(err), "op", uowOp)
		// 	return fmt.Errorf("%s: get image repository generator is unable: %v", uowOp, err)
		// }

		// imageRepo, ok := imageRepoGen.(*postgres.ImageRepo)
		// if !ok {
		// 	s.logger.Error("Conversion problem, not contained expected convesion", "op", op)
		// 	return fmt.Errorf("%s: %w", uowOp, crud_errors.ErrConversionProblem)
		// }

		// // Check image is exist or not. If not exist give error invalid payload on controller
		// if _, err := imageRepo.GetById(ctx, product.Image.Id); err != nil {
		// 	if errors.Is(err, crud_errors.ErrNotFound) {
		// 		s.logger.Warn("supplier not found", "op", uowOp)
		// 		return fmt.Errorf("%s: iamge id: %w", uowOp, err)
		// 	}

		// 	s.logger.Error("unable to create image with creating product", logger.Err(err), "op", uowOp)
		// 	return fmt.Errorf("%s: failed get image by id %w", uowOp, err)
		// }

		// Creating a product repository, type assertion is necessary, since the function returns any
		productRepoGen, err := getReposiotry(tx, uow.ProductRepoName, s.logger)
		if err != nil {
			s.logger.Error("get product repository generator is unable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: get product repository generator is unable: %v", uowOp, err)
		}

		// Type assertion for extract product reposiotry
		productRepo, ok := productRepoGen.(productWriter)
		if !ok {
			s.logger.Error("Conversion problem, not contained expected convesion", "op", op)
			return fmt.Errorf("%s: %w", uowOp, crud_errors.ErrConversionProblem)
		}

		if err := productRepo.Create(ctx, product); err != nil {
			s.logger.Error("failed to create product", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: failed to create product: %v", uowOp, err)
		}

		return nil
	})

	if err != nil {
		s.logger.Error("something wrong with UOW creating", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unit of work creating problem: %w", op, err)
	}

	return nil
}

func (s *productService) GetAll(ctx context.Context, limit, offset int) ([]domain.Product, error) {
	op := "services.productService.GetAll"

	if limit <= 0 || offset < 0 {
		s.logger.Error("limit cannot be 0 or less and offset cannot be less by 0", "op", op)
		return nil, fmt.Errorf("%s: %w", op, crud_errors.ErrInvalidParam)
	}

	products, err := s.reader.GetAll(ctx, limit, offset)
	if err != nil {
		if errors.Is(err, crud_errors.ErrNotFound) {
			s.logger.Warn("No content", "op", op)
		} else {
			s.logger.Error("error detected", logger.Err(err), "op", op)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return products, nil
}

func (s *productService) GetById(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	op := "services.productService.GetById"
	product, err := s.reader.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, crud_errors.ErrNotFound) {
			s.logger.Debug("product not found", "op", op)
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		s.logger.Error("failed get product data by id", logger.Err(err), "op", op)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return product, nil
}

func (s *productService) Update(ctx context.Context, id uuid.UUID, decrease int) error {
	op := "services.productsService.Update"
	if decrease <= 0 {
		s.logger.Error("not valid input value", "value", decrease, "op", op)
		return fmt.Errorf("%s: %w", op, crud_errors.ErrInvalidParam)
	}

	err := s.uow.Do(ctx, func(ctx context.Context, tx uow.Transaction) error {
		uowOp := op + ".uow"
		productRepoGen, err := getReposiotry(tx, uow.ProductRepoName, s.logger)
		if err != nil {
			s.logger.Error("get product repository generator is unable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: get product repository generator is unable: %v", uowOp, err)
		}

		productRepo, ok := productRepoGen.(productWriter)
		if !ok {
			s.logger.Error("Conversion problem, not contained expected convesion", "op", op)
			return fmt.Errorf("%s: %w", uowOp, crud_errors.ErrConversionProblem)
		}

		if err := productRepo.Update(ctx, id, decrease); err != nil {
			if errors.Is(err, crud_errors.ErrNotFound) {
				s.logger.Debug("product not found", "op", uowOp)
				return fmt.Errorf("%s, %w", uowOp, err)
			}

			s.logger.Error("failed to update stock with product", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: failed to update stock with product: %w", uowOp, err)
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, crud_errors.ErrNotFound) {
			s.logger.Warn("update initialize is unable: product not found", "op", op)
			return fmt.Errorf("%s: %w", op, err)
		}

		s.logger.Error("something wrong with UOW updating", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unit of work update problem: %w", op, err)
	}

	return nil
}

func (s *productService) Delete(ctx context.Context, id uuid.UUID) error {
	op := "services.productService.Delete"

	err := s.uow.Do(ctx, func(ctx context.Context, tx uow.Transaction) error {
		uowOp := op + ".uow"
		productRepoGen, err := getReposiotry(tx, uow.ProductRepoName, s.logger)
		if err != nil {
			s.logger.Error("get product repository generator is unable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: get product repository generator is unable: %v", uowOp, err)
		}

		productRepo, ok := productRepoGen.(productWriter)
		if !ok {
			s.logger.Error("Conversion problem, not contained expected convesion", "op", op)
			return fmt.Errorf("%s: %w", uowOp, crud_errors.ErrConversionProblem)
		}

		savepoint := `sp_delete_product`
		err = safeDelete(ctx, tx.GetTX(), id, productRepo.Delete, s.logger, uowOp, savepoint)
		if err != nil {
			if errors.Is(err, crud_errors.ErrNotFound) {
				s.logger.Debug("product not found", "op", op)
				return nil
			}

			s.logger.Error("unable to safe delete product", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: unable to safe delete product: %v", uowOp, err)
		}

		return nil
	})

	if err != nil {
		s.logger.Error("something wrong with UOW deleting", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unit of work delete problem: %v", op, err)
	}

	return nil
}
