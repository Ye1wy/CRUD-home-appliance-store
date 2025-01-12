package repositories

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/database"
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SuppliersRepository struct {
	Collection *mongo.Collection
}

func NewSuppliersRepository(db *mongo.Database) *SuppliersRepository {
	return &SuppliersRepository{
		Collection: db.Collection(database.SUPPLIERS),
	}
}

func (r *SuppliersRepository) AddSupplier(ctx context.Context, supplier *model.Supplier) error {
	_, err := r.Collection.InsertOne(ctx, supplier)
	return err
}

func (r *SuppliersRepository) GetAllSuppliers(ctx context.Context, limit, offset int) ([]model.Supplier, error) {
	findOption := options.Find()
	findOption.SetLimit(int64(limit))
	findOption.SetSkip(int64(offset))

	cursor, err := r.Collection.Find(ctx, bson.M{}, findOption)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var supplier []model.Supplier
	if err = cursor.All(ctx, &supplier); err != nil {
		return nil, err
	}

	return supplier, nil
}

func (r *SuppliersRepository) GetSupplierById(ctx context.Context, id int) (*model.Supplier, error) {
	var supplier model.Supplier
	err := r.Collection.FindOne(ctx, bson.M{"id": id}).Decode(&supplier)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	return &supplier, nil
}

func (r *SuppliersRepository) UpdateAddress(ctx context.Context, id int, newAddressId int) error {
	_, err := r.Collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": bson.M{"address_id": newAddressId}})
	return err
}

func (r *SuppliersRepository) DeleteSupplierById(ctx context.Context, id int) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}
