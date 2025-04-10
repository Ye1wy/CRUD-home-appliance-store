package services

import (
	"context"

	"github.com/google/uuid"
)

type ImageRepository interface {
	Create(ctx context.Context) error
	GetAll(ctx context.Context) error
	GetById(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context) error
	Delete(ctx context.Context, id uuid.UUID) error
}
