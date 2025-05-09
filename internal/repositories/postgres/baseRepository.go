package postgres

import (
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"

	"github.com/jackc/pgx/v5"
)

type BaseRepository interface{}

type basePostgresRepository struct {
	db     pgx.Tx
	logger *logger.Logger
}

func newBasePostgresRepository(db pgx.Tx, logger *logger.Logger) *basePostgresRepository {
	logger.Debug("Postgres base repo is created")
	return &basePostgresRepository{
		db:     db,
		logger: logger,
	}
}
