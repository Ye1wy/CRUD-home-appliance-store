package mongodb

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/database"
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func Connect(mongoURL string, log *slog.Logger) (*MongoStorage, error) {
	const op = "database.mongodb.Connect"

	clientOption := options.Client().ApplyURI(mongoURL)
	client, err := mongo.Connect(context.Background(), clientOption)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	db := client.Database(database.DATABASE)

	err = CreateCollections(client.Database(database.DATABASE))
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return &MongoStorage{Client: client, Database: db}, nil
}

func CreateCollections(db *mongo.Database) error {
	const op = "database.mongodb.CreateCollections"

	collectionToCreate := []string{database.CLIENTS, database.PRODUCTS, database.SUPPLIERS, database.IMAGES, database.ADDRESSES}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	existingCollection, err := db.ListCollectionNames(ctx, struct{}{})
	if err != nil {
		return err
	}

	for _, collectionName := range collectionToCreate {
		if contains(existingCollection, collectionName) {
			continue
		}

		err := db.CreateCollection(ctx, collectionName)
		if err != nil {
			return fmt.Errorf("%s: %s: %v", op, collectionName, err)
		}
	}

	return nil
}

func contains(slice []string, item string) bool {
	for _, elem := range slice {
		if elem == item {
			return true
		}
	}

	return false
}
