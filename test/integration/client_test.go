package integration_test

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/mapper"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
	"CRUD-HOME-APPLIANCE-STORE/internal/repositories/postgres"
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateClient(t *testing.T) {
	givedData := dto.Client{
		Name:     "Adrianna",
		Surname:  "Gopher",
		Birthday: "2001-01-01",
		Gender:   "female",
		Address: dto.Address{
			Country: "Japan",
			City:    "Tokyo",
			Street:  "Godzilla",
		},
	}

	send, err := json.Marshal(&givedData)
	if err != nil {
		panic("Something wrong with marshal struct")
	}

	url := fmt.Sprintf("http://%s:%s/api/v1/clients", cfg.CrudService.Address, cfg.CrudService.Port)
	resp, err := http.Post(url, "application/json", strings.NewReader(
		string(send),
	))

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	log := logger.NewLogger(cfg.Env)

	clientRepo := postgres.NewClientRepository(db, log)

	client, repoErr := clientRepo.GetByNameAndSurname(context.Background(), givedData.Name, givedData.Surname)
	if repoErr != nil {
		panic("Invalid test")
	}

	require.Equal(t, len(client), 1)

	basicData, err := mapper.ClientToDomain(givedData)
	if err != nil {
		panic("Basic data is invalid")
	}

	for _, c := range client {
		require.Equal(t, c.Name, basicData.Name)
		require.Equal(t, c.Surname, basicData.Surname)
		require.Equal(t, c.Birthday, basicData.Birthday)
		require.Equal(t, c.Gender, basicData.Gender)
		require.Equal(t, c.Address.Country, basicData.Address.Country)
		require.Equal(t, c.Address.City, basicData.Address.City)
		require.Equal(t, c.Address.Street, basicData.Address.Street)
	}
}

func TestCreateClient2(t *testing.T) {
	require.Equal(t, 1, 1)
}
