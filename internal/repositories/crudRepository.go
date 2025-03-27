package repositories

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/logger"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CrudRepositoryInterface[T any] interface {
	Create(ctx context.Context, obj T) (*mongo.InsertOneResult, error)
	GetAll(ctx context.Context, limit, offset int) ([]T, error)
	GetById(ctx context.Context, id primitive.ObjectID) (*T, error)
	Update(ctx context.Context, id primitive.ObjectID, item any) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type CrudRepository[T any] struct {
	*BaseMongoRepository
}

func NewCrudRepository[T any](db *mongo.Database, collection string, logger *logger.Logger) *CrudRepository[T] {
	repo := NewBaseRepository(db.Collection(collection), logger)
	logger.Info("Crud Repository is created")
	return &CrudRepository[T]{
		BaseMongoRepository: repo,
	}
}

func (r *CrudRepository[T]) Create(ctx context.Context, obj T) (*mongo.InsertOneResult, error) {
	return r.Collection.InsertOne(ctx, obj)
}

func (r *CrudRepository[T]) GetAll(ctx context.Context, limit, offset int) ([]T, error) {
	op := "repositories.crudRepository.GetAll"
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))

	cursor, err := r.Collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		r.Logger.Debug("Dead on Find item", logger.Err(err), "op", op)
		return nil, fmt.Errorf("CrudRepository: Error find object: %v", err)
	}
	defer cursor.Close(ctx)

	var obj []T
	if err = cursor.All(ctx, &obj); err != nil {
		r.Logger.Debug("Dead on interate cursor for all data", logger.Err(err), "op", op)
		return nil, fmt.Errorf("Error binding data to object: %v", err)
	}

	r.Logger.Debug("All data is retrieved", "op", op)

	return obj, nil
}

func (r *CrudRepository[T]) GetById(ctx context.Context, id primitive.ObjectID) (*T, error) {
	op := "repository.crudRepository.GetById"
	var obj T
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&obj)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.Logger.Debug("No one document is not found", "op", op)
			return nil, nil
		}

		r.Logger.Debug("Failed on find one item", logger.Err(err), "op", op)
		return nil, fmt.Errorf("Repository: Failed find one: %v", err)
	}

	r.Logger.Debug("Data is retrieved", "op", op)
	return &obj, nil
}

func (r *CrudRepository[T]) Update(ctx context.Context, id primitive.ObjectID, item any) error {
	return nil
}

func (r *CrudRepository[T]) Delete(ctx context.Context, id primitive.ObjectID) error {
	op := "repositories.crudRepository.Delete"
	_, err := r.Collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		r.Logger.Debug("Failed on DeleteOne", logger.Err(err), "op", op)
		return fmt.Errorf("Repository: Failed to delete object by id: %v", err)
	}

	r.Logger.Debug("Crud Repository: Item deleted", "op", op)
	return nil
}
