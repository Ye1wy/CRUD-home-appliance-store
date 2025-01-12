package services

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
	"CRUD-HOME-APPLIANCE-STORE/internal/repositories"
	"context"
	"errors"
)

type ClientsService interface {
	AddClient(ctx context.Context, client *model.Client) error
	GetAllClients(ctx context.Context, limit, offset int) ([]model.Client, error)
	GetClientById(ctx context.Context, id int) (*model.Client, error)
	ChangeAddressParameter(ctx context.Context, id int, newAddressId int) error
	DeleteClientById(ctx context.Context, id int) error
}

type ClientsServiceImpl struct {
	Repo *repositories.ClientRepository
}

func NewClientService(rep *repositories.ClientRepository) *ClientsServiceImpl {
	return &ClientsServiceImpl{
		Repo: rep,
	}
}

func (s *ClientsServiceImpl) AddClient(ctx context.Context, client *model.Client) error {
	if client.ClientName == "" || client.ClientSurname == "" {
		return errors.New("client name and surname cannot be empty")
	}

	return s.Repo.AddClient(ctx, client)
}

func (s *ClientsServiceImpl) GetAllClients(ctx context.Context, limit, offset int) ([]model.Client, error) {
	if limit < 0 || offset < 0 {
		return nil, errors.New("limit and offset cannot be less of 0")
	}

	return s.Repo.GetAllClients(ctx, limit, offset)
}

func (s *ClientsServiceImpl) GetClientById(ctx context.Context, name, surname string) (*model.Client, error) {
	if name == "" || surname == "" {
		return nil, errors.New("client name and surname cannot be empty")
	}

	client, err := s.Repo.GetClientByNameAndSurname(ctx, name, surname)
	if err != nil {
		return nil, err
	}

	if client == nil {
		return nil, errors.New("client not found")
	}

	return client, nil
}

func (s *ClientsServiceImpl) ChangeAddressParameter(ctx context.Context, id int, newAddressId int) error {
	if newAddressId < 0 {
		return errors.New("invalid new address id")
	}

	err := s.Repo.UpdateAddress(ctx, id, newAddressId)
	if err != nil {
		return err
	}

	return nil
}

func (s *ClientsServiceImpl) DeleteClientById(ctx context.Context, id int) error {
	return s.Repo.DeleteClientById(ctx, id)
}
