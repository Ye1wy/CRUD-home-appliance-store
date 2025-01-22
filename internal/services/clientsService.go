package services

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/dto"
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
	"CRUD-HOME-APPLIANCE-STORE/internal/repositories"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClientsService interface {
	AddClient(ctx context.Context, dto *dto.ClientDTO) (*model.Client, error)
	GetAllClients(ctx context.Context, limit, offset int) ([]model.Client, error)
	GetClientByNameAndSurname(ctx context.Context, name, surname string) (*model.Client, error)
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
	client := &model.Client{
		ClientName:    dto.Name,
		ClientSurname: dto.Surname,
		Gender:        dto.Gender,
		AddressId:     dto.AddressID,
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

func (s *clientsServiceImpl) GetAllClients(ctx context.Context, limit, offset int) ([]model.Client, error) {
	if limit < 0 || offset < 0 {
		return nil, errors.New("limit and offset cannot be less of 0")
	}

	return s.Repo.GetAllClients(ctx, limit, offset)
}

func (s *clientsServiceImpl) GetClientByNameAndSurname(ctx context.Context, name, surname string) (*model.Client, error) {
	if name == "" || surname == "" {
		return nil, errors.New("client name and surname cannot be empty")
	}

	client, err := s.Repo.GetClientByNameAndSurname(ctx, name, surname)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return client, nil
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
