package mapper

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
	"fmt"
	"time"
)

const dateFormat = "2006-01-02"

func ClientToDTO(client domain.Client) dto.Client {
	return dto.Client{
		Name:     client.Name,
		Surname:  client.Surname,
		Birthday: client.Birthday.Format(dateFormat),
		Gender:   client.Gender,
		Address: dto.Address{
			Country: client.Address.Country,
			City:    client.Address.City,
			Street:  client.Address.Street,
		},
	}
}

func ClientToDomain(dto dto.Client) (domain.Client, error) {
	dtoBirthday, err := time.Parse(dateFormat, dto.Birthday)
	if err != nil {
		return domain.Client{}, fmt.Errorf("clinet mapper: %v", err)
	}

	return domain.Client{
		Name:     dto.Name,
		Surname:  dto.Surname,
		Birthday: dtoBirthday,
		Gender:   dto.Gender,
		Address: domain.Address{
			Country: dto.Address.Country,
			City:    dto.Address.City,
			Street:  dto.Address.Street,
		},
	}, nil
}
