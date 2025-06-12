package main

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/config"
	"CRUD-HOME-APPLIANCE-STORE/internal/consul"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	cfg := config.MustLoad()
	err := consul.WaitForService(cfg)
	if err != nil {
		panic("Service not ready in Consul " + err.Error())
	}

	os.Exit(m.Run())
}

func TestCreateClient(t *testing.T) {
	resp, err := http.Post("http://crud-service:8080/api/v1/clients", "application/json", strings.NewReader(
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
