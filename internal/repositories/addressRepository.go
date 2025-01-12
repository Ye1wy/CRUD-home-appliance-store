package repositories

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/database"
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AddressRepository struct {
	Collection *mongo.Collection
}

func NewAddressRepository(db *mongo.Database) *AddressRepository {
	return &AddressRepository{
		Collection: db.Collection(database.ADDRESSES),
	}
}

func (r *AddressRepository) AddAddress(ctx context.Context, address *model.Address) error {
	_, err := r.Collection.InsertOne(ctx, address)
	return err
}

func (r *AddressRepository) GetAllAddresses(ctx context.Context) ([]model.Address, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var addresses []model.Address
	if err = cursor.All(ctx, &addresses); err != nil {
		return nil, err
	}

	return addresses, nil
}

func (r *AddressRepository) GetAddressById(ctx context.Context, id int) (*model.Address, error) {
	var address model.Address
	err := r.Collection.FindOne(ctx, bson.M{"id": id}).Decode(&address)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	return &address, nil
}

func (r *AddressRepository) DeleteImageById(ctx context.Context, id int) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}
