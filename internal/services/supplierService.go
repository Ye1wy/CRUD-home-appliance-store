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

type supplierReader interface {
	GetAll(ctx context.Context, limit, offset int) ([]domain.Supplier, error)
	GetById(ctx context.Context, id uuid.UUID) (*domain.Supplier, error)
}

type supplierService struct {
	uow    uow.UOW
	reader supplierReader
	logger *logger.Logger
}

func NewSupplierService(reader supplierReader, unit uow.UOW, logger *logger.Logger) *supplierService {
	logger.Debug("Supplier service is created")
	return &supplierService{
		uow:    unit,
		reader: reader,
		logger: logger,
	}
}

func (s *supplierService) Create(ctx context.Context, supplier *domain.Supplier) error {
	op := "services.supplierService.Create"

	err := s.uow.Do(ctx, func(ctx context.Context, tx uow.Transaction) error {
		uowOp := op + ".uow"
		addressRepoGen, err := getReposiotry(tx, uow.AddressRepoName, s.logger)
		if err != nil {
			s.logger.Debug("get address repository generator is unable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: get address repository generator is unable: %v", uowOp, err)
		}

		addressRepo := addressRepoGen.(*postgres.AddressRepo)
		err = addressRepo.Create(ctx, &supplier.Address)
		if err != nil {
			s.logger.Debug("address creation is unavailable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: unable to create address: %v", uowOp, err)
		}

		supplierRepoGen, err := getReposiotry(tx, uow.SupplierRepoName, s.logger)
		if err != nil {
			s.logger.Debug("get supplier repository generator is unable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: get supplier repository generator is unable: %v", uowOp, err)
		}

		supplierRepo := supplierRepoGen.(*postgres.SupplierRepo)
		if err := supplierRepo.Create(ctx, supplier); err != nil {
			s.logger.Debug("failed to create supplier", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: failed to create supplier: %v", uowOp, err)
		}

		return nil
	})

	if err != nil {
		s.logger.Debug("something wrong with UOW creating", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unit of work creating problem: %v", op, err)
	}

	return nil
}

func (s *supplierService) GetAll(ctx context.Context, limit, offset int) ([]domain.Supplier, error) {
	op := "services.supplierService.GetAll"

	if limit <= 0 || offset < 0 {
		s.logger.Debug("invalid parameter limit and offset", "limit", limit, "offset", offset, "op", op)
		return nil, fmt.Errorf("%s: %w", op, crud_errors.ErrInvalidParam)
	}

	supplier, err := s.reader.GetAll(ctx, limit, offset)
	if err != nil {
		s.logger.Debug("error detected", logger.Err(err), "op", op)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return supplier, nil
}

func (s *supplierService) GetById(ctx context.Context, id uuid.UUID) (*domain.Supplier, error) {
	op := "services.supplierService.GetById"
	supplier, err := s.reader.GetById(ctx, id)
	if err != nil {
		s.logger.Debug("error detected", logger.Err(err), "op", op)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return supplier, nil
}

func (s *supplierService) UpdateAddress(ctx context.Context, id uuid.UUID, address *domain.Address) error {
	op := "services.supplierService.UpdateAddress"
	err := s.uow.Do(ctx, func(ctx context.Context, tx uow.Transaction) error {
		uowOp := op + ".uow"
		addressRepoGen, err := getReposiotry(tx, uow.AddressRepoName, s.logger)
		if err != nil {
			s.logger.Debug("get address repository generator is unable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: get address repository generator is unable: %v", uowOp, err)
		}

		addressRepo := addressRepoGen.(*postgres.AddressRepo)
		err = addressRepo.Create(ctx, address)
		if err != nil {
			s.logger.Debug("unable to create address", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: unable to create address: %v", uowOp, err)
		}

		supplierRepoGen, err := getReposiotry(tx, uow.SupplierRepoName, s.logger)
		if err != nil {
			s.logger.Debug("get supplier repository generator is unable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: get supplier repository generator is unable: %v", uowOp, err)
		}

		supplierRepo := supplierRepoGen.(*postgres.SupplierRepo)

		if err := supplierRepo.Update(ctx, id, address.Id); err != nil {
			s.logger.Debug("failed to update address with supplier", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: failed to update address with supplier: %v", uowOp, err)
		}

		return nil
	})

	if err != nil {
		s.logger.Debug("something wrong with UOW creating", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unit of work creating problem: %v", op, err)
	}

	return nil
}

func (s *supplierService) Delete(ctx context.Context, id uuid.UUID) error {
	op := "services.supplierService.Delete"
	err := s.uow.Do(ctx, func(ctx context.Context, tx uow.Transaction) error {
		uowOp := op + ".uow"
		supplierRepoGen, err := getReposiotry(tx, uow.SupplierRepoName, s.logger)
		if err != nil {
			s.logger.Debug("get supplier repository generator is unable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: get supplier repository generator is unable: %v", uowOp, err)
		}

		supplierRepo := supplierRepoGen.(*postgres.SupplierRepo)
		suppler, err := supplierRepo.GetById(ctx, id)
		if err != nil {
			if errors.Is(err, crud_errors.ErrNotFound) {
				s.logger.Debug("supplier not found", "op", uowOp)
				return nil
			}

			s.logger.Debug("unable to get supplier data", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: unable to get supplier data: %v", uowOp, err)
		}

		if err := supplierRepo.Delete(ctx, id); err != nil {
			s.logger.Debug("unable to delete supplier", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: unable to delete supplier: %v", uowOp, err)
		}

		addressRepoGen, err := getReposiotry(tx, uow.SupplierRepoName, s.logger)
		if err != nil {
			s.logger.Debug("get address repository generator is unable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: get address repository generator is unable: %v", uowOp, err)
		}

		addressRepo := addressRepoGen.(*postgres.AddressRepo)

		savepoint := `sp_delete_address`
		err = safeDelete(ctx, tx.GetTX(), suppler.Address.Id, addressRepo.Delete, s.logger, uowOp, savepoint)
		if err != nil {
			s.logger.Debug("unable to safe delete address", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: unable to safe delete address: %v", uowOp, err)
		}

		return nil
	})

	if err != nil {
		s.logger.Debug("somethin wrong with UOW deleting", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unit of work delete problem: %v", op, err)
	}

	return nil
}
