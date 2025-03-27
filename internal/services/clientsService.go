package services

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/dto"
	"CRUD-HOME-APPLIANCE-STORE/internal/logger"
	"CRUD-HOME-APPLIANCE-STORE/internal/mapper"
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
	"CRUD-HOME-APPLIANCE-STORE/internal/repositories"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClientsService interface {
	CrudServiceInterface[model.Client, dto.ClientDTO]
	// Create(ctx context.Context, dto dto.ClientDTO) (*model.Client, error)
	GetClientByNameAndSurname(ctx context.Context, name, surname string) ([]dto.ClientDTO, error)
	ChangeAddressParameter(ctx context.Context, id string, newAddressId string) error
}

type clientsServiceImpl struct {
	*CrudService[model.Client, dto.ClientDTO]
	repository repositories.ClientRepositoryInterface
}

func NewClientService(rep repositories.ClientRepositoryInterface, logger *logger.Logger) *clientsServiceImpl {
	crudService := NewCrudService(rep, mapper.ClientToDTO, mapper.ClientToModel, logger)
	return &clientsServiceImpl{
		CrudService: crudService,
		repository:  rep,
	}
}

// func (s *clientsServiceImpl) Create(ctx context.Context, dto dto.ClientDTO) (*model.Client, error) {
// 	client, err := mapper.ClientToModel(dto)
// 	if err != nil {
// 		return nil, fmt.Errorf("Client service: Error adding a client: %v", err)
// 	}

// 	client.RegistrationDate = time.Now()

// 	_, err = s.repository.Create(ctx, *client)
// 	if err != nil {
// 		return nil, fmt.Errorf("Client service: Error adding a client: %v", err)
// 	}

// 	return client, nil
// }

func (s *clientsServiceImpl) GetClientByNameAndSurname(ctx context.Context, name, surname string) ([]dto.ClientDTO, error) {
	if name == "" || surname == "" {
		return nil, fmt.Errorf("Service: client name and surname cannot be empty")
	}

	clients, err := s.repository.GetClientByNameAndSurname(ctx, name, surname)
	if clients == nil && err == nil {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("Service: Error when taking a client by first and last name: %v", err)
	}

	dtos := make([]dto.ClientDTO, len(clients))

	for i, item := range clients {
		dtos[i] = s.mapperD(item)
	}

	return dtos, nil
}

func (s *clientsServiceImpl) ChangeAddressParameter(ctx context.Context, id string, newAddressId string) error {
	if newAddressId == "" {
		return fmt.Errorf("Service: Invalid new address id")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("Service: Error change address parameter: %v", err)
	}

	objectAddressID, err := primitive.ObjectIDFromHex(newAddressId)
	if err != nil {
		return fmt.Errorf("Service: Error change address parameter: %v", err)
	}

	err = s.repository.Update(ctx, objectID, objectAddressID)
	if err != nil {
		return fmt.Errorf("Service: Error change address parameter: %v", err)
	}

	return nil
}
