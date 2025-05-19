package mapper

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/dto"
	"time"
)

const dateFormat = "2006-01-02"

func ClientToDTO(client *domain.Client) (dto.Client, error) {
	if client == nil {
		return dto.Client{}, ErrNoContent
	}

	return dto.Client{
		Name:      client.Name,
		Surname:   client.Surname,
		Birthday:  client.Birthday.Format(dateFormat),
		Gender:    client.Gender,
		AddressID: client.AddressId,
	}, nil
}

func ClientToDomain(dto *dto.Client) (domain.Client, error) {
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
