package repositories

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/database"
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ImagesRepository struct {
	Collection    *mongo.Collection
	SecCollection *mongo.Collection
}

func NewImagesRepository(db *mongo.Database) *ImagesRepository {
	return &ImagesRepository{
		Collection:    db.Collection(database.IMAGES),
		SecCollection: db.Collection(database.PRODUCTS),
	}
}

func (r *ImagesRepository) AddImage(ctx context.Context, image *model.Image) error {
	_, err := r.Collection.InsertOne(ctx, image)
	return err
}

func (r *ImagesRepository) GetProductImages(ctx context.Context, id int) (*model.ProductImage, error) {
	var image model.ProductImage
	err := r.SecCollection.FindOne(ctx, bson.M{"id": id}).Decode(&image)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	return &image, nil
}

func (r *ImagesRepository) GetImageById(ctx context.Context, id string) (*model.Image, error) {
	var image model.Image
	err := r.Collection.FindOne(ctx, bson.M{"id": id}).Decode(&image)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	return &image, nil
}

func (r *ImagesRepository) ChangeImage(ctx context.Context, id string, newImage []byte) error {
	_, err := r.Collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": bson.M{"image": newImage}})
	return err
}

func (r *ImagesRepository) DeleteImageById(ctx context.Context, id string) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}
