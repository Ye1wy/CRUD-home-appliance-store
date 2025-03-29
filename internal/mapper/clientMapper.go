package mapper

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/dto"
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
	"time"
)

const dateFormat = "2006-01-02"

func ClientToDTO(client *model.Client) (*dto.ClientDTO, error) {
	if client == nil {
		return nil, nil
	}

	return &dto.ClientDTO{
		Name:      client.ClientName,
		Surname:   client.ClientSurname,
		Birthday:  client.Birthday.Format(dateFormat),
		Gender:    client.Gender,
		AddressID: client.AddressId,
	}, nil
}

func ClientToModel(dto *dto.ClientDTO) (*model.Client, error) {
	if dto == nil {
		return nil, nil
	}

	dtoBirthday, err := time.Parse(dateFormat, dto.Birthday)
	if err != nil {
		return nil, err
	}

	return &model.Client{
		ClientName:    dto.Name,
		ClientSurname: dto.Surname,
		Birthday:      dtoBirthday,
		Gender:        dto.Gender,
		AddressId:     dto.AddressID,
	}, nil
}

//// func ToClientDTOs(clients []model.Client) []dto.ClientDTO {
//// 	dtos := make([]dto.ClientDTO, len(clients))

//// 	for i, client := range clients {
//// 		dtos[i] = *ToClientDTO(&client)
//// }

//// 	return dtos
//// }

//// func ToClientModels(dtos []dto.ClientDTO) ([]model.Client, error) {
//// 	models := make([]model.Client, len(dtos))

//// 	for i, dto := range dtos {
//// 		client, err := ToClientModel(&dto)
// // 		if err != nil {
// // 			return nil, err
// // 		}

// // 		models[i] = *client
// // 	}

// // 	return models, nil
// // }
