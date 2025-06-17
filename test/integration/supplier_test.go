package integration

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/mapper"
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
)

func createSuppliers(data []dto.Supplier, url string) error {
	for _, s := range data {
		payload, err := json.Marshal(&s)
		if err != nil {
			return err
		}

		resp, err := http.Post(url, "application/json", strings.NewReader(string(payload)))
		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("status: %d", resp.StatusCode)
		}
	}

	return nil
}

func (s *TestSuite) TestCreateSupplier() {
	s.CleanTable()
	givedData := dto.Supplier{
		Name:        "Aboba Inc.",
		PhoneNumber: "8-800-555-35-35",
		Address: &dto.Address{
			Country: "Japan",
			City:    "Tokyo",
			Street:  "Godzilla",
		},
	}

	payload, err := json.Marshal(&givedData)
	s.Require().NoError(err)
	url := fmt.Sprintf("http://%s:%s/api/v1/suppliers", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	resp, err := http.Post(url, "application/json", strings.NewReader(
		string(payload),
	))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	url = fmt.Sprintf("http://%s:%s/api/v1/suppliers", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	getResp, err := http.Get(url)
	s.Require().NoError(err)
	defer getResp.Body.Close()

	s.Require().Equal(http.StatusOK, getResp.StatusCode)
	data, err := io.ReadAll(getResp.Body)
	s.Require().NoError(err)

	exctractedData := []dto.Supplier{}
	err = json.Unmarshal(data, &exctractedData)
	s.Require().NoError(err)

	s.Require().Contains(exctractedData, givedData)
}

func (s *TestSuite) TestCreateSupplierWithoutAddress() {
	s.CleanTable()
	givedData := dto.Supplier{
		Name:        "Aboba Inc.",
		PhoneNumber: "8-800-555-35-35",
	}

	payload, err := json.Marshal(&givedData)
	s.Require().NoError(err)
	url := fmt.Sprintf("http://%s:%s/api/v1/suppliers", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	resp, err := http.Post(url, "application/json", strings.NewReader(
		string(payload),
	))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)

	url = fmt.Sprintf("http://%s:%s/api/v1/suppliers", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	getResp, err := http.Get(url)
	s.Require().NoError(err)
	defer getResp.Body.Close()

	s.Require().Equal(http.StatusNotFound, getResp.StatusCode)
}

func (s *TestSuite) TestCreateEmptySupplier() {
	s.CleanTable()
	data := dto.Supplier{}

	payload, err := json.Marshal(data)
	s.Require().NoError(err)
	url := fmt.Sprintf("http://%s:%s/api/v1/suppliers", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	resp, err := http.Post(url, "application/json", strings.NewReader(string(payload)))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)

	getResp, err := http.Get(url)
	s.Require().NoError(err)
	defer getResp.Body.Close()

	s.Require().Equal(http.StatusNotFound, getResp.StatusCode)
}

func (s *TestSuite) TestCreateSupplierDuplicate() {
	s.CleanTable()
	supplier := dto.Supplier{
		Name:        "Aboba Inc.",
		PhoneNumber: "8-800-555-35-35",
		Address: &dto.Address{
			Country: "Japan",
			City:    "Tokyo",
			Street:  "Godzilla",
		},
	}

	duplicate := dto.Supplier{
		Name:        "Aboba Inc.",
		PhoneNumber: "8-800-555-35-35",
		Address: &dto.Address{
			Country: "Japan",
			City:    "Tokyo",
			Street:  "Godzilla",
		},
	}

	payload, err := json.Marshal(&supplier)
	s.Require().NoError(err)
	url := fmt.Sprintf("http://%s:%s/api/v1/suppliers", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	resp, err := http.Post(url, "application/json", strings.NewReader(
		string(payload),
	))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	payload, err = json.Marshal(&duplicate)
	s.Require().NoError(err)
	resp, err = http.Post(url, "application/json", strings.NewReader(
		string(payload),
	))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)

	url = fmt.Sprintf("http://%s:%s/api/v1/suppliers", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	getResp, err := http.Get(url)
	s.Require().NoError(err)
	defer getResp.Body.Close()

	s.Require().Equal(http.StatusOK, getResp.StatusCode)

	rawData, err := io.ReadAll(getResp.Body)
	s.Require().NoError(err)

	var check []dto.Supplier
	err = json.Unmarshal(rawData, &check)
	s.Require().NoError(err)

	s.Require().Len(check, 1)
}

func (s *TestSuite) TestCreateSupplierWithSameAddress() {
	s.CleanTable()
	first := dto.Supplier{
		Name:        "Aboba Inc.",
		PhoneNumber: "8-800-555-35-35",
		Address: &dto.Address{
			Country: "Japan",
			City:    "Tokyo",
			Street:  "Godzilla",
		},
	}

	second := dto.Supplier{
		Name:        "Aboba Tech Inc.",
		PhoneNumber: "8-800-532-35-35",
		Address: &dto.Address{
			Country: "Japan",
			City:    "Tokyo",
			Street:  "Godzilla",
		},
	}

	payload, err := json.Marshal(&first)
	s.Require().NoError(err)
	url := fmt.Sprintf("http://%s:%s/api/v1/suppliers", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	resp, err := http.Post(url, "application/json", strings.NewReader(
		string(payload),
	))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	payload, err = json.Marshal(&second)
	s.Require().NoError(err)
	resp, err = http.Post(url, "application/json", strings.NewReader(
		string(payload),
	))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	url = fmt.Sprintf("http://%s:%s/api/v1/suppliers", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	getResp, err := http.Get(url)
	s.Require().NoError(err)
	defer getResp.Body.Close()

	s.Require().Equal(http.StatusOK, getResp.StatusCode)

	rawData, err := io.ReadAll(getResp.Body)
	s.Require().NoError(err)

	var check []dto.Supplier
	err = json.Unmarshal(rawData, &check)
	s.Require().NoError(err)

	s.Require().Len(check, 2)
	s.Require().Contains(check, first)
	s.Require().Contains(check, second)

	query := `SELECT COUNT(id) FROM address`
	var count int
	err = s.db.QueryRow(context.Background(), query).Scan(&count)
	s.Require().NoError(err)
	s.Require().Equal(1, count)
}

func (s *TestSuite) TestCreateSupplierWithDifferentAddress() {
	s.CleanTable()
	first := dto.Supplier{
		Name:        "Aboba Inc.",
		PhoneNumber: "8-800-555-35-35",
		Address: &dto.Address{
			Country: "Japan",
			City:    "Tokyo",
			Street:  "Godzilla",
		},
	}

	second := dto.Supplier{
		Name:        "Aboba Tech Inc.",
		PhoneNumber: "8-800-532-35-35",
		Address: &dto.Address{
			Country: "Korea",
			City:    "Seoul",
			Street:  "Gangnam",
		},
	}

	payload, err := json.Marshal(&first)
	s.Require().NoError(err)
	url := fmt.Sprintf("http://%s:%s/api/v1/suppliers", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	resp, err := http.Post(url, "application/json", strings.NewReader(
		string(payload),
	))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	payload, err = json.Marshal(&second)
	s.Require().NoError(err)
	resp, err = http.Post(url, "application/json", strings.NewReader(
		string(payload),
	))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	url = fmt.Sprintf("http://%s:%s/api/v1/suppliers", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	getResp, err := http.Get(url)
	s.Require().NoError(err)
	defer getResp.Body.Close()

	s.Require().Equal(http.StatusOK, getResp.StatusCode)

	rawData, err := io.ReadAll(getResp.Body)
	s.Require().NoError(err)

	var check []dto.Supplier
	err = json.Unmarshal(rawData, &check)
	s.Require().NoError(err)

	s.Require().Len(check, 2)
	s.Require().Contains(check, first)
	s.Require().Contains(check, second)

	query := `SELECT COUNT(id) FROM address`
	var count int
	err = s.db.QueryRow(context.Background(), query).Scan(&count)
	s.Require().NoError(err)
	s.Require().Equal(2, count)
}

func (s *TestSuite) TestGetAllSupplier() {
	s.CleanTable()
	suppliers := []dto.Supplier{
		{
			Name:        "Global Supplies Inc.",
			PhoneNumber: "+1-202-555-0123",
			Address: &dto.Address{
				Country: "USA",
				City:    "New York",
				Street:  "5th Avenue, 101",
			},
		},
		{
			Name:        "Berlin Tech Parts",
			PhoneNumber: "+49-30-1234567",
			Address: &dto.Address{
				Country: "Germany",
				City:    "Berlin",
				Street:  "Alexanderplatz 5",
			},
		},
		{
			Name:        "Tokyo Machinery Co.",
			PhoneNumber: "+81-3-1234-5678",
			Address: &dto.Address{
				Country: "Japan",
				City:    "Tokyo",
				Street:  "Chiyoda 2-1-1",
			},
		},
		{
			Name:        "Paris Electronics",
			PhoneNumber: "+33-1-2345-6789",
			Address: &dto.Address{
				Country: "France",
				City:    "Paris",
				Street:  "Rue de Rivoli 77",
			},
		},
		{
			Name:        "Sydney Auto Parts",
			PhoneNumber: "+61-2-9876-5432",
			Address: &dto.Address{
				Country: "Australia",
				City:    "Sydney",
				Street:  "George St 55",
			},
		},
		{
			Name:        "London Office Supplies",
			PhoneNumber: "+44-20-7946-0958",
			Address: &dto.Address{
				Country: "UK",
				City:    "London",
				Street:  "Baker Street 221B",
			},
		},
		{
			Name:        "Moscow Tools Ltd.",
			PhoneNumber: "+7-495-123-4567",
			Address: &dto.Address{
				Country: "Russia",
				City:    "Moscow",
				Street:  "Arbat 12",
			},
		},
		{
			Name:        "Toronto Packaging",
			PhoneNumber: "+1-416-555-7890",
			Address: &dto.Address{
				Country: "Canada",
				City:    "Toronto",
				Street:  "King St W 300",
			},
		},
		{
			Name:        "Beijing Textiles",
			PhoneNumber: "+86-10-1234-5678",
			Address: &dto.Address{
				Country: "China",
				City:    "Beijing",
				Street:  "Chang’an Ave 200",
			},
		},
		{
			Name:        "Delhi Agro Export",
			PhoneNumber: "+91-11-2345-6789",
			Address: &dto.Address{
				Country: "India",
				City:    "Delhi",
				Street:  "Connaught Place 19",
			},
		},
		{
			Name:        "Rome Steelworks",
			PhoneNumber: "+39-06-1234-5678",
			Address: &dto.Address{
				Country: "Italy",
				City:    "Rome",
				Street:  "Via del Corso 15",
			},
		},
		{
			Name:        "Madrid Chemicals",
			PhoneNumber: "+34-91-123-4567",
			Address: &dto.Address{
				Country: "Spain",
				City:    "Madrid",
				Street:  "Gran Via 42",
			},
		},
		{
			Name:        "São Paulo Imports",
			PhoneNumber: "+55-11-91234-5678",
			Address: &dto.Address{
				Country: "Brazil",
				City:    "São Paulo",
				Street:  "Avenida Paulista 1000",
			},
		},
		{
			Name:        "Seoul Electronics Hub",
			PhoneNumber: "+82-2-555-1234",
			Address: &dto.Address{
				Country: "South Korea",
				City:    "Seoul",
				Street:  "Gangnam-daero 432",
			},
		},
		{
			Name:        "Cape Town Minerals",
			PhoneNumber: "+27-21-123-4567",
			Address: &dto.Address{
				Country: "South Africa",
				City:    "Cape Town",
				Street:  "Long Street 88",
			},
		},
		{
			Name:        "Amsterdam Bikes Co.",
			PhoneNumber: "+31-20-123-4567",
			Address: &dto.Address{
				Country: "Netherlands",
				City:    "Amsterdam",
				Street:  "Damrak 45",
			},
		},
		{
			Name:        "Zurich Precision Tools",
			PhoneNumber: "+41-44-123-4567",
			Address: &dto.Address{
				Country: "Switzerland",
				City:    "Zurich",
				Street:  "Bahnhofstrasse 10",
			},
		},
		{
			Name:        "Vienna Food Logistics",
			PhoneNumber: "+43-1-2345-678",
			Address: &dto.Address{
				Country: "Austria",
				City:    "Vienna",
				Street:  "Mariahilfer Strasse 23",
			},
		},
		{
			Name:        "Stockholm CleanTech",
			PhoneNumber: "+46-8-1234-5678",
			Address: &dto.Address{
				Country: "Sweden",
				City:    "Stockholm",
				Street:  "Sveavägen 12",
			},
		},
		{
			Name:        "Helsinki Timber Group",
			PhoneNumber: "+358-9-1234-567",
			Address: &dto.Address{
				Country: "Finland",
				City:    "Helsinki",
				Street:  "Mannerheimintie 10",
			},
		},
	}

	url := fmt.Sprintf("http://%s:%s/api/v1/suppliers", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	err := createSuppliers(suppliers, url)
	s.Require().NoError(err)

	resp, err := http.Get(url)
	s.Require().NoError(err)
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	var check []dto.Supplier
	err = json.Unmarshal(data, &check)
	s.Require().NoError(err)
	s.Require().Len(check, 10)

	for i := range 10 {
		s.Require().Contains(check, suppliers[i])
	}
}

func (s *TestSuite) TestGetAllSupplierWithLimitAndOffset() {
	s.CleanTable()
	suppliers := []dto.Supplier{
		{
			Name:        "Global Supplies Inc.",
			PhoneNumber: "+1-202-555-0123",
			Address: &dto.Address{
				Country: "USA",
				City:    "New York",
				Street:  "5th Avenue, 101",
			},
		},
		{
			Name:        "Berlin Tech Parts",
			PhoneNumber: "+49-30-1234567",
			Address: &dto.Address{
				Country: "Germany",
				City:    "Berlin",
				Street:  "Alexanderplatz 5",
			},
		},
		{
			Name:        "Tokyo Machinery Co.",
			PhoneNumber: "+81-3-1234-5678",
			Address: &dto.Address{
				Country: "Japan",
				City:    "Tokyo",
				Street:  "Chiyoda 2-1-1",
			},
		},
		{
			Name:        "Paris Electronics",
			PhoneNumber: "+33-1-2345-6789",
			Address: &dto.Address{
				Country: "France",
				City:    "Paris",
				Street:  "Rue de Rivoli 77",
			},
		},
		{
			Name:        "Sydney Auto Parts",
			PhoneNumber: "+61-2-9876-5432",
			Address: &dto.Address{
				Country: "Australia",
				City:    "Sydney",
				Street:  "George St 55",
			},
		},
		{
			Name:        "London Office Supplies",
			PhoneNumber: "+44-20-7946-0958",
			Address: &dto.Address{
				Country: "UK",
				City:    "London",
				Street:  "Baker Street 221B",
			},
		},
		{
			Name:        "Moscow Tools Ltd.",
			PhoneNumber: "+7-495-123-4567",
			Address: &dto.Address{
				Country: "Russia",
				City:    "Moscow",
				Street:  "Arbat 12",
			},
		},
		{
			Name:        "Toronto Packaging",
			PhoneNumber: "+1-416-555-7890",
			Address: &dto.Address{
				Country: "Canada",
				City:    "Toronto",
				Street:  "King St W 300",
			},
		},
		{
			Name:        "Beijing Textiles",
			PhoneNumber: "+86-10-1234-5678",
			Address: &dto.Address{
				Country: "China",
				City:    "Beijing",
				Street:  "Chang’an Ave 200",
			},
		},
		{
			Name:        "Delhi Agro Export",
			PhoneNumber: "+91-11-2345-6789",
			Address: &dto.Address{
				Country: "India",
				City:    "Delhi",
				Street:  "Connaught Place 19",
			},
		},
		{
			Name:        "Rome Steelworks",
			PhoneNumber: "+39-06-1234-5678",
			Address: &dto.Address{
				Country: "Italy",
				City:    "Rome",
				Street:  "Via del Corso 15",
			},
		},
		{
			Name:        "Madrid Chemicals",
			PhoneNumber: "+34-91-123-4567",
			Address: &dto.Address{
				Country: "Spain",
				City:    "Madrid",
				Street:  "Gran Via 42",
			},
		},
		{
			Name:        "São Paulo Imports",
			PhoneNumber: "+55-11-91234-5678",
			Address: &dto.Address{
				Country: "Brazil",
				City:    "São Paulo",
				Street:  "Avenida Paulista 1000",
			},
		},
		{
			Name:        "Seoul Electronics Hub",
			PhoneNumber: "+82-2-555-1234",
			Address: &dto.Address{
				Country: "South Korea",
				City:    "Seoul",
				Street:  "Gangnam-daero 432",
			},
		},
		{
			Name:        "Cape Town Minerals",
			PhoneNumber: "+27-21-123-4567",
			Address: &dto.Address{
				Country: "South Africa",
				City:    "Cape Town",
				Street:  "Long Street 88",
			},
		},
		{
			Name:        "Amsterdam Bikes Co.",
			PhoneNumber: "+31-20-123-4567",
			Address: &dto.Address{
				Country: "Netherlands",
				City:    "Amsterdam",
				Street:  "Damrak 45",
			},
		},
		{
			Name:        "Zurich Precision Tools",
			PhoneNumber: "+41-44-123-4567",
			Address: &dto.Address{
				Country: "Switzerland",
				City:    "Zurich",
				Street:  "Bahnhofstrasse 10",
			},
		},
		{
			Name:        "Vienna Food Logistics",
			PhoneNumber: "+43-1-2345-678",
			Address: &dto.Address{
				Country: "Austria",
				City:    "Vienna",
				Street:  "Mariahilfer Strasse 23",
			},
		},
		{
			Name:        "Stockholm CleanTech",
			PhoneNumber: "+46-8-1234-5678",
			Address: &dto.Address{
				Country: "Sweden",
				City:    "Stockholm",
				Street:  "Sveavägen 12",
			},
		},
		{
			Name:        "Helsinki Timber Group",
			PhoneNumber: "+358-9-1234-567",
			Address: &dto.Address{
				Country: "Finland",
				City:    "Helsinki",
				Street:  "Mannerheimintie 10",
			},
		},
	}

	url := fmt.Sprintf("http://%s:%s/api/v1/suppliers", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	err := createSuppliers(suppliers, url)
	s.Require().NoError(err)

	url = fmt.Sprintf("http://%s:%s/api/v1/suppliers?limit=5&offset=3", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	resp, err := http.Get(url)
	s.Require().NoError(err)
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	var check []dto.Supplier
	err = json.Unmarshal(data, &check)
	s.Require().NoError(err)
	s.Require().Len(check, 5)

	for i := range 5 {
		s.Require().Contains(check, suppliers[i+3])
	}
}

func (s *TestSuite) TestGetByIdSupplier() {
	s.CleanTable()
	suppliers := []dto.Supplier{
		{
			Name:        "Global Supplies Inc.",
			PhoneNumber: "+1-202-555-0123",
			Address: &dto.Address{
				Country: "USA",
				City:    "New York",
				Street:  "5th Avenue, 101",
			},
		},
		{
			Name:        "Berlin Tech Parts",
			PhoneNumber: "+49-30-1234567",
			Address: &dto.Address{
				Country: "Germany",
				City:    "Berlin",
				Street:  "Alexanderplatz 5",
			},
		},
		{
			Name:        "Tokyo Machinery Co.",
			PhoneNumber: "+81-3-1234-5678",
			Address: &dto.Address{
				Country: "Japan",
				City:    "Tokyo",
				Street:  "Chiyoda 2-1-1",
			},
		},
		{
			Name:        "Paris Electronics",
			PhoneNumber: "+33-1-2345-6789",
			Address: &dto.Address{
				Country: "France",
				City:    "Paris",
				Street:  "Rue de Rivoli 77",
			},
		},
		{
			Name:        "Sydney Auto Parts",
			PhoneNumber: "+61-2-9876-5432",
			Address: &dto.Address{
				Country: "Australia",
				City:    "Sydney",
				Street:  "George St 55",
			},
		},
		{
			Name:        "London Office Supplies",
			PhoneNumber: "+44-20-7946-0958",
			Address: &dto.Address{
				Country: "UK",
				City:    "London",
				Street:  "Baker Street 221B",
			},
		},
		{
			Name:        "Moscow Tools Ltd.",
			PhoneNumber: "+7-495-123-4567",
			Address: &dto.Address{
				Country: "Russia",
				City:    "Moscow",
				Street:  "Arbat 12",
			},
		},
		{
			Name:        "Toronto Packaging",
			PhoneNumber: "+1-416-555-7890",
			Address: &dto.Address{
				Country: "Canada",
				City:    "Toronto",
				Street:  "King St W 300",
			},
		},
		{
			Name:        "Beijing Textiles",
			PhoneNumber: "+86-10-1234-5678",
			Address: &dto.Address{
				Country: "China",
				City:    "Beijing",
				Street:  "Chang’an Ave 200",
			},
		},
		{
			Name:        "Delhi Agro Export",
			PhoneNumber: "+91-11-2345-6789",
			Address: &dto.Address{
				Country: "India",
				City:    "Delhi",
				Street:  "Connaught Place 19",
			},
		},
		{
			Name:        "Rome Steelworks",
			PhoneNumber: "+39-06-1234-5678",
			Address: &dto.Address{
				Country: "Italy",
				City:    "Rome",
				Street:  "Via del Corso 15",
			},
		},
		{
			Name:        "Madrid Chemicals",
			PhoneNumber: "+34-91-123-4567",
			Address: &dto.Address{
				Country: "Spain",
				City:    "Madrid",
				Street:  "Gran Via 42",
			},
		},
		{
			Name:        "São Paulo Imports",
			PhoneNumber: "+55-11-91234-5678",
			Address: &dto.Address{
				Country: "Brazil",
				City:    "São Paulo",
				Street:  "Avenida Paulista 1000",
			},
		},
	}

	url := fmt.Sprintf("http://%s:%s/api/v1/suppliers", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	err := createSuppliers(suppliers, url)
	s.Require().NoError(err)

	supRepo := postgres.NewSupplierRepository(s.db, s.logger)
	sup, err := supRepo.GetByName(context.Background(), "São Paulo Imports")
	s.Require().NoError(err)

	getUrl := fmt.Sprintf("http://%s:%s/api/v1/suppliers/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, sup.Id.String())
	resp, err := http.Get(getUrl)
	s.Require().NoError(err)
	defer resp.Body.Close()

	rawData, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	var check dto.Supplier
	err = json.Unmarshal(rawData, &check)
	s.Require().NoError(err)

	s.Require().Equal(suppliers[len(suppliers)-1], check)
}

func (s *TestSuite) TestGetByIdGhostSupplier() {
	s.CleanTable()
	suppliers := []dto.Supplier{
		{
			Name:        "Global Supplies Inc.",
			PhoneNumber: "+1-202-555-0123",
			Address: &dto.Address{
				Country: "USA",
				City:    "New York",
				Street:  "5th Avenue, 101",
			},
		},
		{
			Name:        "Berlin Tech Parts",
			PhoneNumber: "+49-30-1234567",
			Address: &dto.Address{
				Country: "Germany",
				City:    "Berlin",
				Street:  "Alexanderplatz 5",
			},
		},
		{
			Name:        "Tokyo Machinery Co.",
			PhoneNumber: "+81-3-1234-5678",
			Address: &dto.Address{
				Country: "Japan",
				City:    "Tokyo",
				Street:  "Chiyoda 2-1-1",
			},
		},
		{
			Name:        "Paris Electronics",
			PhoneNumber: "+33-1-2345-6789",
			Address: &dto.Address{
				Country: "France",
				City:    "Paris",
				Street:  "Rue de Rivoli 77",
			},
		},
		{
			Name:        "Sydney Auto Parts",
			PhoneNumber: "+61-2-9876-5432",
			Address: &dto.Address{
				Country: "Australia",
				City:    "Sydney",
				Street:  "George St 55",
			},
		},
		{
			Name:        "London Office Supplies",
			PhoneNumber: "+44-20-7946-0958",
			Address: &dto.Address{
				Country: "UK",
				City:    "London",
				Street:  "Baker Street 221B",
			},
		},
		{
			Name:        "Moscow Tools Ltd.",
			PhoneNumber: "+7-495-123-4567",
			Address: &dto.Address{
				Country: "Russia",
				City:    "Moscow",
				Street:  "Arbat 12",
			},
		},
		{
			Name:        "Toronto Packaging",
			PhoneNumber: "+1-416-555-7890",
			Address: &dto.Address{
				Country: "Canada",
				City:    "Toronto",
				Street:  "King St W 300",
			},
		},
		{
			Name:        "Beijing Textiles",
			PhoneNumber: "+86-10-1234-5678",
			Address: &dto.Address{
				Country: "China",
				City:    "Beijing",
				Street:  "Chang’an Ave 200",
			},
		},
		{
			Name:        "Delhi Agro Export",
			PhoneNumber: "+91-11-2345-6789",
			Address: &dto.Address{
				Country: "India",
				City:    "Delhi",
				Street:  "Connaught Place 19",
			},
		},
		{
			Name:        "Rome Steelworks",
			PhoneNumber: "+39-06-1234-5678",
			Address: &dto.Address{
				Country: "Italy",
				City:    "Rome",
				Street:  "Via del Corso 15",
			},
		},
		{
			Name:        "Madrid Chemicals",
			PhoneNumber: "+34-91-123-4567",
			Address: &dto.Address{
				Country: "Spain",
				City:    "Madrid",
				Street:  "Gran Via 42",
			},
		},
		{
			Name:        "São Paulo Imports",
			PhoneNumber: "+55-11-91234-5678",
			Address: &dto.Address{
				Country: "Brazil",
				City:    "São Paulo",
				Street:  "Avenida Paulista 1000",
			},
		},
	}

	url := fmt.Sprintf("http://%s:%s/api/v1/suppliers", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	err := createSuppliers(suppliers, url)
	s.Require().NoError(err)

	randomId, err := uuid.NewRandom()
	s.Require().NoError(err)

	getUrl := fmt.Sprintf("http://%s:%s/api/v1/suppliers/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, randomId.String())
	resp, err := http.Get(getUrl)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *TestSuite) TestGetByIdSupplierInvalidUUID() {
	s.CleanTable()
	getUrl := fmt.Sprintf("http://%s:%s/api/v1/suppliers/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, "312321")
	resp, err := http.Get(getUrl)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *TestSuite) TestUpdateAddressSupplier() {
	s.CleanTable()
	supplier := dto.Supplier{
		Name:        "Aboba Tech Inc.",
		PhoneNumber: "1232913",
		Address: &dto.Address{
			Country: "Korea",
			City:    "Seoul",
			Street:  "Myeongdong",
		},
	}

	postUrl := fmt.Sprintf("http://%s:%s/api/v1/suppliers", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	err := createSuppliers([]dto.Supplier{supplier}, postUrl)
	s.Require().NoError(err)

	sR := postgres.NewSupplierRepository(s.db, s.logger)
	takedData, err := sR.GetByName(context.Background(), supplier.Name)
	check := mapper.SupplierToDTO(*takedData)
	s.Require().NoError(err)
	s.Require().Equal(supplier, check)

	httpClient := http.Client{}
	patchUrl := fmt.Sprintf("http://%s:%s/api/v1/suppliers/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, takedData.Id)

	newAddress := &dto.Address{
		Country: "China",
		City:    "Beijing",
		Street:  "Wangfujing",
	}

	payload, err := json.Marshal(newAddress)
	s.Require().NoError(err)
	req, err := http.NewRequest(http.MethodPatch, patchUrl, bytes.NewBuffer(payload))
	s.Require().NoError(err)
	req.Header.Set("Content-type", "application/json")

	resp, err := httpClient.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	takedData, err = sR.GetByName(s.T().Context(), supplier.Name)
	s.Require().NoError(err)

	query := `SELECT COUNT(id) FROM address`
	var count int
	err = s.db.QueryRow(context.Background(), query).Scan(&count)
	s.Require().NoError(err)
	s.Require().Equal(1, count)

	getUrl := fmt.Sprintf("http://%s:%s/api/v1/suppliers/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, takedData.Id)
	getResp, err := http.Get(getUrl)
	s.Require().NoError(err)
	defer getResp.Body.Close()
	s.Require().Equal(http.StatusOK, getResp.StatusCode)

	raw, err := io.ReadAll(getResp.Body)
	s.Require().NoError(err)

	var controllCheck dto.Supplier
	err = json.Unmarshal(raw, &controllCheck)
	s.Require().NoError(err)
	validData := mapper.SupplierToDTO(*takedData)
	s.Require().Equal(validData, controllCheck)
}

func (s *TestSuite) TestUpdateAddressSupplierNotLinkedAddressBehavior() {
	s.CleanTable()
	first := dto.Supplier{
		Name:        "Servo Inc.",
		PhoneNumber: "+7 999 233 13 23",
		Address: &dto.Address{
			Country: "Russia",
			City:    "Moscow",
			Street:  "Tverskaya Street",
		},
	}

	second := dto.Supplier{
		Name:        "Servo Tech Inc.",
		PhoneNumber: "+231312312312",
		Address: &dto.Address{
			Country: "Japan",
			City:    "Tokyo",
			Street:  "Godzilla",
		},
	}

	suppliers := []dto.Supplier{first, second}

	postUrl := fmt.Sprintf("http://%s:%s/api/v1/suppliers", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	err := createSuppliers(suppliers, postUrl)
	s.Require().NoError(err)

	supRepo := postgres.NewSupplierRepository(s.db, s.logger)
	secondSupplier, err := supRepo.GetByName(context.Background(), second.Name)
	verifiable := mapper.SupplierToDTO(*secondSupplier)
	s.Require().NoError(err)
	s.Require().Equal(verifiable.Address, second.Address)

	newAddress := dto.Address{
		Country: "Russia",
		City:    "Moscow",
		Street:  "Tverskaya Street",
	}

	payload, err := json.Marshal(newAddress)
	s.Require().NoError(err)
	patchUrl := fmt.Sprintf("http://%s:%s/api/v1/suppliers/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, secondSupplier.Id)
	httpClient := http.Client{}
	req, err := http.NewRequest(http.MethodPatch, patchUrl, bytes.NewBuffer(payload))
	s.Require().NoError(err)
	req.Header.Set("Content-type", "application/json")

	patchResp, err := httpClient.Do(req)
	s.Require().NoError(err)

	s.Require().Equal(http.StatusOK, patchResp.StatusCode)

	query := `SELECT COUNT(id) FROM address`
	var count int
	err = s.db.QueryRow(context.Background(), query).Scan(&count)
	s.Require().NoError(err)
	s.Require().Equal(1, count)

	getUrl := fmt.Sprintf("http://%s:%s/api/v1/suppliers/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, secondSupplier.Id)
	getResp, err := http.Get(getUrl)
	s.Require().NoError(err)
	defer getResp.Body.Close()
	s.Require().Equal(http.StatusOK, getResp.StatusCode)

	raw, err := io.ReadAll(getResp.Body)
	s.Require().NoError(err)

	var check dto.Supplier
	err = json.Unmarshal(raw, &check)
	s.Require().NoError(err)
	s.Require().Equal(newAddress, *check.Address)
}

func (s *TestSuite) TestUpdateAddressSupplierFromSameAddress() {
	s.CleanTable()
	first := dto.Supplier{
		Name:        "Mech Inc.",
		PhoneNumber: "+7 999 233 13 23",
		Address: &dto.Address{
			Country: "Russia",
			City:    "Moscow",
			Street:  "Tverskaya Street",
		},
	}

	second := dto.Supplier{
		Name:        "Mech Tech Inc.",
		PhoneNumber: "+231312312312",
		Address: &dto.Address{
			Country: "Russia",
			City:    "Moscow",
			Street:  "Tverskaya Street",
		},
	}

	suppliers := []dto.Supplier{first, second}

	postUrl := fmt.Sprintf("http://%s:%s/api/v1/suppliers", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	err := createSuppliers(suppliers, postUrl)
	s.Require().NoError(err)

	supRepo := postgres.NewSupplierRepository(s.db, s.logger)
	takedData, err := supRepo.GetByName(context.Background(), second.Name)
	s.Require().NoError(err)
	checkTakedData := mapper.SupplierToDTO(*takedData)
	s.Require().Equal(checkTakedData.Address, second.Address)

	newAddress := dto.Address{
		Country: "Korea",
		City:    "Seoul",
		Street:  "Gangnam",
	}

	payload, err := json.Marshal(newAddress)
	s.Require().NoError(err)
	patchUrl := fmt.Sprintf("http://%s:%s/api/v1/suppliers/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, takedData.Id)
	httpClient := http.Client{}
	req, err := http.NewRequest(http.MethodPatch, patchUrl, bytes.NewBuffer(payload))
	s.Require().NoError(err)
	req.Header.Set("Content-type", "application/json")

	patchResp, err := httpClient.Do(req)
	s.Require().NoError(err)

	s.Require().Equal(http.StatusOK, patchResp.StatusCode)

	query := `SELECT COUNT(id) FROM address`
	var count int
	err = s.db.QueryRow(context.Background(), query).Scan(&count)
	s.Require().NoError(err)
	s.Require().Equal(2, count)

	getUrl := fmt.Sprintf("http://%s:%s/api/v1/suppliers/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, takedData.Id)
	getResp, err := http.Get(getUrl)
	s.Require().NoError(err)
	defer getResp.Body.Close()
	s.Require().Equal(http.StatusOK, getResp.StatusCode)

	raw, err := io.ReadAll(getResp.Body)
	s.Require().NoError(err)

	var check dto.Supplier
	err = json.Unmarshal(raw, &check)
	s.Require().NoError(err)
	s.Require().Equal(newAddress, *check.Address)
}

func (s *TestSuite) TestUpdateAddressSupplier2() {
	s.CleanTable()
	first := dto.Supplier{
		Name:        "Terra Inc.",
		PhoneNumber: "+7 999 233 13 23",
		Address: &dto.Address{
			Country: "Russia",
			City:    "Moscow",
			Street:  "Tverskaya Street",
		},
	}

	second := dto.Supplier{
		Name:        "Terra Tech Inc.",
		PhoneNumber: "+231312312312",
		Address: &dto.Address{
			Country: "German",
			City:    "Berlin",
			Street:  "Stag",
		},
	}

	third := dto.Supplier{
		Name:        "Dominion Inc.",
		PhoneNumber: "+231312312312",
		Address: &dto.Address{
			Country: "German",
			City:    "Berlin",
			Street:  "Stag",
		},
	}

	suppliers := []dto.Supplier{first, second, third}

	postUrl := fmt.Sprintf("http://%s:%s/api/v1/suppliers", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	err := createSuppliers(suppliers, postUrl)
	s.Require().NoError(err)

	supRepo := postgres.NewSupplierRepository(s.db, s.logger)
	firstCheck, err := supRepo.GetByName(context.Background(), third.Name)
	s.Require().NoError(err)
	firstCheckConv := mapper.SupplierToDTO(*firstCheck)
	s.Require().Equal(third.Address, firstCheckConv.Address)

	newAddress := dto.Address{
		Country: "Russia",
		City:    "Moscow",
		Street:  "Tverskaya Street",
	}

	payload, err := json.Marshal(newAddress)
	s.Require().NoError(err)
	patchUrl := fmt.Sprintf("http://%s:%s/api/v1/suppliers/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, firstCheck.Id)
	httpClient := http.Client{}
	req, err := http.NewRequest(http.MethodPatch, patchUrl, bytes.NewBuffer(payload))
	s.Require().NoError(err)
	req.Header.Set("Content-type", "application/json")

	patchResp, err := httpClient.Do(req)
	s.Require().NoError(err)

	s.Require().Equal(http.StatusOK, patchResp.StatusCode)

	query := `SELECT COUNT(id) FROM address`
	var count int
	err = s.db.QueryRow(context.Background(), query).Scan(&count)
	s.Require().NoError(err)
	s.Require().Equal(2, count)

	getUrl := fmt.Sprintf("http://%s:%s/api/v1/suppliers/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, firstCheck.Id)
	getResp, err := http.Get(getUrl)
	s.Require().NoError(err)
	defer getResp.Body.Close()
	s.Require().Equal(http.StatusOK, getResp.StatusCode)

	raw, err := io.ReadAll(getResp.Body)
	s.Require().NoError(err)

	var check dto.Supplier
	err = json.Unmarshal(raw, &check)
	s.Require().NoError(err)
	s.Require().Equal(newAddress, *check.Address)
}

func (s *TestSuite) TestDeleteSupplier() {
	s.CleanTable()

	first := dto.Supplier{
		Name:        "Terra Inc.",
		PhoneNumber: "+7 999 233 13 23",
		Address: &dto.Address{
			Country: "Russia",
			City:    "Moscow",
			Street:  "Tverskaya Street",
		},
	}

	second := dto.Supplier{
		Name:        "Terra Tech Inc.",
		PhoneNumber: "+231312312312",
		Address: &dto.Address{
			Country: "German",
			City:    "Berlin",
			Street:  "Stag",
		},
	}

	third := dto.Supplier{
		Name:        "Dominion Inc.",
		PhoneNumber: "+231312312312",
		Address: &dto.Address{
			Country: "German",
			City:    "Berlin",
			Street:  "Stag",
		},
	}

	postUrl := fmt.Sprintf("http://%s:%s/api/v1/suppliers", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	err := createSuppliers([]dto.Supplier{first, second, third}, postUrl)
	s.Require().NoError(err)

	getUrl := fmt.Sprintf("http://%s:%s/api/v1/suppliers", s.cfg.CrudService.Address, s.cfg.CrudService.Port)
	resp, err := http.Get(getUrl)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	temp, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)
	var firstCheck []dto.Supplier
	err = json.Unmarshal(temp, &firstCheck)
	s.Require().NoError(err)
	s.Require().Len(firstCheck, 3)

	supRepo := postgres.NewSupplierRepository(s.db, s.logger)
	tempSup, err := supRepo.GetByName(context.Background(), second.Name)
	s.Require().NoError(err)

	httpClient := http.Client{}
	deleteUrl := fmt.Sprintf("http://%s:%s/api/v1/suppliers/%s", s.cfg.CrudService.Address, s.cfg.CrudService.Port, tempSup.Id.String())
	req, err := http.NewRequest(http.MethodPatch, deleteUrl, nil)
	s.Require().NoError(err)

	deleteResp, err := httpClient.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNoContent, deleteResp.StatusCode)

	query := `SELECT COUNT(id) FROM address`
	var count int
	err = s.db.QueryRow(context.Background(), query).Scan(&count)
	s.Require().NoError(err)
	s.Require().Equal(2, count)

	query = `SELECT COUNT(id) FROM suppliers`
	count = 0
	err = s.db.QueryRow(context.Background(), query).Scan(&count)
	s.Require().NoError(err)
	s.Require().Equal(2, count)

	resp, err = http.Get(getUrl)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	temp, err = io.ReadAll(resp.Body)
	s.Require().NoError(err)
	var secondCheck []dto.Supplier
	err = json.Unmarshal(temp, &secondCheck)
	s.Require().NoError(err)
	s.Require().Len(secondCheck, 2)
	s.Require().NotContains(secondCheck, second)
}

func (s *TestSuite) TestDeleteSupplierEmptyTable() {

}
