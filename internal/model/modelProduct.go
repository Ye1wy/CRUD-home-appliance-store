package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name           string             `json:"name" bson:"name"`
	Category       string             `json:"category" bson:"category"`
	Price          float32            `json:"price" bson:"price"`
	AvailableStock int64              `json:"available_stock" bson:"available_stock"`
	LastUpdateDate time.Time          `json:"last_update_date" bson:"last_update_date"`
	SupplierId     string             `json:"supplier_id" bson:"supplier_id"`
	ImageId        string             `json:"image_id" bson:"image_id"`
}
