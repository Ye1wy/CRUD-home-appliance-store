package postgres

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"context"

	"github.com/google/uuid"
)

type ImageRepository struct {
	*basePostgresRepository
}

func NewImageRepository(db DB, logger *logger.Logger) *ImageRepository {
	repo := newBasePostgresRepository(db, logger)
	return &ImageRepository{
		repo,
	}
}

func (r *ImageRepository) Create(ctx context.Context, image domain.Image) error {
	return nil
}

func (r *ImageRepository) GetAll(ctx context.Context, limit, offset int) ([]domain.Product, error) {
	return nil, nil
}

func (r *ImageRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Image, error) {
	return nil, nil
}

func (r *ImageRepository) Update(ctx context.Context, image domain.Image) error {
	return nil
}

func (r *ImageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
