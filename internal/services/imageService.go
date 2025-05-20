package services

import (
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
		return fmt.Errorf("Image service: unit of work problem %w", err)
	}

	s.logger.Debug("Image is created", "op", op)
	return nil
}

func (s *imageService) GetAll(ctx context.Context, limit, offset int) ([]domain.Image, error) {
	op := "services.imageService.GetAll"

	if limit <= 0 || offset <= 0 {
		s.logger.Debug("Invalid parameter limit and offset", "op", op)
		return nil, ErrInvalidParam
	}

	images, err := s.reader.GetAll(ctx, limit, offset)
	if err != nil {
		s.logger.Debug("Somthing wrong in repository", logger.Err(err), "op", op)
		return nil, err
	}

	return images, nil
}

func (s *imageService) GetById(ctx context.Context, id uuid.UUID) (*domain.Image, error) {
	op := "services.imageService.GetById"

	return nil, nil
}

func (s *imageService) Update(ctx context.Context, image *domain.Image) error {
	return nil
}

func (s *imageService) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
