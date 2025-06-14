package integration

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

func (s *TestSuite) TestCreateClient() {
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
	s.Require().NoError(err)

	url := fmt.Sprintf("http://%s:%s/api/v1/clients", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	resp, err := http.Post(url, "application/json", strings.NewReader(
		string(send),
	))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	log := logger.NewLogger(s.cfg.Env)

	clientRepo := postgres.NewClientRepository(s.db, log)

	client, err := clientRepo.GetByNameAndSurname(context.Background(), givedData.Name, givedData.Surname)
	s.Require().NoError(err)

	s.Require().Len(client, 1)

	basicData, err := mapper.ClientToDomain(givedData)
	s.Require().NoError(err)

	for _, c := range client {
		s.Require().Equal(c.Name, basicData.Name)
		s.Require().Equal(c.Surname, basicData.Surname)
		s.Require().Equal(c.Birthday, basicData.Birthday)
		s.Require().Equal(c.Gender, basicData.Gender)
		s.Require().Equal(c.Address.Country, basicData.Address.Country)
		s.Require().Equal(c.Address.City, basicData.Address.City)
		s.Require().Equal(c.Address.Street, basicData.Address.Street)
	}
}

func TestCreateClient2(t *testing.T) {
	require.Equal(t, 1, 1)
}
