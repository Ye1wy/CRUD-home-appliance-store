package integration

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/mapper"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
	"CRUD-HOME-APPLIANCE-STORE/internal/repositories/postgres"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func clientResponseToRequest(c dto.ClientResponse) dto.ClientRequest {
	output := dto.ClientRequest{
		Name:     c.Name,
		Surname:  c.Surname,
		Birthday: c.Birthday,
		Gender:   c.Gender,
	}

	if c.Address != nil {
		output.Address = c.Address
	}

	return output
}

func (s *TestSuite) TestCreateClient() {
	s.CleanTable()
	givedData := dto.ClientRequest{
		Name:     "Adrianna",
		Surname:  "Gopher",
		Birthday: "2001-01-01",
		Gender:   "female",
		Address: &dto.Address{
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

	clientRepo := postgres.NewClientRepository(s.db, s.logger)

	client, err := clientRepo.GetByNameAndSurname(context.Background(), givedData.Name, givedData.Surname)
	s.Require().NoError(err)

	s.Require().Len(client, 1)

	basicData, err := mapper.ClientRequestToDomain(givedData)
	s.Require().NoError(err)

	for _, c := range client {
		s.Require().Equal(c.Name, basicData.Name)
		s.Require().Equal(c.Surname, basicData.Surname)
		s.Require().Equal(c.Birthday, basicData.Birthday)
		s.Require().Equal(c.Gender, basicData.Gender)
		s.Require().Equal(c.Address.Country, basicData.Address.Country)
		s.Require().Equal(c.Address.City, basicData.Address.City)
		s.Require().Equal(c.Address.Street, basicData.Address.Street)
		s.Require().NotEmpty(c.Id)
	}
}

func (s *TestSuite) TestCreateClientWithoutAddress() {
	s.CleanTable()
	givedData := dto.ClientRequest{
		Name:     "Adrianna",
		Surname:  "Gopher",
		Birthday: "2001-01-01",
		Gender:   "female",
	}

	send, err := json.Marshal(&givedData)
	s.Require().NoError(err)

	url := fmt.Sprintf("http://%s:%s/api/v1/clients", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	resp, err := http.Post(url, "application/json", strings.NewReader(
		string(send),
	))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	clientRepo := postgres.NewClientRepository(s.db, s.logger)

	client, err := clientRepo.GetByNameAndSurname(context.Background(), givedData.Name, givedData.Surname)
	s.Require().NoError(err)

	s.Require().Len(client, 1)

	var checkData []dto.ClientResponse

	for _, c := range client {
		checkData = append(checkData, mapper.ClientDomainToClientResponse(c))
	}

	for _, c := range checkData {
		s.Require().Equal(givedData.Name, c.Name)
		s.Require().Equal(givedData.Surname, c.Surname)
		s.Require().Equal(givedData.Gender, c.Gender)
		s.Require().Equal(givedData.Birthday, c.Birthday)

		// check id is contains or not
		s.Require().NotEmpty(c.Id)
	}

	query := `SELECT COUNT(*) FROM address`

	var count int

	err = s.db.QueryRow(context.Background(), query).Scan(&count)
	s.Require().NoError(err)

	s.Require().EqualValues(0, count)

}

func (s *TestSuite) TestCreateClientWithOneAddress() {
	s.CleanTable()
	commonAddress := dto.Address{
		Country: "Japan",
		City:    "Tokyo",
		Street:  "Godzilla",
	}
	givedData := dto.ClientRequest{
		Name:     "Adrianna",
		Surname:  "Gopher",
		Birthday: "2001-01-01",
		Gender:   "female",
		Address:  &commonAddress,
	}

	send, err := json.Marshal(&givedData)
	s.Require().NoError(err)

	url := fmt.Sprintf("http://%s:%s/api/v1/clients", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	resp, err := http.Post(url, "application/json", strings.NewReader(
		string(send),
	))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	givedDataAnother := dto.ClientRequest{
		Name:     "Adrian",
		Surname:  "Gopher",
		Birthday: "2005-01-01",
		Gender:   "male",
		Address:  &commonAddress,
	}

	send, err = json.Marshal(&givedDataAnother)
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

	var clients []dto.ClientResponse
	err = json.Unmarshal(body, &clients)
	s.Require().NoError(err)
	s.Require().Len(clients, 2)

	var check []dto.ClientRequest

	for _, c := range clients {
		tmp := clientResponseToRequest(c)
		check = append(check, tmp)
	}

	s.Require().Contains(check, givedData)
	s.Require().Contains(check, givedDataAnother)

	// check id is contains or not
	for _, c := range clients {
		s.Require().NotEmpty(c.Id)
	}
}

func (s *TestSuite) TestGetClient() {
	s.CleanTable()
	commonAddress := dto.Address{
		Country: "Japan",
		City:    "Tokyo",
		Street:  "Godzilla",
	}
	givedData := dto.ClientRequest{
		Name:     "Adrianna",
		Surname:  "Gopher",
		Birthday: "2001-01-01",
		Gender:   "female",
		Address:  &commonAddress,
	}

	send, err := json.Marshal(&givedData)
	s.Require().NoError(err)

	url := fmt.Sprintf("http://%s:%s/api/v1/clients", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	resp, err := http.Post(url, "application/json", strings.NewReader(
		string(send),
	))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	givedDataAnother := dto.ClientRequest{
		Name:     "Adrian",
		Surname:  "Gopher",
		Birthday: "2005-01-01",
		Gender:   "male",
		Address:  &commonAddress,
	}

	send, err = json.Marshal(&givedDataAnother)
	s.Require().NoError(err)

	url = fmt.Sprintf("http://%s:%s/api/v1/clients", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	resp, err = http.Post(url, "application/json", strings.NewReader(
		string(send),
	))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	url = fmt.Sprintf("http://%s:%s/api/v1/clients", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	resp, err = http.Get(url)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Require().Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	var clients []dto.ClientResponse
	err = json.Unmarshal(body, &clients)
	s.Require().NoError(err)
	s.Require().Len(clients, 2)

	var check []dto.ClientRequest

	for _, c := range clients {
		tmp := clientResponseToRequest(c)
		check = append(check, tmp)
	}

	s.Require().Contains(check, givedData)
	s.Require().Contains(check, givedDataAnother)

	// check id is contains or not
	for _, c := range clients {
		s.Require().NotEmpty(c.Id)
	}
}

func (s *TestSuite) TestGetClientByNameAndSurname() {
	s.CleanTable()
	commonAddress := dto.Address{
		Country: "Japan",
		City:    "Tokyo",
		Street:  "Godzilla",
	}
	givedData := dto.ClientRequest{
		Name:     "Adrianna",
		Surname:  "Gopher",
		Birthday: "2001-01-01",
		Gender:   "female",
		Address:  &commonAddress,
	}

	send, err := json.Marshal(&givedData)
	s.Require().NoError(err)

	url := fmt.Sprintf("http://%s:%s/api/v1/clients", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	resp, err := http.Post(url, "application/json", strings.NewReader(
		string(send),
	))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	givedDataAnother := dto.ClientRequest{
		Name:     "Adrian",
		Surname:  "Gopher",
		Birthday: "2005-01-01",
		Gender:   "male",
		Address:  &commonAddress,
	}

	send, err = json.Marshal(&givedDataAnother)
	s.Require().NoError(err)

	url = fmt.Sprintf("http://%s:%s/api/v1/clients", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	resp, err = http.Post(url, "application/json", strings.NewReader(
		string(send),
	))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	url = fmt.Sprintf("http://%s:%s/api/v1/clients/search?name=%s&surname=%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, givedData.Name, givedData.Surname)
	resp, err = http.Get(url)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Require().Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	var clients []dto.ClientResponse
	err = json.Unmarshal(body, &clients)
	s.Require().NoError(err)
	s.Require().Len(clients, 1)

	var check []dto.ClientRequest

	for _, c := range clients {
		check = append(check, clientResponseToRequest(c))
	}

	s.Require().Contains(check, givedData)

	for _, c := range clients {
		s.Require().NotEmpty(c.Id)
	}
}

func (s *TestSuite) TestGetClientByNameAndSurnameMulty() {
	s.CleanTable()
	commonAddress := dto.Address{
		Country: "Japan",
		City:    "Tokyo",
		Street:  "Godzilla",
	}
	first := dto.ClientRequest{
		Name:     "Adrianna",
		Surname:  "Gopher",
		Birthday: "2001-01-01",
		Gender:   "female",
		Address:  &commonAddress,
	}

	send, err := json.Marshal(&first)
	s.Require().NoError(err)

	url := fmt.Sprintf("http://%s:%s/api/v1/clients", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	resp, err := http.Post(url, "application/json", strings.NewReader(
		string(send),
	))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	second := dto.ClientRequest{
		Name:     "Adrian",
		Surname:  "Gopher",
		Birthday: "2005-01-01",
		Gender:   "male",
		Address:  &commonAddress,
	}

	send, err = json.Marshal(&second)
	s.Require().NoError(err)

	url = fmt.Sprintf("http://%s:%s/api/v1/clients", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	resp, err = http.Post(url, "application/json", strings.NewReader(
		string(send),
	))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	third := dto.ClientRequest{
		Name:     "Adrian",
		Surname:  "Gopher",
		Birthday: "2001-01-01",
		Gender:   "male",
		Address:  &commonAddress,
	}

	send, err = json.Marshal(&third)
	s.Require().NoError(err)

	url = fmt.Sprintf("http://%s:%s/api/v1/clients", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	resp, err = http.Post(url, "application/json", strings.NewReader(
		string(send),
	))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	name := "Adrian"
	surname := "Gopher"

	url = fmt.Sprintf("http://%s:%s/api/v1/clients/search?name=%s&surname=%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, name, surname)
	resp, err = http.Get(url)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Require().Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	var clients []dto.ClientResponse
	err = json.Unmarshal(body, &clients)
	s.Require().NoError(err)
	s.Require().Len(clients, 2)
	var check []dto.ClientRequest

	for _, c := range clients {
		check = append(check, clientResponseToRequest(c))
	}

	s.Require().Contains(check, second)
	s.Require().Contains(check, third)

	for _, c := range clients {
		s.Require().NotEmpty(c.Id)
	}
}

func (s *TestSuite) TestClientUpdateAddress() {
	s.CleanTable()

	first := dto.ClientRequest{
		Name:     "Adrianna",
		Surname:  "Gopher",
		Birthday: "2001-01-01",
		Gender:   "female",
		Address: &dto.Address{
			Country: "Japan",
			City:    "Tokyo",
			Street:  "Godzilla",
		},
	}

	second := dto.ClientRequest{
		Name:     "Adrian",
		Surname:  "Gopher",
		Birthday: "2005-01-01",
		Gender:   "male",
		Address: &dto.Address{
			Country: "Japan",
			City:    "Tokyo",
			Street:  "Godzilla",
		},
	}

	third := dto.ClientRequest{
		Name:     "Kazui",
		Surname:  "Franclin",
		Birthday: "2001-01-01",
		Gender:   "male",
	}

	dataBank := []dto.ClientRequest{first, second, third}

	for _, data := range dataBank {
		send, err := json.Marshal(&data)
		s.Require().NoError(err)

		url := fmt.Sprintf("http://%s:%s/api/v1/clients", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
		resp, err := http.Post(url, "application/json", strings.NewReader(
			string(send),
		))

		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, resp.StatusCode)
	}

	clientRepo := postgres.NewClientRepository(s.db, s.logger)

	testClients, err := clientRepo.GetByNameAndSurname(context.Background(), second.Name, second.Surname)
	s.Require().NoError(err)

	s.Require().Len(testClients, 1)

	neededId := testClients[0].Id

	url := fmt.Sprintf("http://%s:%s/api/v1/clients/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, neededId.String())
	payload, err := json.Marshal(dto.Address{
		Country: "Korea",
		City:    "Seoul",
		Street:  "Gangnam",
	})
	s.Require().NoError(err)

	httpClient := http.Client{}
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(payload))
	req.Header.Set("Content-type", "application/json")
	s.Require().NoError(err)

	resp, err := httpClient.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	query := `SELECT COUNT(id) FROM address`

	var count int
	err = s.db.QueryRow(context.Background(), query).Scan(&count)
	s.Require().NoError(err)
	s.Require().EqualValues(2, count)

	query = `SELECT * FROM address WHERE country=@newCountry AND city=@newCity AND street=@newStreet`
	args := pgx.NamedArgs{
		"newCountry": "Korea",
		"newCity":    "Seoul",
		"newStreet":  "Gangnam",
	}

	temp := domain.Address{}
	err = s.db.QueryRow(context.Background(), query, args).Scan(
		&temp.Id,
		&temp.Country,
		&temp.City,
		&temp.Street,
	)
	s.Require().NoError(err)

	check := mapper.AddressToDto(temp)
	s.Require().Equal(check, dto.Address{
		Country: "Korea",
		City:    "Seoul",
		Street:  "Gangnam",
	})

	url = fmt.Sprintf("http://%s:%s/api/v1/clients/search?name=%s&surname=%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, second.Name, second.Surname)
	resp, err = http.Get(url)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Require().Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	var clients []dto.ClientResponse
	err = json.Unmarshal(body, &clients)
	s.Require().NoError(err)
	var checkClient []dto.ClientRequest

	for _, c := range clients {
		checkClient = append(checkClient, clientResponseToRequest(c))
	}

	s.Require().Contains(checkClient, dto.ClientRequest{
		Name:     "Adrian",
		Surname:  "Gopher",
		Birthday: "2005-01-01",
		Gender:   "male",
		Address: &dto.Address{
			Country: "Korea",
			City:    "Seoul",
			Street:  "Gangnam",
		},
	})

}

func (s *TestSuite) TestClientUpdateAddressOnEmpty() {
	s.CleanTable()

	first := dto.ClientRequest{
		Name:     "Adrianna",
		Surname:  "Gopher",
		Birthday: "2001-01-01",
		Gender:   "female",
		Address: &dto.Address{
			Country: "Japan",
			City:    "Tokyo",
			Street:  "Godzilla",
		},
	}

	second := dto.ClientRequest{
		Name:     "Adrian",
		Surname:  "Gopher",
		Birthday: "2005-01-01",
		Gender:   "male",
		Address: &dto.Address{
			Country: "Japan",
			City:    "Tokyo",
			Street:  "Godzilla",
		},
	}

	third := dto.ClientRequest{
		Name:     "Kazui",
		Surname:  "Franclin",
		Birthday: "2001-01-01",
		Gender:   "male",
	}

	dataBank := []dto.ClientRequest{first, second, third}

	for _, data := range dataBank {
		send, err := json.Marshal(&data)
		s.Require().NoError(err)

		url := fmt.Sprintf("http://%s:%s/api/v1/clients", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
		resp, err := http.Post(url, "application/json", strings.NewReader(
			string(send),
		))

		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, resp.StatusCode)
	}

	clientRepo := postgres.NewClientRepository(s.db, s.logger)

	testClients, err := clientRepo.GetByNameAndSurname(context.Background(), second.Name, second.Surname)
	s.Require().NoError(err)

	s.Require().Len(testClients, 1)

	neededId := testClients[0].Id

	url := fmt.Sprintf("http://%s:%s/api/v1/clients/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, neededId.String())
	payload, err := json.Marshal(dto.Address{
		Country: "",
		City:    "",
		Street:  "",
	})
	s.Require().NoError(err)

	httpClient := http.Client{}
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(payload))
	req.Header.Set("Content-type", "application/json")
	s.Require().NoError(err)

	resp, err := httpClient.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)

	query := `SELECT COUNT(id) FROM address`

	var count int
	err = s.db.QueryRow(context.Background(), query).Scan(&count)
	s.Require().NoError(err)
	s.Require().EqualValues(1, count)

	query = `SELECT * FROM address WHERE country=@newCountry AND city=@newCity AND street=@newStreet`
	args := pgx.NamedArgs{
		"newCountry": "Korea",
		"newCity":    "Seoul",
		"newStreet":  "Gangnam",
	}

	temp := domain.Address{}
	err = s.db.QueryRow(context.Background(), query, args).Scan(
		&temp.Id,
		&temp.Country,
		&temp.City,
		&temp.Street,
	)
	s.Require().ErrorIs(err, pgx.ErrNoRows)

	url = fmt.Sprintf("http://%s:%s/api/v1/clients/search?name=%s&surname=%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, second.Name, second.Surname)
	resp, err = http.Get(url)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Require().Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	var clients []dto.ClientResponse
	err = json.Unmarshal(body, &clients)
	s.Require().NoError(err)
	var checkClient []dto.ClientRequest

	for _, c := range clients {
		checkClient = append(checkClient, clientResponseToRequest(c))
	}

	s.Require().Contains(checkClient, second)
}

func (s *TestSuite) TestClientUpdateAddressOnEmpty_2() {
	s.CleanTable()

	first := dto.ClientRequest{
		Name:     "Adrianna",
		Surname:  "Gopher",
		Birthday: "2001-01-01",
		Gender:   "female",
		Address: &dto.Address{
			Country: "Japan",
			City:    "Tokyo",
			Street:  "Godzilla",
		},
	}

	second := dto.ClientRequest{
		Name:     "Adrian",
		Surname:  "Gopher",
		Birthday: "2005-01-01",
		Gender:   "male",
		Address: &dto.Address{
			Country: "Japan",
			City:    "Tokyo",
			Street:  "Godzilla",
		},
	}

	third := dto.ClientRequest{
		Name:     "Kazui",
		Surname:  "Franclin",
		Birthday: "2001-01-01",
		Gender:   "male",
	}

	dataBank := []dto.ClientRequest{first, second, third}

	for _, data := range dataBank {
		send, err := json.Marshal(&data)
		s.Require().NoError(err)

		url := fmt.Sprintf("http://%s:%s/api/v1/clients", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
		resp, err := http.Post(url, "application/json", strings.NewReader(
			string(send),
		))

		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, resp.StatusCode)
	}

	clientRepo := postgres.NewClientRepository(s.db, s.logger)

	testClients, err := clientRepo.GetByNameAndSurname(context.Background(), second.Name, second.Surname)
	s.Require().NoError(err)

	s.Require().Len(testClients, 1)

	neededId := testClients[0].Id

	url := fmt.Sprintf("http://%s:%s/api/v1/clients/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, neededId.String())

	httpClient := http.Client{}
	req, err := http.NewRequest(http.MethodPatch, url, nil)
	req.Header.Set("Content-type", "application/json")
	s.Require().NoError(err)

	resp, err := httpClient.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)

	query := `SELECT COUNT(id) FROM address`

	var count int
	err = s.db.QueryRow(context.Background(), query).Scan(&count)
	s.Require().NoError(err)
	s.Require().EqualValues(1, count)

	query = `SELECT * FROM address WHERE country=@newCountry AND city=@newCity AND street=@newStreet`
	args := pgx.NamedArgs{
		"newCountry": "Korea",
		"newCity":    "Seoul",
		"newStreet":  "Gangnam",
	}

	temp := domain.Address{}
	err = s.db.QueryRow(context.Background(), query, args).Scan(
		&temp.Id,
		&temp.Country,
		&temp.City,
		&temp.Street,
	)
	s.Require().ErrorIs(err, pgx.ErrNoRows)

	url = fmt.Sprintf("http://%s:%s/api/v1/clients/search?name=%s&surname=%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, second.Name, second.Surname)
	resp, err = http.Get(url)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Require().Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	var clients []dto.ClientResponse
	err = json.Unmarshal(body, &clients)
	s.Require().NoError(err)
	var checkClients []dto.ClientRequest

	for _, c := range clients {
		checkClients = append(checkClients, clientResponseToRequest(c))
	}

	s.Require().Contains(checkClients, second)
}

func (s *TestSuite) TestClientDelete() {
	s.CleanTable()

	first := dto.ClientRequest{
		Name:     "Adrianna",
		Surname:  "Gopher",
		Birthday: "2001-01-01",
		Gender:   "female",
		Address: &dto.Address{
			Country: "Japan",
			City:    "Tokyo",
			Street:  "Godzilla",
		},
	}

	second := dto.ClientRequest{
		Name:     "Adrian",
		Surname:  "Gopher",
		Birthday: "2005-01-01",
		Gender:   "male",
		Address: &dto.Address{
			Country: "Japan",
			City:    "Tokyo",
			Street:  "Godzilla",
		},
	}

	third := dto.ClientRequest{
		Name:     "Kazui",
		Surname:  "Franclin",
		Birthday: "2001-01-01",
		Gender:   "male",
	}

	dataBank := []dto.ClientRequest{first, second, third}

	for _, data := range dataBank {
		send, err := json.Marshal(&data)
		s.Require().NoError(err)

		url := fmt.Sprintf("http://%s:%s/api/v1/clients", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
		resp, err := http.Post(url, "application/json", strings.NewReader(
			string(send),
		))

		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, resp.StatusCode)
	}

	clientRepo := postgres.NewClientRepository(s.db, s.logger)

	testClients, err := clientRepo.GetByNameAndSurname(context.Background(), second.Name, second.Surname)
	s.Require().NoError(err)

	s.Require().Len(testClients, 1)

	neededId := testClients[0].Id

	url := fmt.Sprintf("http://%s:%s/api/v1/clients/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, neededId.String())

	httpClient := http.Client{}
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	req.Header.Set("Content-type", "application/json")
	s.Require().NoError(err)

	resp, err := httpClient.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNoContent, resp.StatusCode)

	query := `SELECT COUNT(id) FROM address`

	var count int
	err = s.db.QueryRow(context.Background(), query).Scan(&count)
	s.Require().NoError(err)
	s.Require().EqualValues(1, count)

	query = `SELECT * FROM address WHERE country=@newCountry AND city=@newCity AND street=@newStreet`
	args := pgx.NamedArgs{
		"newCountry": "Japan",
		"newCity":    "Tokyo",
		"newStreet":  "Godzilla",
	}

	temp := domain.Address{}
	err = s.db.QueryRow(context.Background(), query, args).Scan(
		&temp.Id,
		&temp.Country,
		&temp.City,
		&temp.Street,
	)
	s.Require().NoError(err)

	url = fmt.Sprintf("http://%s:%s/api/v1/clients", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	resp, err = http.Get(url)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Require().Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	var clients []dto.ClientRequest
	err = json.Unmarshal(body, &clients)
	s.Require().NoError(err)

	s.Require().Contains(clients, first)
	s.Require().Contains(clients, third)
}

func (s *TestSuite) TestClientDeleteGhost() {
	s.CleanTable()

	first := dto.ClientRequest{
		Name:     "Adrianna",
		Surname:  "Gopher",
		Birthday: "2001-01-01",
		Gender:   "female",
		Address: &dto.Address{
			Country: "Japan",
			City:    "Tokyo",
			Street:  "Godzilla",
		},
	}

	second := dto.ClientRequest{
		Name:     "Adrian",
		Surname:  "Gopher",
		Birthday: "2005-01-01",
		Gender:   "male",
		Address: &dto.Address{
			Country: "Japan",
			City:    "Tokyo",
			Street:  "Godzilla",
		},
	}

	third := dto.ClientRequest{
		Name:     "Kazui",
		Surname:  "Franclin",
		Birthday: "2001-01-01",
		Gender:   "male",
	}

	dataBank := []dto.ClientRequest{first, second, third}

	for _, data := range dataBank {
		send, err := json.Marshal(&data)
		s.Require().NoError(err)

		url := fmt.Sprintf("http://%s:%s/api/v1/clients", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
		resp, err := http.Post(url, "application/json", strings.NewReader(
			string(send),
		))

		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, resp.StatusCode)
	}

	randomId, err := uuid.NewRandom()
	s.Require().NoError(err)

	url := fmt.Sprintf("http://%s:%s/api/v1/clients/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, randomId.String())

	httpClient := http.Client{}
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	req.Header.Set("Content-type", "application/json")
	s.Require().NoError(err)

	resp, err := httpClient.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNoContent, resp.StatusCode)

	query := `SELECT COUNT(id) FROM address`

	var count int
	err = s.db.QueryRow(context.Background(), query).Scan(&count)
	s.Require().NoError(err)
	s.Require().EqualValues(1, count)

	query = `SELECT * FROM address WHERE country=@newCountry AND city=@newCity AND street=@newStreet`
	args := pgx.NamedArgs{
		"newCountry": "Japan",
		"newCity":    "Tokyo",
		"newStreet":  "Godzilla",
	}

	temp := domain.Address{}
	err = s.db.QueryRow(context.Background(), query, args).Scan(
		&temp.Id,
		&temp.Country,
		&temp.City,
		&temp.Street,
	)
	s.Require().NoError(err)

	url = fmt.Sprintf("http://%s:%s/api/v1/clients", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	resp, err = http.Get(url)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Require().Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	var clients []dto.ClientResponse
	err = json.Unmarshal(body, &clients)
	s.Require().NoError(err)
	var checkClients []dto.ClientRequest

	for _, c := range clients {
		checkClients = append(checkClients, clientResponseToRequest(c))
	}

	s.Require().Contains(checkClients, first)
	s.Require().Contains(checkClients, second)
	s.Require().Contains(checkClients, third)
}

func (s *TestSuite) TestClientDeleteWithAddress() {
	s.CleanTable()

	first := dto.ClientRequest{
		Name:     "Adrianna",
		Surname:  "Gopher",
		Birthday: "2001-01-01",
		Gender:   "female",
		Address: &dto.Address{
			Country: "Japan",
			City:    "Tokyo",
			Street:  "Godzilla",
		},
	}

	second := dto.ClientRequest{
		Name:     "Adrian",
		Surname:  "Gopher",
		Birthday: "2005-01-01",
		Gender:   "male",
		Address: &dto.Address{
			Country: "Japan",
			City:    "Tokyo",
			Street:  "Jingu-dori",
		},
	}

	third := dto.ClientRequest{
		Name:     "Kazui",
		Surname:  "Franclin",
		Birthday: "2001-01-01",
		Gender:   "male",
	}

	dataBank := []dto.ClientRequest{first, second, third}

	for _, data := range dataBank {
		send, err := json.Marshal(&data)
		s.Require().NoError(err)

		url := fmt.Sprintf("http://%s:%s/api/v1/clients", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
		resp, err := http.Post(url, "application/json", strings.NewReader(
			string(send),
		))

		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, resp.StatusCode)
	}

	clientRepo := postgres.NewClientRepository(s.db, s.logger)

	testClients, err := clientRepo.GetByNameAndSurname(context.Background(), second.Name, second.Surname)
	s.Require().NoError(err)

	s.Require().Len(testClients, 1)

	neededId := testClients[0].Id

	query := `SELECT COUNT(id) FROM address`

	var count int
	err = s.db.QueryRow(context.Background(), query).Scan(&count)
	s.Require().NoError(err)
	s.Require().EqualValues(2, count)

	url := fmt.Sprintf("http://%s:%s/api/v1/clients/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, neededId.String())

	httpClient := http.Client{}
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	req.Header.Set("Content-type", "application/json")
	s.Require().NoError(err)

	resp, err := httpClient.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNoContent, resp.StatusCode)

	query = `SELECT COUNT(id) FROM address`

	count = 0
	err = s.db.QueryRow(context.Background(), query).Scan(&count)
	s.Require().NoError(err)
	s.Require().EqualValues(1, count)

	query = `SELECT * FROM address WHERE country=@newCountry AND city=@newCity AND street=@newStreet`
	args := pgx.NamedArgs{
		"newCountry": "Japan",
		"newCity":    "Tokyo",
		"newStreet":  "Jingu-dori",
	}

	temp := domain.Address{}
	err = s.db.QueryRow(context.Background(), query, args).Scan(
		&temp.Id,
		&temp.Country,
		&temp.City,
		&temp.Street,
	)
	s.Require().ErrorIs(err, pgx.ErrNoRows)

	query = `SELECT * FROM address WHERE country=@newCountry AND city=@newCity AND street=@newStreet`
	args = pgx.NamedArgs{
		"newCountry": "Japan",
		"newCity":    "Tokyo",
		"newStreet":  "Godzilla",
	}

	temp = domain.Address{}
	err = s.db.QueryRow(context.Background(), query, args).Scan(
		&temp.Id,
		&temp.Country,
		&temp.City,
		&temp.Street,
	)
	s.Require().NoError(err)

	url = fmt.Sprintf("http://%s:%s/api/v1/clients", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	resp, err = http.Get(url)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Require().Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	var clients []dto.ClientRequest
	err = json.Unmarshal(body, &clients)
	s.Require().NoError(err)

	s.Require().Contains(clients, first)
	s.Require().Contains(clients, third)
}
