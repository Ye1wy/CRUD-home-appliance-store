package services

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/repositories/postgres"
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type ClientWriter interface {
	Create(ctx context.Context, client domain.Client) error
	UpdateAddress(ctx context.Context, id, address uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ClientReader interface {
	GetAll(ctx context.Context, limit, offset int) ([]domain.Client, error)
	GetByNameAndSurname(ctx context.Context, name, surname string) ([]domain.Client, error)
}

type clientsService struct {
	writer ClientWriter
	reader ClientReader
	logger *logger.Logger
}

func NewClientService(writer ClientWriter, reader ClientReader, logger *logger.Logger) *clientsService {
	logger.Debug("Client service is created")
	return &clientsService{
		writer: writer,
		reader: reader,
		logger: logger,
	}
}

func (s *clientsService) Create(ctx context.Context, client domain.Client) error {
	op := "services.clientService.Create"

	// if err := s.writer.UnitOfWork(ctx, func(tx psgrep.WriteClientRepo) error {
	// 	if err := s.writer.Create(ctx, client); err != nil {
	// 		s.logger.Debug("Failed create client", logger.Err(err), "op", op)
	// 		return fmt.Errorf("Client Service: failed creating client: %v", err)
	// 	}

	// 	return nil

	// }); err != nil {
	// 	return fmt.Errorf("Client Service: unit of work problem: %v", err)
	// }

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

	clients, err := s.reader.GetByNameAndSurname(ctx, name, surname)
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

	// if err := s.writer.UnitOfWork(ctx, func(psgrep.WriteClientRepo) error {
	// 	if err := s.writer.UpdateAddress(ctx, id, address); err != nil {
	// 		s.logger.Debug("Error recieved from Update", logger.Err(err), "op", op)
	// 		return fmt.Errorf("Client Service: change address is uable: %v", err)
	// 	}

	// 	return nil
	// }); err != nil {
	// 	s.logger.Debug("Unit Of Work problem", logger.Err(err), "op", op)
	// 	return fmt.Errorf("Client Service: unit of work problem: %v", err)
	// }

	s.logger.Debug("Data updated", "op", op)
	return nil
}

func (s *clientsService) Delete(ctx context.Context, id uuid.UUID) error {
	// op := "services.clientService.Delete"

	// if err := s.writer.UnitOfWork(ctx, func(w psgrep.WriteClientRepo) error {
	// 	if err := s.writer.Delete(ctx, id); err != nil {
	// 		s.logger.Debug("delete is unable", logger.Err(err), "op", op)
	// 		return fmt.Errorf("Client Service: Delete error from repository: %v", err)
	// 	}

	// 	return nil
	// }); err != nil {
	// 	s.logger.Debug("Unit Of Work problem", logger.Err(err), "op", op)
	// 	return fmt.Errorf("Client Service: unit of work problem: %v", err)
	// }

	return nil
}
