package repositories

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/logger"
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
	"context"
	"fmt"

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
	logger.Debug(" Product Repository is created (mongo)")
	return &mongoProductsRepository{
		CrudRepository: repo,
	}
}

func (r *mongoProductsRepository) DecreaseParameter(ctx context.Context, id string, decrease int) error {
	op := "repositories.productRepository.DecreaseParameter"
	result, err := r.Collection.UpdateOne(
		ctx, bson.M{
			"_id":             id,
			"available_stock": bson.M{"$gte": decrease},
		},
		bson.M{"$inc": bson.M{"available_stock": -decrease}})
	if err != nil {
		r.Logger.Debug("Failed update stock parameter", logger.Err(err), "op", op)
		return fmt.Errorf("Product Repository: failed update: %v", err)
	}

	if result.ModifiedCount == 0 {
		r.Logger.Debug("Insufficient stock or product not found", "op", op)
		return fmt.Errorf("Product Repository: Insufficient stock or product not found")
	}

	r.Logger.Debug("Updated successfully", "op", op)
	return nil
}
