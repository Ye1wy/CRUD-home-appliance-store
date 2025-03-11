package repositories

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/database"
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ClientRepository interface {
	AddClient(ctx context.Context, client *model.Client) (*mongo.InsertOneResult, error)
	GetAllClients(ctx context.Context, limit, offset int) ([]model.Client, error)
	GetClientByNameAndSurname(ctx context.Context, name string, surname string) ([]model.Client, error)
	UpdateAddress(ctx context.Context, id primitive.ObjectID, newAddressId primitive.ObjectID) error
	DeleteClientById(ctx context.Context, id primitive.ObjectID) error
}

type mongoClientRepository struct {
	collection *mongo.Collection
}

func NewMongoClientRepository(db *mongo.Database) *mongoClientRepository {
	return &mongoClientRepository{
		collection: db.Collection(database.CLIENTS),
	}
}

func (r *mongoClientRepository) AddClient(ctx context.Context, client *model.Client) (*mongo.InsertOneResult, error) {
	return r.collection.InsertOne(ctx, client)
}

func (r *mongoClientRepository) GetAllClients(ctx context.Context, limit, offset int) ([]model.Client, error) {
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
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

func (r *mongoClientRepository) GetClientByNameAndSurname(ctx context.Context, name string, surname string) ([]model.Client, error) {
	var clients []model.Client
	cursor, err := r.collection.Find(ctx, bson.M{"client_name": name, "client_surname": surname})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var client model.Client
		if err := cursor.Decode(&client); err != nil {
			return nil, err
		}

		clients = append(clients, client)
	}

	if len(clients) == 0 {
		return nil, nil
	}

	return clients, nil
}

func (r *mongoClientRepository) UpdateAddress(ctx context.Context, id primitive.ObjectID, newAddressId primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"address_id": newAddressId}})
	return err
}

func (r *mongoClientRepository) DeleteClientById(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
