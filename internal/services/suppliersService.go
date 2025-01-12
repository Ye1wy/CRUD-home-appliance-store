package services

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
	"CRUD-HOME-APPLIANCE-STORE/internal/repositories"
	"context"
	"errors"
)

type SuppliersService interface {
	AddSupplier(ctx context.Context, supplier *model.Supplier) error
	GetAllSuppliers(ctx context.Context, limit, offset int) ([]model.Supplier, error)
	GetSupplierById(ctx context.Context, id int) (*model.Supplier, error)
	ChangeAddressParameter(ctx context.Context, id int, newAddressId int) error
	DeleteSupplierById(ctx context.Context, id int) error
}

type SuppliersServiceImpl struct {
	Repo *repositories.SuppliersRepository
}

func NewSuppliersServiceImpl(repo *repositories.SuppliersRepository) *SuppliersServiceImpl {
	return &SuppliersServiceImpl{
		Repo: repo,
	}
}

func (s *SuppliersServiceImpl) AddSupplier(ctx context.Context, supplier *model.Supplier) error {
	if supplier.Name == "" || supplier.PhoneNumber == "" {
		return errors.New("supplier name and surname cannot be empty")
	}

	return s.Repo.AddSupplier(ctx, supplier)
}

func (s *SuppliersServiceImpl) GetAllSuppliers(ctx context.Context, limit, offset int) ([]model.Supplier, error) {
	if limit < 0 || offset < 0 {
		return nil, errors.New("limit and offset cannot be less of 0")
	}

	return s.Repo.GetAllSuppliers(ctx, limit, offset)
}

func (s *SuppliersServiceImpl) GetSupplierById(ctx context.Context, id int) (*model.Supplier, error) {
	supplier, err := s.Repo.GetSupplierById(ctx, id)
	if err != nil {
		return nil, err
	}

	if supplier == nil {
		return nil, errors.New("suppliers not found")
	}

	return supplier, nil
}

func (s *SuppliersServiceImpl) ChangeAddressParameter(ctx context.Context, id int, newAddressId int) error {
	if newAddressId < 0 {
		return errors.New("invalid new address id")
	}

	err := s.Repo.UpdateAddress(ctx, id, newAddressId)
	if err != nil {
		return err
	}

	return nil
}

func (s *SuppliersServiceImpl) DeleteSupplierById(ctx context.Context, id int) error {
	return s.Repo.DeleteSupplierById(ctx, id)
}
