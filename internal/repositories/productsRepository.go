package repositories

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/database"
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepository interface {
	AddProduct(ctx context.Context, product *model.Product) (*mongo.InsertOneResult, error)
	GetAllProducts(ctx context.Context, limit, offset int) ([]model.Product, error)
	GetProductById(ctx context.Context, id string) (*model.Product, error)
	DecreaseParametr(ctx context.Context, id string, decrease int) error
	DeleteProductById(ctx context.Context, id string) error
}

type mongoProductsRepository struct {
	Collection *mongo.Collection
}

func NewMongoProductsRepository(db *mongo.Database) *mongoProductsRepository {
	return &mongoProductsRepository{
		Collection: db.Collection(database.PRODUCTS),
	}
}

func (r *mongoProductsRepository) AddProduct(ctx context.Context, product *model.Product) (*mongo.InsertOneResult, error) {
	product.LastUpdateDate = time.Now()
	return r.Collection.InsertOne(ctx, product)
}

func (r *mongoProductsRepository) GetAllProducts(ctx context.Context, limit, offset int) ([]model.Product, error) {
	findOption := options.Find()
	findOption.SetLimit(int64(limit))
	findOption.SetSkip(int64(offset))

	cursor, err := r.Collection.Find(ctx, bson.M{}, findOption)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var products []model.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *mongoProductsRepository) GetProductById(ctx context.Context, id string) (*model.Product, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var product model.Product
	err = r.Collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	return &product, nil
}

func (r *mongoProductsRepository) DecreaseParametr(ctx context.Context, id string, decrease int) error {
	result, err := r.Collection.UpdateOne(
		ctx, bson.M{
			"_id":             id,
			"available_stock": bson.M{"$gte": decrease},
		},
		bson.M{"$inc": bson.M{"available_stock": -decrease}})
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("insufficient stock or product not found")
	}

	return nil
}

func (r *mongoProductsRepository) DeleteProductById(ctx context.Context, id string) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
