package repositories

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/logger"

	"go.mongodb.org/mongo-driver/mongo"
)

type BaseRepositoryInterface interface {
}

type BaseMongoRepository struct {
	Collection *mongo.Collection
	Logger     *logger.Logger
}

func NewBaseRepository(collection *mongo.Collection, logger *logger.Logger) *BaseMongoRepository {
	logger.Debug("Base Repository is created")
	return &BaseMongoRepository{Collection: collection, Logger: logger}
}
