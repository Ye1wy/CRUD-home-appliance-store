package services

import (
	crud_errors "CRUD-HOME-APPLIANCE-STORE/internal/errors"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/repositories/postgres"
	"CRUD-HOME-APPLIANCE-STORE/internal/uow"
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/google/uuid"
)

type ImageReader interface {
	GetAll(ctx context.Context, limit, offset int) ([]domain.Image, error)
	GetById(ctx context.Context, id uuid.UUID) (*domain.Image, error)
}

func validateImage(data []byte) error {
	_, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("%v :image is corruption %w", err, crud_errors.ErrImageCorruption)
	}

	return nil
}

func generateImageHash(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

type imageService struct {
	uow    uow.UOW
	reader ImageReader
	logger *logger.Logger
}

func NewImageService(reader ImageReader, unit uow.UOW, logger *logger.Logger) *imageService {
	logger.Debug("image service is created")
	return &imageService{
		uow:    unit,
		reader: reader,
		logger: logger,
	}
}

func (s *imageService) Create(ctx context.Context, image *domain.Image) error {
	op := "services.imageService.Create"

	if err := validateImage(image.Data); err != nil {
		s.logger.Debug("image is corruption or taked data is not image", logger.Err(err), "op", op)
		return fmt.Errorf("%s: validation error: %w", op, err)
	}

	err := s.uow.Do(ctx, func(ctx context.Context, tx uow.Transaction) error {
		uowOp := op + ".uow"
		imageRepoGen, err := getReposiotry(tx, uow.ImageRepoName, s.logger)
		if err != nil {
			s.logger.Error("get image repository generator is unable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: get image repository generator is unable: %v", uowOp, err)
		}

		imageRepo := imageRepoGen.(*postgres.ImageRepo)

		image.Hash = generateImageHash(image.Data)

		if err := imageRepo.Create(ctx, image); err != nil {
			s.logger.Error("failed to create image", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: failed to create image: %v", uowOp, err)
		}

		return nil
	})

	if err != nil {
		s.logger.Error("something wrong with UOW creating", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unit of work problem %v", op, err)
	}

	return nil
}

func (s *imageService) GetAll(ctx context.Context, limit, offset int) ([]domain.Image, error) {
	op := "services.imageService.GetAll"

	if limit <= 0 || offset < 0 {
		s.logger.Error("invalid parameter limit and offset", "op", op)
		return nil, fmt.Errorf("%s: %w", op, crud_errors.ErrInvalidParam)
	}

	images, err := s.reader.GetAll(ctx, limit, offset)
	if err != nil {
		if errors.Is(err, crud_errors.ErrNotFound) {
			s.logger.Debug("no content", "op", op)
		} else {
			s.logger.Error("extract data failed", logger.Err(err), "op", op)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return images, nil
}

func (s *imageService) GetById(ctx context.Context, id uuid.UUID) (*domain.Image, error) {
	op := "services.imageService.GetById"

	image, err := s.reader.GetById(ctx, id)

	if err != nil {
		if errors.Is(err, crud_errors.ErrNotFound) {
			s.logger.Debug("image not found", "op", op)
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		s.logger.Error("extract data failed", logger.Err(err), "op", op)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return image, nil
}

func (s *imageService) Update(ctx context.Context, image *domain.Image) error {
	op := "services.imageService.Update"

	if err := validateImage(image.Data); err != nil {
		s.logger.Error("image is corruption or taked data is not image", logger.Err(err), "op", op)
		return fmt.Errorf("%s: validation error: %v", op, err)
	}

	err := s.uow.Do(ctx, func(ctx context.Context, tx uow.Transaction) error {
		uowOp := op + ".uow"
		imageRepoGen, err := getReposiotry(tx, uow.ImageRepoName, s.logger)
		if err != nil {
			s.logger.Error("get image repository generator is unable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: get image repository generator is unable: %v", uowOp, err)
		}

		imageRepo := imageRepoGen.(*postgres.ImageRepo)

		image.Hash = generateImageHash(image.Data)

		if err := imageRepo.Update(ctx, image); err != nil {
			if errors.Is(err, crud_errors.ErrNotFound) {
				s.logger.Debug("update initialize is unable", "op", uowOp)
				return fmt.Errorf("%s: %w", uowOp, err)
			}

			s.logger.Error("failed to update image", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: failed to update image: %v", uowOp, err)
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, crud_errors.ErrNotFound) {
			s.logger.Warn("update initialize is unable: image not found", "op", op)
			return fmt.Errorf("%s: %w", op, err)
		}

		s.logger.Error("something wrong with UOW updating", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unit of work update problem: %w", op, err)
	}

	return nil
}

func (s *imageService) Delete(ctx context.Context, id uuid.UUID) error {
	op := "services.imageService.Delete"

	err := s.uow.Do(ctx, func(ctx context.Context, tx uow.Transaction) error {
		uowOp := op + ".uow"
		imageRepoGen, err := getReposiotry(tx, uow.ImageRepoName, s.logger)
		if err != nil {
			s.logger.Error("image transaction problem on creating", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: get image repository generator is unable: %v", uowOp, err)
		}

		imageRepo := imageRepoGen.(*postgres.ImageRepo)

		savepoint := `sp_delete_image`
		err = safeDelete(ctx, tx.GetTX(), id, imageRepo.Delete, s.logger, uowOp, savepoint)
		if err != nil {
			if errors.Is(err, crud_errors.ErrNotFound) {
				s.logger.Debug("image not found", "op", uowOp)
				return nil
			}

			s.logger.Error("unable to safe delete address", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: unable to safe delete image: %v", uowOp, err)
		}

		return nil
	})

	if err != nil {
		s.logger.Error("something wrong with UOW deleting", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unit of work delete problem: %w", op, err)
	}

	return nil
}
