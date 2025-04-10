package services

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/logger"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	psgrep "CRUD-HOME-APPLIANCE-STORE/internal/repositories/postgres"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrInvalidParam = errors.New("invalid parameter")
)

type ClientRepositoryInterface interface {
	Create(ctx context.Context, client *domain.Client) error
	GetAll(ctx context.Context, limit, offset int) ([]domain.Client, error)
	GetByNameAndSurname(ctx context.Context, client domain.Client) ([]domain.Client, error)
	UpdateAddress(ctx context.Context, client *domain.Client) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type clientsService struct {
	repo   ClientRepositoryInterface
	logger *logger.Logger
}

func NewClientService(rep ClientRepositoryInterface, logger *logger.Logger) *clientsService {
	logger.Debug("Client service is created")
	return &clientsService{
		repo:   rep,
		logger: logger,
	}
}

func (s *clientsService) Create(ctx context.Context, client *domain.Client) error {
	op := "services.clientService.Create"

	if err := s.repo.Create(ctx, client); err != nil {
		s.logger.Debug("Failed create client", logger.Err(err), "op", op)
		return fmt.Errorf("Client Service: failed creating client: %v", err)
	}

	s.logger.Debug("Client is created", "op", op)
	return nil
}

func (s *clientsService) GetAll(ctx context.Context, limit, offer int) ([]domain.Client, error) {
	op := "services.clientService.GetAll"

	if limit <= 0 || offer <= 0 {
		return nil, ErrInvalidParam
	}

	clients, err := s.repo.GetAll(ctx, limit, offer)
	if errors.Is(err, psgrep.ErrClientNotFound) {
		s.logger.Debug("Clients not found", logger.Err(err), "op", op)
		return nil, err
	}

	if err != nil {
		s.logger.Debug("Get all unable", logger.Err(err), "op", op)
		return nil, err
	}

	return clients, nil
}

func (s *clientsService) GetByNameAndSurname(ctx context.Context, name, surname string) ([]domain.Client, error) {
	op := "services.clientsService.GetByNameAndSurname"

	if name == "" || surname == "" {
		s.logger.Debug("Name or Surname is empty", "op", op)
		return nil, ErrInvalidParam
	}

	model := domain.Client{
		Name:    name,
		Surname: surname,
	}

	clients, err := s.repo.GetByNameAndSurname(ctx, model)
	if errors.Is(err, psgrep.ErrClientNotFound) {
		s.logger.Debug("Client not found", "op", op)
		return nil, psgrep.ErrClientNotFound
	}

	if err != nil {
		s.logger.Debug("Error recieved from repository", logger.Err(err), "op", op)
		return nil, fmt.Errorf("Client Service: Error when taking a client by first and last name: %v", err)
	}

	s.logger.Debug("All clients is retrieved", "op", op)
	return clients, nil
}

func (s *clientsService) UpdateAddress(ctx context.Context, object *domain.Client) error {
	op := "services.clientsService.UpdateAddress"

	if err := s.repo.UpdateAddress(ctx, object); err != nil {
		s.logger.Debug("Error recieved from Update", logger.Err(err), "op", op)
		return fmt.Errorf("Client Service: change address is uable: %v", err)
	}

	s.logger.Debug("Data updated", "op", op)
	return nil
}

func (s *clientsService) Delete(ctx context.Context, id uuid.UUID) error {
	op := "services.clientService.Delete"

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Debug("delete is unable", logger.Err(err), "op", op)
		return fmt.Errorf("Client Service: Delete error from repository: %v", err)
	}

	return nil
}
