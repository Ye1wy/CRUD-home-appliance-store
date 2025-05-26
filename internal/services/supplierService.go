package services

import (
	crud_errors "CRUD-HOME-APPLIANCE-STORE/internal/errors"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/internal/repositories/postgres"
	"CRUD-HOME-APPLIANCE-STORE/internal/uow"
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"context"
	"fmt"

	"github.com/google/uuid"
)

var supplierRepoName = uow.RepositoryName("supplier")

type supplierReader interface {
	GetAll(ctx context.Context, limit, offset int) ([]domain.Supplier, error)
	GetById(ctx context.Context, id uuid.UUID) (*domain.Supplier, error)
}

type supplierService struct {
	uow    uow.UOW
	reader supplierReader
	logger *logger.Logger
}

func NewSupplierService(reader supplierReader, uow uow.UOW, logger *logger.Logger) *supplierService {
	logger.Debug("Supplier service is created")
	return &supplierService{
		uow:    uow,
		reader: reader,
		logger: logger,
	}
}

func (s *supplierService) Create(ctx context.Context, supplier domain.Supplier) error {
	op := "services.supplierService.Create"

	err := s.uow.Do(ctx, func(ctx context.Context, tx uow.Transaction) error {
		repo, err := tx.Get(addressRepoName)
		if err != nil {
			s.logger.Debug("Address get repo problem on create supplier transaction", logger.Err(err), "op", op)
			return err
		}

		repoGen := repo.(uow.RepositoryGenerator)(tx.GetTX(), s.logger)
		addressRepo := repoGen.(*postgres.AddressRepo)
		addressId, err := addressRepo.Create(ctx, supplier.Address)
		if err != nil {
			s.logger.Debug("Address creation problem in supplier uow", logger.Err(err), "op", op)
			return err
		}

		supplier.Address.Id = addressId
		repo, err = tx.Get(supplierRepoName)
		if err != nil {
			s.logger.Debug("Supplier get repo problem on create supplier transaction", logger.Err(err), "op", op)
			return err
		}

		repoGen = repo.(uow.RepositoryGenerator)(tx.GetTX(), s.logger)
		supplierRepo := repoGen.(*postgres.SupplierRepo)
		return supplierRepo.Create(ctx, supplier)
	})

	if err != nil {
		s.logger.Debug("Something wrong with UOW creating", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unit of work creating problem: %v", op, err)
	}

	return nil
}

func (s *supplierService) GetAll(ctx context.Context, limit, offset int) ([]domain.Supplier, error) {
	op := "services.supplierService.GetAll"

	if limit <= 0 || offset < 0 {
		s.logger.Debug("Invalid parameter limit and offset", "limit", limit, "offset", offset, "op", op)
		return nil, fmt.Errorf("%s: %w", op, crud_errors.ErrInvalidParam)
	}

	supplier, err := s.reader.GetAll(ctx, limit, offset)
	if err != nil {
		s.logger.Debug("Error detected", logger.Err(err), "op", op)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return supplier, nil
}

func (s *supplierService) GetById(ctx context.Context, id uuid.UUID) (*domain.Supplier, error) {
	op := "services.supplierService.GetById"
	supplier, err := s.reader.GetById(ctx, id)
	if err != nil {
		s.logger.Debug("Error detected", logger.Err(err), "op", op)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return supplier, nil
}

func (s *supplierService) UpdateAddress(ctx context.Context, id uuid.UUID, address domain.Address) error {
	op := "services.supplierService.UpdateAddress"
	err := s.uow.Do(ctx, func(ctx context.Context, tx uow.Transaction) error {
		repo, err := tx.Get(addressRepoName)
		if err != nil {
			s.logger.Debug("Address get repo problem on create supplier transaction", logger.Err(err), "op", op)
			return err
		}

		repoGen := repo.(uow.RepositoryGenerator)(tx.GetTX(), s.logger)
		addressRepo := repoGen.(*postgres.AddressRepo)
		addressId, err := addressRepo.Create(ctx, address)
		if err != nil {
			s.logger.Debug("Address creation problem in supplier uow", logger.Err(err), "op", op)
			return err
		}

		address.Id = addressId
		repo, err = tx.Get(supplierRepoName)
		if err != nil {
			s.logger.Debug("Get transaction problem on updating", logger.Err(err), "op", op)
			return err
		}

		repoGen = repo.(uow.RepositoryGenerator)(tx.GetTX(), s.logger)
		supplierRepo := repoGen.(*postgres.SupplierRepo)
		return supplierRepo.Update(ctx, id, address.Id)
	})

	if err != nil {
		s.logger.Debug("Something wrong with UOW creating", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unit of work creating problem: %v", op, err)
	}

	return nil
}

func (s *supplierService) Delete(ctx context.Context, id uuid.UUID) error {
	op := "services.supplierService.Delete"
	err := s.uow.Do(ctx, func(ctx context.Context, tx uow.Transaction) error {
		repo, err := tx.Get(supplierRepoName)
		if err != nil {
			s.logger.Debug("Get transaction problem on deleting", logger.Err(err), "op", op)
			return err
		}

		repoGen := repo.(uow.RepositoryGenerator)(tx.GetTX(), s.logger)
		supplierRepo := repoGen.(*postgres.SupplierRepo)
		return supplierRepo.Delete(ctx, id)
	})

	if err != nil {
		s.logger.Debug("Somethin wrong with UOW deleting", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unit of work delete problem: %v", op, err)
	}

	return nil
}
