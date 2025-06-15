package mapper

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
	"fmt"
	"time"
)

const dateFormat = "2006-01-02"

func ClientToDTO(client domain.Client) dto.Client {
	output := dto.Client{
		Name:     client.Name,
		Surname:  client.Surname,
		Birthday: client.Birthday.Format(dateFormat),
		Gender:   client.Gender,
	}

	if client.Address != nil {
		address := AddressToDto(*client.Address)
		output.Address = &address
	}

	return output
}

func ClientToDomain(dto dto.Client) (domain.Client, error) {
	dtoBirthday, err := time.Parse(dateFormat, dto.Birthday)
	if err != nil {
		return domain.Client{}, fmt.Errorf("clinet mapper: %v", err)
	}

	client := domain.Client{
		Name:     dto.Name,
		Surname:  dto.Surname,
		Birthday: dtoBirthday,
		Gender:   dto.Gender,
	}

	if dto.Address != nil {
		address := AddressToDomain(*dto.Address)
		client.Address = &address
	}

	return client, nil
}
