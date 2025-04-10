package mapper

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
	"time"
)

const dateFormat = "2006-01-02"

func ClientToDTO(client *domain.Client) (dto.ClientDTO, error) {
	if client == nil {
		return dto.ClientDTO{}, ErrNoContent
	}

	return dto.ClientDTO{
		Name:      client.Name,
		Surname:   client.Surname,
		Birthday:  client.Birthday.Format(dateFormat),
		Gender:    client.Gender,
		AddressID: client.AddressId,
	}, nil
}

func ClientToDomain(dto *dto.ClientDTO) (domain.Client, error) {
	if dto == nil {
		return domain.Client{}, ErrNoContent
	}

	dtoBirthday, err := time.Parse(dateFormat, dto.Birthday)
	if err != nil {
		return domain.Client{}, err
	}

	return domain.Client{
		Name:      dto.Name,
		Surname:   dto.Surname,
		Birthday:  dtoBirthday,
		Gender:    dto.Gender,
		AddressId: dto.AddressID,
	}, nil
}
