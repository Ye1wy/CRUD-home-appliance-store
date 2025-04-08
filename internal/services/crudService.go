package services

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/dto"
	"CRUD-HOME-APPLIANCE-STORE/internal/logger"
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
	"CRUD-HOME-APPLIANCE-STORE/internal/repositories"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CrudServiceInterface[T any, D any] interface {
	Create(ctx context.Context, dto *D) (*D, error)
	GetAll(ctx context.Context, limit, offset int) ([]D, error)
	GetById(ctx context.Context, id string) (*D, error)
	Update(ctx context.Context, id string, item any) error
	Delete(ctx context.Context, id string) error
}

type CrudService[T any, D any] struct {
	*BaseService
	repository repositories.CrudRepositoryInterface[T]
	mapperD    func(*T) (*D, error)
	mapperT    func(*D) (*T, error)
}

type ClientCrudService struct {
	CrudService[model.Client, dto.ClientDTO]
}

func NewCrudService[T any, D any](
	repository repositories.CrudRepositoryInterface[T],
	mapperD func(*T) (*D, error),
	mapperT func(*D) (*T, error),
	logger *logger.Logger,
) *CrudService[T, D] {
	service := NewBaseService(nil, logger)
	logger.Debug("Crud Service is created")
	return &CrudService[T, D]{
		BaseService: service,
		repository:  repository,
		mapperD:     mapperD,
		mapperT:     mapperT,
	}
}

func (s *CrudService[T, D]) Create(ctx context.Context, dto *D) (*D, error) {
	op := "services.crudService.Create"
	obj, err := s.mapperT(dto)
	if err != nil {
		s.Logger.Debug("Mapping error", logger.Err(err), "op", op)
		return nil, fmt.Errorf("CRUD Service: Error adding a object: %v", err)
	}

	_, err = s.repository.Create(ctx, *obj)
	if err != nil {
		s.Logger.Debug("Create error", logger.Err(err), "op", op)
		return nil, fmt.Errorf("CRUD Service: Error adding a object: %v", err)
	}

	s.Logger.Debug("Create is complete", "op", op)
	return dto, nil
}

func (s *CrudService[T, D]) GetAll(ctx context.Context, limit, offset int) ([]D, error) {
	op := "services.crudService.GetAll"

	if limit < 0 || offset < 0 {
		s.Logger.Debug("Limit and offset is less then 0", "op", op)
		return nil, fmt.Errorf("CRUD Service: Limit and offset cannot be less of 0")
	}

	objs, err := s.repository.GetAll(ctx, limit, offset)
	if err != nil {
		s.Logger.Debug("Failed in GetAll method: Error from repository", logger.Err(err), "op", op)
		return nil, fmt.Errorf("CRUD Service: Error recieved from GetAll: %v", err)
	}

	dtos := make([]D, len(objs))

	for i, item := range objs {
		dto, err := s.mapperD(&item)
		if err != nil {
			s.Logger.Debug("Mapping Error", logger.Err(err), "op", op)
			return nil, fmt.Errorf("CRUD Service: mapping error %v", err)
		}

		dtos[i] = *dto
	}

	s.Logger.Debug("All data is retrieved", "op", op)
	return dtos, nil
}

func (s *CrudService[T, D]) GetById(ctx context.Context, id string) (*D, error) {
	op := "services.crudService.GetById"
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.Logger.Debug("Error converting string to ObjectID", logger.Err(err), "op", op)
		return nil, fmt.Errorf("CRUD Service: Error converting a string to Hex: %v", err)
	}

	s.Logger.Debug("Converted id", "objectID", objectId, "op", op)
	obj, err := s.repository.GetById(ctx, objectId)
	if obj == nil && err == nil {
		s.Logger.Debug("Object not found", "op", op)
		return nil, nil
	}

	if err != nil {
		s.Logger.Debug("Error in method GetById", logger.Err(err), "op", op)
		return nil, fmt.Errorf("CRUD Service: Error recieved from GetById %v", err)
	}

	dto, err := s.mapperD(obj)
	if err != nil {
		s.Logger.Debug("Mapping error", logger.Err(err), "op", op)
		return nil, fmt.Errorf("CRUD Service: Error mapping obj to dto: %v", err)
	}

	s.Logger.Debug("Object is retrieved")
	return dto, nil
}

func (s *CrudService[T, D]) Update(ctx context.Context, id string, item any) error {
	op := "services.crudService.Update"
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.Logger.Debug("Error converting a string to Hex", logger.Err(err), "op", op)
		return fmt.Errorf("CRUD Service: Error converting a string to Hex: %v", err)
	}

	if err := s.repository.Update(ctx, objectId, item); err != nil {
		s.Logger.Debug("Error recieved from repository Update", logger.Err(err), "op", op)
		return fmt.Errorf("CRUD Service: Error updating the object: %v", err)
	}

	s.Logger.Debug("Data updated", "op", op)
	return nil
}

func (s *CrudService[T, D]) Delete(ctx context.Context, id string) error {
	op := "services.crudService.Delete"
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.Logger.Debug("Error converting a string to Hex", logger.Err(err), "op", op)
		return fmt.Errorf("CRUD Service: Error converting a string to Hex: %v", err)
	}

	if err := s.repository.Delete(ctx, objectId); err != nil {
		s.Logger.Debug("Error recieved from repository", logger.Err(err), "op", op)
		return fmt.Errorf("CRUD Service: Error deleting the object: %v", err)
	}

	s.Logger.Debug("Data is deleted", "op", op)
	return nil
}
