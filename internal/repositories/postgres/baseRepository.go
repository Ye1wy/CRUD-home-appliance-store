package psgrep

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/logger"

	"github.com/jackc/pgx/v5"
)

type BaseRepository interface{}

type basePostgresRepository struct {
	db     *pgx.Conn
	logger *logger.Logger
}

func newBasePostgresRepository(db *pgx.Conn, logger *logger.Logger) *basePostgresRepository {
	logger.Debug("Postgres base repo is created")
	return &basePostgresRepository{
		db:     db,
		logger: logger,
	}
}
