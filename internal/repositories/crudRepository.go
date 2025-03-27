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
	return &CrudRepository[T]{
		BaseMongoRepository: repo,
	}
}

func (r *CrudRepository[T]) Create(ctx context.Context, obj T) (*mongo.InsertOneResult, error) {
	return r.Collection.InsertOne(ctx, obj)
}

func (r *CrudRepository[T]) GetAll(ctx context.Context, limit, offset int) ([]T, error) {
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))

	cursor, err := r.Collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("Repository: Error find object: %v", err)
	}
	defer cursor.Close(ctx)

	var obj []T
	if err = cursor.All(ctx, &obj); err != nil {
		return nil, fmt.Errorf("Repository: Error binding data to object: %v", err)
	}

	return obj, nil
}

func (r *CrudRepository[T]) GetById(ctx context.Context, id primitive.ObjectID) (*T, error) {
	var obj T
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&obj)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, fmt.Errorf("Repository: Failed find one: %v", err)
	}

	return &obj, nil
}

func (r *CrudRepository[T]) Update(ctx context.Context, id primitive.ObjectID, item any) error {
	return nil
}

func (r *CrudRepository[T]) Delete(ctx context.Context, id primitive.ObjectID) error {
	op := "repositories.crudRepository.Delete"
	_, err := r.Collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("Repository: Failed to delete object by id: %v", err)
	}

	r.Logger.Debug("Crud Repository: Item deleted", "op", op)

	return nil
}
