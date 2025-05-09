package postgres

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// 170250
type ClientWriter interface {
	Create(ctx context.Context, client domain.Client) error
	UpdateAddress(ctx context.Context, id, address uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	client ClientWriter
}

func NewRepository(tx pgx.Tx, log *logger.Logger) *repository {
	client := NewClientRepository(tx, log)
	return &repository{
		client: client,
	}
}

func (r *repository) ClientRepo() ClientWriter {
	return r.client
}
