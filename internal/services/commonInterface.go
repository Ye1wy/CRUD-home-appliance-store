package services

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"context"

	"github.com/google/uuid"
)

type addressWriter interface {
	Create(ctx context.Context, address *domain.Address) error
	Delete(ctx context.Context, id uuid.UUID) error
}
