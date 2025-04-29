package connection

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type PostgresConfig struct {
	PostgresHost     string        `env:"POSTGRES_HOST" env-default:"localhost"`
	PostgresPort     string        `env:"POSTGRES_PORT" env-default:"5432"`
	PostgresUser     string        `env:"POSTGRES_USER"`
	PostgresPassword string        `env:"POSTGRES_PASSWORD"`
	PostgresDatabase string        `env:"POSTGRES_DB"`
	MaxConn          string        `env:"postgres_db_pool_max_conns" env-default:"20"`
	ConnectTimeout   time.Duration `env:"timeout" env-default:"2s"`
}

var (
	ErrConnectTimeout = errors.New("connect timeout")
)

func NewPostgresStorage(cfg *PostgresConfig) (*pgx.Conn, error) {
	connCtx, cancel := context.WithTimeoutCause(context.Background(), cfg.ConnectTimeout, ErrConnectTimeout)
	defer cancel()

	// connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?pool_max_conns=%s", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDatabase, cfg.MaxConn)
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDatabase)
	conn, err := pgx.Connect(connCtx, connStr)
	if err != nil {
		return nil, fmt.Errorf("postgres connect: %v", err)
	}

	return conn, nil
}
