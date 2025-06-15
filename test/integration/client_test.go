package integration

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/mapper"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
	"CRUD-HOME-APPLIANCE-STORE/internal/repositories/postgres"
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
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

	s.CleanTable()
}

func (s *TestSuite) TestCreateClientWithOneAddress() {
	commonAddress := dto.Address{
		Country: "Japan",
		City:    "Tokyo",
		Street:  "Godzilla",
	}
	givedData := dto.Client{
		Name:     "Adrianna",
		Surname:  "Gopher",
		Birthday: "2001-01-01",
		Gender:   "female",
		Address:  commonAddress,
	}

	send, err := json.Marshal(&givedData)
	s.Require().NoError(err)

	url := fmt.Sprintf("http://%s:%s/api/v1/clients", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	resp, err := http.Post(url, "application/json", strings.NewReader(
		string(send),
	))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	givedData = dto.Client{
		Name:     "Adrian",
		Surname:  "Gopher",
		Birthday: "2005-01-01",
		Gender:   "male",
		Address:  commonAddress,
	}

	send, err = json.Marshal(&givedData)
	s.Require().NoError(err)

	url = fmt.Sprintf("http://%s:%s/api/v1/clients", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	resp, err = http.Post(url, "application/json", strings.NewReader(
		string(send),
	))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	query := `SELECT COUNT(id) FROM address WHERE country=@country AND city=@city AND street=@street`
	args := pgx.NamedArgs{
		"country": commonAddress.Country,
		"city":    commonAddress.City,
		"street":  commonAddress.Street,
	}

	var count int
	err = s.db.QueryRow(context.Background(), query, args).Scan(&count)
	s.Require().NoError(err)

	s.Require().EqualValues(1, count)

	url = fmt.Sprintf("http://%s:%s/api/v1/clients", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	resp, err = http.Get(url)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Require().Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	var clients []dto.Client
	err = json.Unmarshal(body, &clients)
	s.Require().NoError(err)
	s.Require().Len(clients, 2)

	s.Require().Contains(clients, dto.Client{
		Name:     "Adrianna",
		Surname:  "Gopher",
		Birthday: "2001-01-01",
		Gender:   "female",
		Address:  commonAddress,
	})
	s.Require().Contains(clients, dto.Client{
		Name:     "Adrian",
		Surname:  "Gopher",
		Birthday: "2005-01-01",
		Gender:   "male",
		Address:  commonAddress,
	})
}
