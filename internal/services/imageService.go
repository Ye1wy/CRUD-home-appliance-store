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

type ImageReader interface {
	GetAll(ctx context.Context, limit, offset int) ([]domain.Image, error)
	GetById(ctx context.Context, id uuid.UUID) (*domain.Image, error)
}

type imageService struct {
	uow    uow.UOW
	reader ImageReader
	logger *logger.Logger
}

func NewImageService(reader ImageReader, unit uow.UOW, logger *logger.Logger) *imageService {
	return &imageService{
		uow:    unit,
		reader: reader,
		logger: logger,
	}
}

func (s *imageService) Create(ctx context.Context, image domain.Image) error {
	op := "services.imageService.Create"

	err := s.uow.Do(ctx, func(ctx context.Context, tx uow.Transaction) error {
		uowOp := op + ".uow"
		imageRepoGen, err := getReposiotry(tx, uow.ImageRepoName, s.logger)
		if err != nil {
			s.logger.Debug("get image repository generator is unable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: get image repository generator is unable: %v", uowOp, err)
		}

		imageRepo := imageRepoGen.(*postgres.ImageRepository)
		if err := imageRepo.Create(ctx, image); err != nil {
			s.logger.Debug("failed to create image", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: failed to create image: %v", uowOp, err)
		}

		return nil
	})

	if err != nil {
		s.logger.Debug("Somthing wrong with UOW creating", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unit of work problem %v", op, err)
	}

	return nil
}

func (s *imageService) GetAll(ctx context.Context, limit, offset int) ([]domain.Image, error) {
	op := "services.imageService.GetAll"

	if limit <= 0 || offset <= 0 {
		s.logger.Debug("Invalid parameter limit and offset", "op", op)
		return nil, fmt.Errorf("%s: %w", op, crud_errors.ErrInvalidParam)
	}

	images, err := s.reader.GetAll(ctx, limit, offset)
	if err != nil {
		s.logger.Debug("Somthing wrong in repository", logger.Err(err), "op", op)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return images, nil
}

func (s *imageService) GetById(ctx context.Context, id uuid.UUID) (*domain.Image, error) {
	op := "services.imageService.GetById"

	image, err := s.reader.GetById(ctx, id)
	if err != nil {
		s.logger.Debug("Extract data is failed", logger.Err(err), "op", op)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return image, nil
}

func (s *imageService) Update(ctx context.Context, image *domain.Image) error {
	op := "services.imageService.Update"

	err := s.uow.Do(ctx, func(ctx context.Context, tx uow.Transaction) error {
		uowOp := op + ".uow"
		imageRepoGen, err := getReposiotry(tx, uow.ImageRepoName, s.logger)
		if err != nil {
			s.logger.Debug("get image repository generator is unable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: get image repository generator is unable: %v", uowOp, err)
		}

		imageRepo := imageRepoGen.(*postgres.ImageRepository)
		if err := imageRepo.Update(ctx, image); err != nil {
			s.logger.Debug("failed to update image", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: failed to update image: %v", uowOp, err)
		}

		return nil
	})

	if err != nil {
		s.logger.Debug("Somthing wrong with UOW updating", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unit of work update problem: %v", op, err)
	}

	return nil
}

func (s *imageService) Delete(ctx context.Context, id uuid.UUID) error {
	op := "services.imageService.Delete"

	err := s.uow.Do(ctx, func(ctx context.Context, tx uow.Transaction) error {
		uowOp := op + ".uow"
		imageRepoGen, err := getReposiotry(tx, uow.ImageRepoName, s.logger)
		if err != nil {
			s.logger.Debug("Image transaction problem on creating", logger.Err(err), "op", uowOp)
			return err
		}

		imageRepo := imageRepoGen.(*postgres.ImageRepository)
		savepoint := `sp_delete_address`
		err = safeDeleteAddress(ctx, tx.GetTX(), id, imageRepo.Delete, s.logger, uowOp, savepoint)
		if err != nil {
			s.logger.Debug("unable to safe delete address", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: unable to safe delete address: %v", uowOp, err)
		}

		return nil
	})

	if err != nil {
		s.logger.Debug("Somthin wrong with UOW deleting", logger.Err(err), "op", op)
		return fmt.Errorf("Image service: unit of work delete problem: %w", err)
	}

	return nil
}
