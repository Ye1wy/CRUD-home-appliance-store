package repositories

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/database"
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductsRepository struct {
	Collection *mongo.Collection
}

func NewProductsRepository(db *mongo.Database) *ProductsRepository {
	return &ProductsRepository{
		Collection: db.Collection(database.PRODUCTS),
	}
}

func (r *ProductsRepository) AddProduct(ctx context.Context, product *model.Product) error {
	_, err := r.Collection.InsertOne(ctx, product)
	return err
}

func (r *ProductsRepository) GetAllProducts(ctx context.Context) ([]model.Product, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{})
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

func (r *ProductsRepository) GetProductById(ctx context.Context, id int) (*model.Product, error) {
	var product model.Product
	err := r.Collection.FindOne(ctx, bson.M{"id": id}).Decode(&product)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	return &product, nil
}

func (r *ProductsRepository) DecreaseParametr(ctx context.Context, id int, decrease int) error {
	var product model.Product
	err := r.Collection.FindOne(ctx, bson.M{"id": id}).Decode(&product)
	if err == mongo.ErrNoDocuments {
		return fmt.Errorf("product with id %d not found", id)

	} else if err != nil {
		return err
	}

	if product.AvailableStock < int64(decrease) {
		return fmt.Errorf("insufficient stock to decrease by %d", decrease)
	}

	_, err = r.Collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$inc": bson.M{"available_stock": -decrease}})
	return err
}

func (r *ProductsRepository) DeleteProductById(ctx context.Context, id int) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}
