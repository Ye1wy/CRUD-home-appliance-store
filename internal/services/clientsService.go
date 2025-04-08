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

type ClientsServiceInterface interface {
	CrudServiceInterface[model.Client, dto.ClientDTO]
	GetByNameAndSurname(ctx context.Context, name, surname string) ([]dto.ClientDTO, error)
	UpdateAddress(ctx context.Context, id, newAddressId string) error
}

type clientsServiceImpl struct {
	*CrudService[model.Client, dto.ClientDTO]
	repo repositories.ClientRepositoryInterface
}

func NewClientService(rep repositories.ClientRepositoryInterface, logger *logger.Logger) *clientsServiceImpl {
	crudService := NewCrudService(rep, mapper.ClientToDTO, mapper.ClientToModel, logger)
	logger.Debug("Client service is created")
	return &clientsServiceImpl{
		CrudService: crudService,
		repo:        rep,
	}
}

func (s *clientsServiceImpl) GetByNameAndSurname(ctx context.Context, name, surname string) ([]dto.ClientDTO, error) {
	op := "services.clientsService.GetByNameAndSurname"

	if name == "" || surname == "" {
		s.Logger.Debug("Name or Surname is empty", "op", op)
		return nil, fmt.Errorf("Client Service: client name and surname cannot be empty")
	}

	clients, err := s.repo.GetByNameAndSurname(ctx, name, surname)
	if clients == nil && err == nil {
		s.Logger.Debug("Client not found", "op", op)
		return nil, nil
	}

	if err != nil {
		s.Logger.Debug("Error recieved from repository", logger.Err(err), "op", op)
		return nil, fmt.Errorf("Client Service: Error when taking a client by first and last name: %v", err)
	}

	dtos := make([]dto.ClientDTO, len(clients))

	for i, item := range clients {
		dto, err := s.mapperD(&item)
		if err != nil {
			s.Logger.Debug("Mapping error", logger.Err(err), "op", op)
			return nil, fmt.Errorf("Client Service: mapping error %v", err)
		}

		dtos[i] = *dto
	}

	s.Logger.Debug("All clients is retrieved", "op", op)
	return dtos, nil
}

func (s *clientsServiceImpl) UpdateAddress(ctx context.Context, id, newAddressId string) error {
	op := "services.clientsService.UpdateAddress"

	if newAddressId == "" {
		s.Logger.Debug("New address is empty", "op", op)
		return fmt.Errorf("Service: Invalid new address id")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.Logger.Debug("Error converting id to ObjectID", logger.Err(err), "op", op)
		return fmt.Errorf("Service: Error change address parameter: %v", err)
	}

	objectAddressID, err := primitive.ObjectIDFromHex(newAddressId)
	if err != nil {
		s.Logger.Debug("Error converting newAddressId to ObjectID", logger.Err(err), "op", op)
		return fmt.Errorf("Service: Error change address parameter: %v", err)
	}

	err = s.repo.Update(ctx, objectID, objectAddressID)
	if err != nil {
		s.Logger.Debug("Error recieved from Update", logger.Err(err), "op", op)
		return fmt.Errorf("Service: Error change address parameter: %v", err)
	}

	s.Logger.Debug("Data updated", "op", op)
	return nil
}
