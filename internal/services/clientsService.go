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

var addressRepoName = uow.RepositoryName("address")
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
		repo, err := tx.Get(addressRepoName)
		if err != nil {
			s.logger.Debug("Address transaction problem on creating", logger.Err(err), "op", op)
			return err
		}

		repoGen := repo.(uow.RepositoryGenerator)(tx.GetTX(), s.logger)
		addressRepo := repoGen.(*postgres.AddressRepo)
		address_id, err := addressRepo.Create(ctx, client.Address)
		if err != nil {
			s.logger.Debug("Address transaction problem on creating", logger.Err(err), "op", op)
			return err
		}

		client.Address.Id = address_id
		repo, err = tx.Get(clientRepoName)
		if err != nil {
			s.logger.Debug("Client transaction problem on creating", logger.Err(err), "op", op)
			return err
		}

		repoGen = repo.(uow.RepositoryGenerator)(tx.GetTX(), s.logger)
		clientRepo := repoGen.(*postgres.ClientRepo)
		return clientRepo.Create(ctx, client)
	})

	if err != nil {
		s.logger.Debug("Somthing wrong with UOW creating", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unit of work creating problem: %v", op, err)
	}

	return nil
}

func (s *clientsService) GetAll(ctx context.Context, limit, offset int) ([]domain.Client, error) {
	op := "services.clientService.GetAll"

	if limit <= 0 || offset < 0 {
		s.logger.Debug("Invalid parameter limit and offset", "limit", limit, "offset", offset, "op", op)
		return nil, fmt.Errorf("%s: %w", op, crud_errors.ErrInvalidParam)
	}

	clients, err := s.reader.GetAll(ctx, limit, offset)
	if err != nil {
		s.logger.Debug("Get all unable", logger.Err(err), "op", op)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return clients, nil
}

func (s *clientsService) GetByNameAndSurname(ctx context.Context, name, surname string) ([]domain.Client, error) {
	op := "services.clientsService.GetByNameAndSurname"

	if name == "" || surname == "" {
		s.logger.Debug("Name or Surname is empty", "op", op)
		return nil, fmt.Errorf("%s: %w", op, crud_errors.ErrInvalidParam)
	}

	clients, err := s.reader.GetByNameAndSurname(ctx, name, surname)
	if err != nil {
		s.logger.Debug("Error recieved from repository", logger.Err(err), "op", op)
		return nil, fmt.Errorf("%s: Error when taking a client by first and last name: %w", op, err)
	}

	return clients, nil
}

func (s *clientsService) UpdateAddress(ctx context.Context, id uuid.UUID, address domain.Address) error {
	op := "services.clientsService.UpdateAddress"
	err := s.uow.Do(ctx, func(ctx context.Context, tx uow.Transaction) error {
		repo, err := tx.Get(addressRepoName)
		if err != nil {
			s.logger.Debug("Get address transaction problem on updating", logger.Err(err), "op", op)
			return err
		}

		repoGen := repo.(uow.RepositoryGenerator)(tx.GetTX(), s.logger)
		addressRepo := repoGen.(*postgres.AddressRepo)
		addressId, err := addressRepo.Create(ctx, address)
		if err != nil {
			return err
		}

		repo, err = tx.Get(clientRepoName)
		if err != nil {
			s.logger.Debug("Get transaction problem on updating", logger.Err(err), "op", op)
			return err
		}

		repoGen = repo.(uow.RepositoryGenerator)(tx.GetTX(), s.logger)
		clientRepo := repoGen.(*postgres.ClientRepo)
		return clientRepo.UpdateAddress(ctx, id, addressId)
	})

	if err != nil {
		s.logger.Debug("Somthin wrong with UOW updating", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unit of work update problem: %v", op, err)
	}

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
		client, err := clientRepo.GetById(ctx, id)
		if err != nil {
			s.logger.Debug("Not found or something went wrong", logger.Err(err), "op", op+".uow")
			return err
		}

		err = clientRepo.Delete(ctx, id)
		if err != nil {
			s.logger.Debug("something went wrong in delete", logger.Err(err), "op", op+".uow")
			return err
		}

		repo, err = tx.Get(addressRepoName)
		if err != nil {
			s.logger.Debug("get address repositroy is unable", logger.Err(err), "op", op+".uow")
			return err
		}

		repoGen = repo.(uow.RepositoryGenerator)(tx.GetTX(), s.logger)
		addressRepo := repoGen.(*postgres.AddressRepo)
		savepoint := `SAVEPOINT sq_delete_address;`
		_, err = tx.GetTX().Exec(ctx, savepoint)
		if err != nil {
			s.logger.Debug("unable to set savepoint before delete address", logger.Err(err), "op", op+".uow")
			return fmt.Errorf("%s.uow: unable to set savepoint: %v", op, err)
		}

		err = addressRepo.Delete(ctx, client.Address.Id)
		if err != nil {

			if errors.Is(err, crud_errors.ErrForeignKeyViolation) {
				backToSave := `ROLLBACK TO SAVEPOINT sq_delete_address;`
				_, err = tx.GetTX().Exec(ctx, backToSave)
				if err != nil {
					s.logger.Debug("unable back to savepoint after try delete address", logger.Err(err), "op", op+".uow")
					return fmt.Errorf("%s.uow: unable back to savepoint: %v", op, err)
				}

				return nil
			}

			s.logger.Debug("deleting is unable: something went wrong. can't rollback", logger.Err(err), "op", op)
			return nil
		}

		return nil
	})

	if err != nil {
		s.logger.Debug("Something wrong with UOW deleting", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unit of work delete problem: %v", op, err)
	}

	return nil
}
