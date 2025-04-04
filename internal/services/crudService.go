package services

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/logger"
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

func NewCrudService[T any, D any](repository repositories.CrudRepositoryInterface[T], mapperD func(*T) (*D, error), mapperT func(*D) (*T, error), logger *logger.Logger) *CrudService[T, D] {
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
	obj, err := s.mapperT(dto)
	if err != nil {
		return nil, fmt.Errorf("Service: Error adding a client: %v", err)
	}

	_, err = s.repository.Create(ctx, *obj)
	if err != nil {
		return nil, fmt.Errorf("Service: Error adding a client: %v", err)
	}

	return dto, nil
}

func (s *CrudService[T, D]) GetAll(ctx context.Context, limit, offset int) ([]D, error) {
	if limit < 0 || offset < 0 {
		return nil, fmt.Errorf("Client service: limit and offset cannot be less of 0")
	}

	objs, err := s.repository.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("Client service: Error receriving all the client: %v", err)
	}

	dtos := make([]D, len(objs))

	for i, item := range objs {
		dto, err := s.mapperD(&item)
		if err != nil {
			return nil, nil
		}

		dtos[i] = *dto
	}

	return dtos, nil
}

func (s *CrudService[T, D]) GetById(ctx context.Context, id string) (*D, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("Service: Error converting a string to Hex: %v", err)
	}

	obj, err := s.repository.GetById(ctx, objectId)
	if obj == nil && err == nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	dto, err := s.mapperD(obj)
	if err != nil {
		return nil, fmt.Errorf("Service: Error mapping obj to dto: %v", err)
	}

	return dto, nil
}

func (s *CrudService[T, D]) Update(ctx context.Context, id string, item any) error {
	op := "services.crudService.Update"
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.Logger.Debug("Error converting a string to Hex", logger.Err(err), "op", op)
		return fmt.Errorf("Service Update: Error converting a string to Hex: %v", err)
	}

	if err := s.repository.Update(ctx, objectId, item); err != nil {
		s.Logger.Debug("Error from repository", logger.Err(err), "op", op)
		return fmt.Errorf("Service Update: Error updating the object: %v", err)
	}

	return nil
}

func (s *CrudService[T, D]) Delete(ctx context.Context, id string) error {
	op := "services.crudService.Delete"
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.Logger.Debug("Error converting a string to Hex", logger.Err(err), "op", op)
		return fmt.Errorf("Service: Error converting a string to Hex: %v", err)
	}

	if err := s.repository.Delete(ctx, objectId); err != nil {
		s.Logger.Debug("Error from repository", logger.Err(err), "op", op)
		return fmt.Errorf("Service: Error deleting the object: %v", err)
	}
	return nil
}
