package services

import (
	crud_errors "CRUD-HOME-APPLIANCE-STORE/internal/errors"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/repositories/postgres"
	"CRUD-HOME-APPLIANCE-STORE/internal/uow"
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

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
	logger.Debug("Client service is created")
	return &clientsService{
		uow:    unit,
		reader: reader,
		logger: logger,
	}
}

func (s *clientsService) Create(ctx context.Context, client *domain.Client) error {
	op := "services.clientService.Create"

	err := s.uow.Do(ctx, func(ctx context.Context, tx uow.Transaction) error {
		uowOp := op + ".uow"
		addressRepoGen, err := getReposiotry(tx, uow.AddressRepoName, s.logger)
		if err != nil {
			s.logger.Error("get address repository generator is unable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: get address repository generator is unable: %v", uowOp, err)
		}

		addressRepo := addressRepoGen.(*postgres.AddressRepo)
		err = addressRepo.Create(ctx, &client.Address)
		if err != nil {
			s.logger.Error("address creation is unavailable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: unable to create address: %v", uowOp, err)
		}

		clientRepoGen, err := getReposiotry(tx, uow.ClientRepoName, s.logger)
		if err != nil {
			s.logger.Error("get client repository generator is unable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: error when try to get repository generator: %v", uowOp, err)
		}

		clientRepo := clientRepoGen.(*postgres.ClientRepo)

		if err := clientRepo.Create(ctx, client); err != nil {
			s.logger.Error("failed to create client", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: failed to create client: %v", uowOp, err)
		}

		return nil
	})

	if err != nil {
		s.logger.Error("something wrong with UOW creating", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unit of work creating problem: %v", op, err)
	}

	return nil
}

func (s *clientsService) GetAll(ctx context.Context, limit, offset int) ([]domain.Client, error) {
	op := "services.clientService.GetAll"

	if limit <= 0 || offset < 0 {
		s.logger.Error("invalid parameter limit and offset", "limit", limit, "offset", offset, "op", op)
		return nil, fmt.Errorf("%s: %w", op, crud_errors.ErrInvalidParam)
	}

	clients, err := s.reader.GetAll(ctx, limit, offset)
	if err != nil {
		s.logger.Error("error detected", logger.Err(err), "op", op)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return clients, nil
}

func (s *clientsService) GetByNameAndSurname(ctx context.Context, name, surname string) ([]domain.Client, error) {
	op := "services.clientsService.GetByNameAndSurname"

	if name == "" || surname == "" {
		s.logger.Error("name or surname is empty", "op", op)
		return nil, fmt.Errorf("%s: %w", op, crud_errors.ErrInvalidParam)
	}

	clients, err := s.reader.GetByNameAndSurname(ctx, name, surname)
	if err != nil {
		s.logger.Error("error recieved from repository", logger.Err(err), "op", op)
		return nil, fmt.Errorf("%s: error when taking a client by first and last name: %w", op, err)
	}

	return clients, nil
}

func (s *clientsService) UpdateAddress(ctx context.Context, id uuid.UUID, address *domain.Address) error {
	op := "services.clientsService.UpdateAddress"
	err := s.uow.Do(ctx, func(ctx context.Context, tx uow.Transaction) error {
		uowOp := op + ".uow"
		addressRepoGen, err := getReposiotry(tx, uow.AddressRepoName, s.logger)
		if err != nil {
			s.logger.Error("get address repository generator is unable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: get address repository generator is unable: %v", uowOp, err)
		}

		addressRepo := addressRepoGen.(*postgres.AddressRepo)
		err = addressRepo.Create(ctx, address)
		if err != nil {
			s.logger.Error("address creation is unavailable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: unable to create address: %v", uowOp, err)
		}

		clientRepoGen, err := getReposiotry(tx, uow.ClientRepoName, s.logger)
		if err != nil {
			s.logger.Error("get client repository generator is unable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: error when try to get repository generator: %v", uowOp, err)
		}

		clientRepo := clientRepoGen.(*postgres.ClientRepo)
		if err := clientRepo.UpdateAddress(ctx, id, address.Id); err != nil {
			if errors.Is(err, crud_errors.ErrNotFound) {
				s.logger.Debug("update initialize is unable", logger.Err(err), "op", uowOp)
				return fmt.Errorf("%s: %w", uowOp, err)
			}

			s.logger.Error("failed to update address with client", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: failed to update address with client: %v", uowOp, err)
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, crud_errors.ErrNotFound) {
			s.logger.Warn("update initialize is unable: client not found", logger.Err(err), "op", op)
			return fmt.Errorf("%s: %w", op, err)
		}

		s.logger.Error("something wrong with UOW updating", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unit of work update problem: %w", op, err)
	}

	return nil
}

func (s *clientsService) Delete(ctx context.Context, id uuid.UUID) error {
	op := "services.clientService.Delete"
	err := s.uow.Do(ctx, func(ctx context.Context, tx uow.Transaction) error {
		uowOp := op + ".uow"
		clientRepoGen, err := getReposiotry(tx, uow.ClientRepoName, s.logger)
		if err != nil {
			s.logger.Error("get client repository generator is unable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: error when try to get repository generator: %v", uowOp, err)
		}

		clientRepo := clientRepoGen.(*postgres.ClientRepo)

		client, err := clientRepo.GetById(ctx, id)
		if err != nil {
			if errors.Is(err, crud_errors.ErrNotFound) {
				s.logger.Debug("client not found", "op", uowOp)
				return fmt.Errorf("%s: %w", op, err)
			}

			s.logger.Error("unable to get client data", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: unable to get client data: %v", uowOp, err)
		}

		if err := clientRepo.Delete(ctx, id); err != nil {
			s.logger.Error("unable to delete client", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: unable to delete client: %v", uowOp, err)
		}

		addressRepoGen, err := getReposiotry(tx, uow.AddressRepoName, s.logger)
		if err != nil {
			s.logger.Error("get address repository generator is unable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: get address repository generator is unable: %v", uowOp, err)
		}

		addressRepo := addressRepoGen.(*postgres.AddressRepo)

		savepoint := `sp_delete_address`
		err = safeDelete(ctx, tx.GetTX(), client.Address.Id, addressRepo.Delete, s.logger, uowOp, savepoint)
		if err != nil {
			if errors.Is(err, crud_errors.ErrNotFound) {
				s.logger.Debug("client not found", "op", uowOp)
				return nil
			}

			s.logger.Error("unable to safe delete address", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: unable to safe delete address: %v", uowOp, err)
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, crud_errors.ErrNotFound) {
			s.logger.Debug("client not found", "op", op)
			return fmt.Errorf("%s: %w", op, err)
		}

		s.logger.Error("something wrong with UOW deleting", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unit of work delete problem: %v", op, err)
	}

	return nil
}
