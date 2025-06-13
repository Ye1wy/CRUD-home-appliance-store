package main

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/config"
	"CRUD-HOME-APPLIANCE-STORE/internal/consul"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
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

func TestCreateClient(t *testing.T) {
	url := fmt.Sprintf("http://%s:%s/api/v1/clients", cfg.CrudService.Address, cfg.CrudService.Port)
	resp, err := http.Post(url, "application/json", strings.NewReader(
		`{
		"name": "Adrianna",
		"surname": "Gopher",
		"birthday": "2001-01-01",
		"gender": "female",
		"country": "Japan",
		"city": "Tokyo",
		"street": "Godzilla"
	}`,
	))

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
}
