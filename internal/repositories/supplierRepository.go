package repositories

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/logger"
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SupplierRepositoryInterface interface {
	CrudRepositoryInterface[model.Supplier]
	UpdateAddress(ctx context.Context, id primitive.ObjectID, newAddressId primitive.ObjectID) error
}

type mongoSupplierRepository struct {
	*CrudRepository[model.Supplier]
}

func NewMongoSupplierRepository(db *mongo.Database, collection string, logger *logger.Logger) *mongoSupplierRepository {
	repo := NewCrudRepository[model.Supplier](db, collection, logger)
	logger.Debug("Supplier Repository is created (mongo)")
	return &mongoSupplierRepository{
		CrudRepository: repo,
	}
}

func (r *mongoSupplierRepository) UpdateAddress(ctx context.Context, id primitive.ObjectID, newAddressId primitive.ObjectID) error {
	op := "repositorires.supplierRepository.UpdateAddress"
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"address_id": newAddressId}})
	if err != nil {
		r.Logger.Debug("Failed on update address", logger.Err(err), "op", op)
		return fmt.Errorf("Supplier Repository: Error update supplier: %v", err)
	}

	r.Logger.Debug("Updated successfully", "op", op)
	return nil
}
