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

type ClientRepositoryInterface interface {
	CrudRepositoryInterface[model.Client]
	GetClientByNameAndSurname(ctx context.Context, name string, surname string) ([]model.Client, error)
	UpdateAddress(ctx context.Context, id primitive.ObjectID, newAddressId primitive.ObjectID) error
}

type MongoClientRepository struct {
	*CrudRepository[model.Client]
}

func NewMongoClientRepository(db *mongo.Database, collection string, logger *logger.Logger) *MongoClientRepository {
	rep := NewCrudRepository[model.Client](db, collection, logger)
	return &MongoClientRepository{
		CrudRepository: rep,
	}
}

func (r *MongoClientRepository) GetClientByNameAndSurname(ctx context.Context, name string, surname string) ([]model.Client, error) {
	op := "repositories.clientRepository.GetClientByNameAndSurname"

	var clients []model.Client
	cursor, err := r.Collection.Find(ctx, bson.M{"client_name": name, "client_surname": surname})
	if err != nil {
		return nil, fmt.Errorf("Client repository: Error find client: %v", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var client model.Client
		if err := cursor.Decode(&client); err != nil {
			return nil, fmt.Errorf("Client repository: Error in decode cursor to struct: %v", err)
		}

		clients = append(clients, client)
	}

	if len(clients) == 0 {
		return nil, nil
	}

	r.Logger.Debug("ClientRepository: all clients with name and surname is retrived", "op", op)
	return clients, nil
}

func (r *MongoClientRepository) UpdateAddress(ctx context.Context, id primitive.ObjectID, newAddressId primitive.ObjectID) error {
	op := "repositories.clientRepository.UpdateAddress"
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"address_id": newAddressId}})
	if err != nil {
		return fmt.Errorf("Client repository: Error update client: %v", err)
	}

	r.Logger.Debug("ClientRepository: Updated succesufully", "op", op)
	return nil
}
