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

var imageRepoName = uow.RepositoryName("image")

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
		repo, err := tx.Get(clientRepoName)
		if err != nil {
			s.logger.Debug("Image transaction problem on creating", logger.Err(err), "op", op)
			return err
		}

		repoGen := repo.(uow.RepositoryGenerator)(tx.GetTX(), s.logger)
		imageRepo := repoGen.(*postgres.ImageRepository)
		return imageRepo.Create(ctx, image)
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
		return nil, fmt.Errorf("%s: %v", op, err)
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
		repo, err := tx.Get(imageRepoName)
		if err != nil {
			s.logger.Debug("Product transaction problem on updating", logger.Err(err), "op", op)
			return err
		}

		repoGen := repo.(uow.RepositoryGenerator)(tx.GetTX(), s.logger)
		imageRepo := repoGen.(postgres.ImageRepository)
		return imageRepo.Update(ctx, image)
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
		repo, err := tx.Get(imageRepoName)
		if err != nil {
			s.logger.Debug("Get transaction problem", logger.Err(err), "op", op)
			return err
		}

		repoGen := repo.(uow.RepositoryGenerator)(tx.GetTX(), s.logger)
		imageRepo := repoGen.(postgres.ImageRepository)
		return imageRepo.Delete(ctx, id)
	})

	if err != nil {
		s.logger.Debug("Somthin wrong with UOW deleting", logger.Err(err), "op", op)
		return fmt.Errorf("Image service: unit of work delete problem: %w", err)
	}

	return nil
}
