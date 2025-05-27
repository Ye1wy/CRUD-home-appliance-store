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
		uowOp := op + ".uow"
		repo, err := tx.Get(addressRepoName)
		if err != nil {
			s.logger.Debug("get address repository generator is unable", logger.Err(err), "op", op)
			return fmt.Errorf("%s: get address repository generator is unable: %v", uowOp, err)
		}

		repoGen := repo.(uow.RepositoryGenerator)(tx.GetTX(), s.logger)
		addressRepo := repoGen.(*postgres.AddressRepo)
		addressId, err := addressRepo.Create(ctx, supplier.Address)
		if err != nil {
			s.logger.Debug("Address creation problem in supplier uow", logger.Err(err), "op", op)
			return fmt.Errorf("%s: unable to create address: %v", uowOp, err)
		}

		supplier.Address.Id = addressId
		repo, err = tx.Get(supplierRepoName)
		if err != nil {
			s.logger.Debug("get supplier repository generator is unable", logger.Err(err), "op", op)
			return fmt.Errorf("%s: get supplier repository generator is unable: %v", uowOp, err)
		}

		repoGen = repo.(uow.RepositoryGenerator)(tx.GetTX(), s.logger)
		supplierRepo := repoGen.(*postgres.SupplierRepo)
		if err := supplierRepo.Create(ctx, supplier); err != nil {
			s.logger.Debug("failed to create supplier", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: failed to create supplier: %v", uowOp, err)
		}

		return nil
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
		uowOp := op + ".uow"
		repo, err := tx.Get(addressRepoName)
		if err != nil {
			s.logger.Debug("get address repository generator is unable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: get address repository generator is unable: %v", uowOp, err)
		}

		repoGen := repo.(uow.RepositoryGenerator)(tx.GetTX(), s.logger)
		addressRepo := repoGen.(*postgres.AddressRepo)
		addressId, err := addressRepo.Create(ctx, address)
		if err != nil {
			s.logger.Debug("unable to create address", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: unable to create address: %v", uowOp, err)
		}

		address.Id = addressId
		repo, err = tx.Get(supplierRepoName)
		if err != nil {
			s.logger.Debug("get supplier repository generator is unable ", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: get supplier repository generator is unable: %v", uowOp, err)
		}

		repoGen = repo.(uow.RepositoryGenerator)(tx.GetTX(), s.logger)
		supplierRepo := repoGen.(*postgres.SupplierRepo)
		if err := supplierRepo.Update(ctx, id, address.Id); err != nil {
			s.logger.Debug("failed to update address with supplier", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: failed to update address with supplier: %v", uowOp, err)
		}

		return nil
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
		uowOp := op + ".uow"
		repo, err := tx.Get(supplierRepoName)
		if err != nil {
			s.logger.Debug("get supplier repository generator is unable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: get supplier repository generator is unable: %v", uowOp, err)
		}

		repoGen := repo.(uow.RepositoryGenerator)(tx.GetTX(), s.logger)
		supplierRepo := repoGen.(*postgres.SupplierRepo)
		suppler, err := supplierRepo.GetById(ctx, id)
		if err != nil {
			if errors.Is(err, crud_errors.ErrNotFound) {
				s.logger.Debug("supplier not found", "op", uowOp)
				return fmt.Errorf("%s: %w", uowOp, err)
			}

			s.logger.Debug("unable to get supplier data", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: unable to get supplier data: %v", uowOp, err)
		}

		if err := supplierRepo.Delete(ctx, id); err != nil {
			s.logger.Debug("unable to delete supplier", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: unable to delete supplier: %v", uowOp, err)
		}

		repo, err = tx.Get(addressRepoName)
		if err != nil {
			s.logger.Debug("get address repository generator is unable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: get address repository generator is unable: %v", uowOp, err)
		}

		repoGen = repo.(uow.RepositoryGenerator)(tx.GetTX(), s.logger)
		addressRepo := repoGen.(*postgres.AddressRepo)

		sqlStatement := `SAVEPOINT sq_delete_address;`
		_, err = tx.GetTX().Exec(ctx, sqlStatement)
		if err != nil {
			s.logger.Debug("unable to set savepoint before delete address", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: unable to set savepoint: %v", uowOp, err)
		}

		if err := addressRepo.Delete(ctx, suppler.Address.Id); err != nil {
			if errors.Is(err, crud_errors.ErrForeignKeyViolation) {
				backToSave := `ROLLBACK TO SAVEPOINT sq_delete_address;`
				_, err := tx.GetTX().Exec(ctx, backToSave)
				if err != nil {
					s.logger.Debug("unable back to savepoint after try delete address", logger.Err(err), "op", uowOp)
					return fmt.Errorf("%s: unable back to savepoint: %v", uowOp, err)
				}

				return nil
			}

			s.logger.Debug("deleting is unable: unexpected error from delete address: rollback to savepoint is unavailable", logger.Err(err), "op", uowOp)
			return fmt.Errorf("%s: unexpected error from delete address: %v", uowOp, err)
		}

		return nil
	})

	if err != nil {
		s.logger.Debug("Somethin wrong with UOW deleting", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unit of work delete problem: %v", op, err)
	}

	return nil
}
