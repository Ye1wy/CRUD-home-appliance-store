package repositories

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/database"
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClientRepository struct {
	Collection *mongo.Collection
}

func NewClientRepository(db *mongo.Database) *ClientRepository {
	return &ClientRepository{
		Collection: db.Collection(database.CLIENTS),
	}
}

func (r *ClientRepository) AddClient(ctx context.Context, client *model.Client) error {
	_, err := r.Collection.InsertOne(ctx, client)
	return err
}

func (r *ClientRepository) GetAllClients(ctx context.Context) ([]model.Client, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var client []model.Client
	if err = cursor.All(ctx, &client); err != nil {
		return nil, err
	}

	return client, nil
}

func (r *ClientRepository) GetClientByNameAndSurname(ctx context.Context, name string, surname string) (*model.Client, error) {
	var client model.Client
	err := r.Collection.FindOne(ctx, bson.M{"client_name": name, "client_surname": surname}).Decode(&client)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	return &client, nil
}

func (r *ClientRepository) UpdateAddress(ctx context.Context, id int, newAddressId int) error {
	_, err := r.Collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": bson.M{"address_id": newAddressId}})
	return err
}

func (r *ClientRepository) DeleteClientById(ctx context.Context, id int) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}
