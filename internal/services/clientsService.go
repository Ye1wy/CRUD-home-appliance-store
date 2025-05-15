package services

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/repositories/postgres"
	"CRUD-HOME-APPLIANCE-STORE/internal/uow"
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var clientRepoName = uow.RepositoryName("client")

type ClientReader interface {
	GetAll(ctx context.Context, limit, offset int) ([]domain.Client, error)
	GetByNameAndSurname(ctx context.Context, name, surname string) ([]domain.Client, error)
}

type clientsService struct {
	uow    uow.UOW
	reader ClientReader
	logger *logger.Logger
}

func NewClientService(reader ClientReader, unit uow.UOW, logger *logger.Logger) *clientsService {
	return &clientsService{
		uow:    unit,
		reader: reader,
		logger: logger,
	}
}

func (s *clientsService) Create(ctx context.Context, client domain.Client) error {
	op := "services.clientService.Create"

	err := s.uow.Do(ctx, func(ctx context.Context, tx uow.Transaction) error {
		repo, err := tx.Get(clientRepoName)
		if err != nil {
			s.logger.Debug("Client transaction problem on creating", logger.Err(err), "op", op)
			return err
		}

		repoGen := repo.(uow.RepositoryGenerator)(tx.GetTX(), s.logger)
		clientRepo := repoGen.(*postgres.ClientRepo)
		return clientRepo.Create(ctx, client)
	})

	if err != nil {
		s.logger.Debug("Somthing wrong with UOW creating", logger.Err(err), "op", op)
		return fmt.Errorf("Client service: unit of work creating problem: %v", err)
	}

	s.logger.Debug("Client is created", "op", op)
	return nil
}

func (s *clientsService) GetAll(ctx context.Context, limit, offset int) ([]domain.Client, error) {
	op := "services.clientService.GetAll"

	if limit <= 0 || offset < 0 {
		s.logger.Debug("Invalid parameter limit and offset", "limit", limit, "offset", offset, "op", op)
		return nil, ErrInvalidParam
	}

	clients, err := s.reader.GetAll(ctx, limit, offset)
	if errors.Is(err, postgres.ErrNotFound) {
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

	clients, err := s.reader.GetByNameAndSurname(ctx, name, surname)
	if errors.Is(err, postgres.ErrNotFound) {
		s.logger.Debug("Client not found", "op", op)
		return nil, postgres.ErrNotFound
	}

	if err != nil {
		s.logger.Debug("Error recieved from repository", logger.Err(err), "op", op)
		return nil, fmt.Errorf("Client Service: Error when taking a client by first and last name: %v", err)
	}

	s.logger.Debug("All clients is retrieved", "op", op)
	return clients, nil
}

func (s *clientsService) UpdateAddress(ctx context.Context, id, address uuid.UUID) error {
	op := "services.clientsService.UpdateAddress"
	err := s.uow.Do(ctx, func(ctx context.Context, tx uow.Transaction) error {
		repo, err := tx.Get(clientRepoName)
		if err != nil {
			s.logger.Debug("Get transaction problem on updating", logger.Err(err), "op", op)
			return err
		}

		repoGen := repo.(uow.RepositoryGenerator)(tx.GetTX(), s.logger)
		clientRepo := repoGen.(*postgres.ClientRepo)
		return clientRepo.UpdateAddress(ctx, id, address)
	})

	if err != nil {
		s.logger.Debug("Somthin wrong with UOW updating", logger.Err(err), "op", op)
		return fmt.Errorf("Client service: unit of work update problem: %v", err)
	}

	s.logger.Debug("Data updated", "op", op)
	return nil
}

func (s *clientsService) Delete(ctx context.Context, id uuid.UUID) error {
	op := "services.clientService.Delete"

	err := s.uow.Do(ctx, func(ctx context.Context, tx uow.Transaction) error {
		repo, err := tx.Get(clientRepoName)
		if err != nil {
			s.logger.Debug("Get transaction problem", logger.Err(err), "op", op)
			return err
		}

		repoGen := repo.(uow.RepositoryGenerator)(tx.GetTX(), s.logger)
		clientRepo := repoGen.(*postgres.ClientRepo)
		return clientRepo.Delete(ctx, id)
	})

	if err != nil {
		s.logger.Debug("Somthin wrong with UOW deleting", logger.Err(err), "op", op)
		return fmt.Errorf("Client service: unit of work delete problem: %v", err)
	}

	s.logger.Debug("Client is deleted", "op", op)
	return nil
}
