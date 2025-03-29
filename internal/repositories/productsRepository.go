package repositories

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/logger"
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepositoryInterface interface {
	CrudRepositoryInterface[model.Product]
	DecreaseParameter(ctx context.Context, id string, decrease int) error
}

type mongoProductsRepository struct {
	*CrudRepository[model.Product]
}

func NewMongoProductsRepository(db *mongo.Database, collection string, logger *logger.Logger) *mongoProductsRepository {
	repo := NewCrudRepository[model.Product](db, collection, logger)
	return &mongoProductsRepository{
		CrudRepository: repo,
	}
}

func (r *mongoProductsRepository) DecreaseParameter(ctx context.Context, id string, decrease int) error {
	result, err := r.Collection.UpdateOne(
		ctx, bson.M{
			"_id":             id,
			"available_stock": bson.M{"$gte": decrease},
		},
		bson.M{"$inc": bson.M{"available_stock": -decrease}})
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("insufficient stock or product not found")
	}

	return nil
}
