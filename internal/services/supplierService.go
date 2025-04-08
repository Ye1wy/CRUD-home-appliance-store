package services

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/dto"
	"CRUD-HOME-APPLIANCE-STORE/internal/logger"
	"CRUD-HOME-APPLIANCE-STORE/internal/mapper"
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
	"CRUD-HOME-APPLIANCE-STORE/internal/repositories"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SupplierServiceInterface interface {
	CrudServiceInterface[model.Supplier, dto.SupplierDTO]
	UpdateAddress(ctx context.Context, id, newAddressId string) error
}

type supplierServiceImpl struct {
	*CrudService[model.Supplier, dto.SupplierDTO]
	repo repositories.SupplierRepositoryInterface
}

func NewSupplierService(repo repositories.SupplierRepositoryInterface, logger *logger.Logger) *supplierServiceImpl {
	service := NewCrudService(repo, mapper.SupplierToDTO, mapper.SupplierToModel, logger)
	logger.Debug("Supplier serivce created")
	return &supplierServiceImpl{service, repo}
}

func (s *supplierServiceImpl) UpdateAddress(ctx context.Context, id, newAddressId string) error {
	op := "services.supplierService.UpdateAddress"

	if newAddressId == "" {
		s.Logger.Debug("New address is empty", "op", op)
		return fmt.Errorf("Supplier Service: Invalid new address id")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.Logger.Debug("Error converting id to ObjectID", logger.Err(err), "op", op)
		return fmt.Errorf("Supplier Service: Error change address parameter: %v", err)
	}

	objectAddressID, err := primitive.ObjectIDFromHex(newAddressId)
	if err != nil {
		s.Logger.Debug("Failed converting newAddressId to ObjectID", logger.Err(err), "op", op)
		return fmt.Errorf("Supplier Service: Error change address parameter: %v", err)
	}

	err = s.repo.Update(ctx, objectID, objectAddressID)
	if err != nil {
		s.Logger.Debug("Error recieved from repository Update", logger.Err(err), "op", op)
		return fmt.Errorf("Supplier Service: Error change address parameter: %v", err)
	}

	s.Logger.Debug("Supplier updated", "op", op)
	return nil
}
