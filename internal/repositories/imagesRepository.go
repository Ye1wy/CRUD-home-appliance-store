package repositories

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/database"
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ImagesRepository struct {
	Collection *mongo.Collection
}

func NewImagesRepository(db *mongo.Database) *ImagesRepository {
	return &ImagesRepository{
		Collection: db.Collection(database.IMAGES),
	}
}

func (r *ImagesRepository) AddImage(ctx context.Context, image *model.Image) error {
	_, err := r.Collection.InsertOne(ctx, image)
	return err
}

func (r *ImagesRepository) GetAllImages(ctx context.Context) ([]model.Image, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var images []model.Image
	if err = cursor.All(ctx, &images); err != nil {
		return nil, err
	}

	return images, nil
}

func (r *ImagesRepository) GetImageById(ctx context.Context, id int) (*model.Image, error) {
	var image model.Image
	err := r.Collection.FindOne(ctx, bson.M{"id": id}).Decode(&image)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	return &image, nil
}

func (r *ImagesRepository) ChangeImage(ctx context.Context, id int, newImage []byte) error {
	_, err := r.Collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": bson.M{"image": newImage}})
	return err
}

func (r *ImagesRepository) DeleteImageById(ctx context.Context, id int) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}
