package postgres

import (
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type BaseRepository interface{}

type DB interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type basePostgresRepository struct {
	db     DB
	logger *logger.Logger
}

func newBasePostgresRepository(db DB, logger *logger.Logger) *basePostgresRepository {
	logger.Debug("postgres base repo is created")
	return &basePostgresRepository{
		db:     db,
		logger: logger,
	}
}
