package integration_test

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/config"
	"CRUD-HOME-APPLIANCE-STORE/internal/consul"
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

var (
	cfg *config.Config
	db  *pgx.Conn
)

func TestMain(m *testing.M) {
	cfg = config.MustLoad()
	err := consul.WaitForService(cfg)
	if err != nil {
		panic("Service not ready in Consul " + err.Error())
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.PostgresConfig.PostgresUser, cfg.PostgresConfig.PostgresPassword, cfg.PostgresConfig.PostgresHost, cfg.PostgresConfig.PostgresPort, cfg.PostgresConfig.PostgresDatabase)
	db, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		panic("Postgres not vailable")
	}

	os.Exit(m.Run())
}
