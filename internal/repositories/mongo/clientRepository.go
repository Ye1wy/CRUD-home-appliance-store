package mongoRep

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
	GetByNameAndSurname(ctx context.Context, name string, surname string) ([]model.Client, error)
	UpdateAddress(ctx context.Context, id primitive.ObjectID, newAddressId primitive.ObjectID) error
}

type mongoClientRepository struct {
	*CrudRepository[model.Client]
}

func NewMongoClientRepository(db *mongo.Database, collection string, logger *logger.Logger) *mongoClientRepository {
	rep := NewCrudRepository[model.Client](db, collection, logger)
	logger.Debug("Client Repository is created (mongo)")
	return &mongoClientRepository{
		CrudRepository: rep,
	}
}

func (r *mongoClientRepository) GetByNameAndSurname(ctx context.Context, name string, surname string) ([]model.Client, error) {
	op := "repositories.clientRepository.GetByNameAndSurname"
	var clients []model.Client
	cursor, err := r.Collection.Find(ctx, bson.M{"client_name": name, "client_surname": surname})
	if err != nil {
		r.Logger.Debug("Failed on find", logger.Err(err), "op", op)
		return nil, fmt.Errorf("Client repository: Error find client: %v", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var client model.Client
		if err := cursor.Decode(&client); err != nil {
			r.Logger.Debug("Failed on decode cursor", logger.Err(err), "op", op)
			return nil, fmt.Errorf("Client repository: Error in decode cursor to struct: %v", err)
		}

		clients = append(clients, client)
	}

	if len(clients) == 0 {
		return nil, nil
	}

	r.Logger.Debug("All clients with name and surname is retrived", "op", op)
	return clients, nil
}

func (r *mongoClientRepository) UpdateAddress(ctx context.Context, id primitive.ObjectID, newAddressId primitive.ObjectID) error {
	op := "repositories.clientRepository.UpdateAddress"
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"address_id": newAddressId}})
	if err != nil {
		r.Logger.Debug("Failed on update address", logger.Err(err), "op", op)
		return fmt.Errorf("Client repository: Error update client: %v", err)
	}

	r.Logger.Debug("Updated succesufully", "op", op)
	return nil
}
