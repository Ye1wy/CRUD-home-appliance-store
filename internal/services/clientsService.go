package services

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/dto"
	"CRUD-HOME-APPLIANCE-STORE/internal/mapper"
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
	"CRUD-HOME-APPLIANCE-STORE/internal/repositories"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClientsService interface {
	AddClient(ctx context.Context, dto *dto.ClientDTO) (*model.Client, error)
	GetAllClients(ctx context.Context, limit, offset int) ([]dto.ClientDTO, error)
	GetClientByNameAndSurname(ctx context.Context, name, surname string) ([]dto.ClientDTO, error)
	ChangeAddressParameter(ctx context.Context, id primitive.ObjectID, newAddressId string) error
	DeleteClientById(ctx context.Context, id primitive.ObjectID) error
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
		return nil, err
	}

	result, err := s.Repo.AddClient(ctx, client)
	if err != nil {
		return nil, err
	}

	if objectID, ok := result.InsertedID.(primitive.ObjectID); ok {
		client.Id = objectID.Hex()

	} else {
		return nil, fmt.Errorf("failed to parse inserted ID")
	}

	return client, nil
}

func (s *clientsServiceImpl) GetAllClients(ctx context.Context, limit, offset int) ([]dto.ClientDTO, error) {
	if limit < 0 || offset < 0 {
		return nil, errors.New("limit and offset cannot be less of 0")
	}

	clients, err := s.Repo.GetAllClients(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	dto := mapper.ToClientDTOs(clients)

	return dto, nil
}

func (s *clientsServiceImpl) GetClientByNameAndSurname(ctx context.Context, name, surname string) ([]dto.ClientDTO, error) {
	if name == "" || surname == "" {
		return nil, errors.New("client name and surname cannot be empty")
	}

	clients, err := s.Repo.GetClientByNameAndSurname(ctx, name, surname)
	if clients == nil && err == nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	dtos := mapper.ToClientDTOs(clients)

	return dtos, nil
}

func (s *clientsServiceImpl) ChangeAddressParameter(ctx context.Context, id primitive.ObjectID, newAddressId string) error {
	if newAddressId == "" {
		return errors.New("invalid new address id")
	}

	err := s.Repo.UpdateAddress(ctx, id, newAddressId)
	if err != nil {
		return err
	}

	return nil
}

func (s *clientsServiceImpl) DeleteClientById(ctx context.Context, id primitive.ObjectID) error {
	return s.Repo.DeleteClientById(ctx, id)
}
