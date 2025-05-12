package services

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/repositories/postgres"
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type RepositoryName string
type Repository any
type RepositoryGenerator func(tx *pgx.Tx, log *logger.Logger) Repository

// type ClientWriter interface {
// 	Create(ctx context.Context, client domain.Client) error
// 	UpdateAddress(ctx context.Context, id, address uuid.UUID) error
// 	Delete(ctx context.Context, id uuid.UUID) error
// }

type ClientReader interface {
	GetAll(ctx context.Context, limit, offset int) ([]domain.Client, error)
	GetByNameAndSurname(ctx context.Context, name, surname string) ([]domain.Client, error)
}

// type clientsService struct {
// 	// writer ClientWriter
// 	// reader ClientReader
// 	logger *logger.Logger
// }

// func NewClientService(writer ClientWriter, reader ClientReader, logger *logger.Logger) *clientsService {
// 	logger.Debug("Client service is created")
// 	return &clientsService{
// 		writer: writer,
// 		reader: reader,
// 		logger: logger,
// 	}
// }

type Transaction interface {
	Get(name RepositoryName) (Repository, error)
}

type UOW interface {
	Register(name RepositoryName, gen RepositoryGenerator) error
	Remove(name RepositoryName) error
	Clear()
	Do(ctx context.Context, fn func(ctx context.Context, tx Transaction) error) error
}

type clientsService struct {
	uow    UOW
	repo   ClientReader
	logger *logger.Logger
}

func NewClientService(reader ClientReader, unit UOW, logger *logger.Logger) *clientsService {
	return &clientsService{
		uow:    unit,
		repo:   reader,
		logger: logger,
	}
}

func (s *clientsService) Create(ctx context.Context, client domain.Client) error {
	op := "services.clientService.Create"

	err := s.uow.Do(ctx, func(ctx context.Context, tx Transaction) error {
		repo, err := tx.Get("client")
		if err != nil {
			s.logger.Debug("Get transaction problem on creating", logger.Err(err), "op", op)
			return err
		}

		userRepo := repo.(postgres.ClientRepo)
		return userRepo.Create(ctx, client)
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

	clients, err := s.repo.GetAll(ctx, limit, offset)
	if errors.Is(err, postgres.ErrClientNotFound) {
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

	clients, err := s.repo.GetByNameAndSurname(ctx, name, surname)
	if errors.Is(err, postgres.ErrClientNotFound) {
		s.logger.Debug("Client not found", "op", op)
		return nil, postgres.ErrClientNotFound
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
	err := s.uow.Do(ctx, func(ctx context.Context, tx Transaction) error {
		repo, err := tx.Get("client")
		if err != nil {
			s.logger.Debug("Get transaction problem on updating", logger.Err(err), "op", op)
			return err
		}

		userRepo := repo.(postgres.ClientRepo)
		return userRepo.UpdateAddress(ctx, id, address)
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

	err := s.uow.Do(ctx, func(ctx context.Context, tx Transaction) error {
		repo, err := tx.Get("client")
		if err != nil {
			s.logger.Debug("Get transaction problem", logger.Err(err), "op", op)
			return err
		}

		userRepo := repo.(postgres.ClientRepo)
		return userRepo.Delete(ctx, id)
	})

	if err != nil {
		s.logger.Debug("Somthin wrong with UOW deleting", logger.Err(err), "op", op)
		return fmt.Errorf("Client service: unit of work delete problem: %v", err)
	}

	s.logger.Debug("Client is deleted", "op", op)
	return nil
}
