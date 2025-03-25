package services

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/dto"
	"CRUD-HOME-APPLIANCE-STORE/internal/mapper"
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
	"CRUD-HOME-APPLIANCE-STORE/internal/repositories"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClientsService interface {
	AddClient(ctx context.Context, dto *dto.ClientDTO) (*model.Client, error)
	GetAllClients(ctx context.Context, limit, offset int) ([]dto.ClientDTO, error)
	GetClientByNameAndSurname(ctx context.Context, name, surname string) ([]dto.ClientDTO, error)
	ChangeAddressParameter(ctx context.Context, id string, newAddressId string) error
	DeleteClientById(ctx context.Context, id string) error
}

type clientsServiceImpl struct {
	Repo repositories.ClientRepository
}

func NewClientService(rep repositories.ClientRepository) *clientsServiceImpl {
	return &clientsServiceImpl{
		Repo: rep,
	}
}

func (s *clientsServiceImpl) AddClient(ctx context.Context, dto *dto.ClientDTO) (*model.Client, error) {
	client, err := mapper.ToClientModel(dto)
	if err != nil {
		return nil, fmt.Errorf("Client service: Error adding a client: %v", err)
	}

	client.RegistrationDate = time.Now()

	_, err = s.Repo.AddClient(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("Client service: Error adding a client: %v", err)
	}

	// if objectID, ok := result.InsertedID.(primitive.ObjectID); ok {
	// 	client.Id = objectID.Hex()

	// } else {
	// 	return nil, fmt.Errorf("failed to parse inserted ID")
	// }

	return client, nil
}

func (s *clientsServiceImpl) GetAllClients(ctx context.Context, limit, offset int) ([]dto.ClientDTO, error) {
	if limit < 0 || offset < 0 {
		return nil, fmt.Errorf("Client service: limit and offset cannot be less of 0")
	}

	clients, err := s.Repo.GetAllClients(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("Client service: Error receriving all the client: %v", err)
	}

	dto := mapper.ToClientDTOs(clients)

	return dto, nil
}

func (s *clientsServiceImpl) GetClientByNameAndSurname(ctx context.Context, name, surname string) ([]dto.ClientDTO, error) {
	if name == "" || surname == "" {
		return nil, fmt.Errorf("Service: client name and surname cannot be empty")
	}

	clients, err := s.Repo.GetClientByNameAndSurname(ctx, name, surname)
	if clients == nil && err == nil {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("Service: Error when taking a client by first and last name: %v", err)
	}

	dtos := mapper.ToClientDTOs(clients)

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

	err = s.Repo.UpdateAddress(ctx, objectID, objectAddressID)
	if err != nil {
		return fmt.Errorf("Service: Error change address parameter: %v", err)
	}

	return nil
}

func (s *clientsServiceImpl) DeleteClientById(ctx context.Context, id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("Service: Error delete client by id: %v", err)
	}

	return s.Repo.DeleteClientById(ctx, objectId)
}
