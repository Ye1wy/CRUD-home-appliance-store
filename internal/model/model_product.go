package model

type Product struct {
	Id             int64   `json:"id,omitempty" bson:"id,omitempty"`
	Name           string  `json:"name" bson:"name"`
	Category       string  `json:"category" bson:"category"`
	Price          float32 `json:"price" bson:"price"`
	AvailableStock int64   `json:"available_stock" bson:"available_stock"`
	LastUpdateDate string  `json:"last_update_date" bson:"last_update_date"`
	SupplierId     int64   `json:"supplier_id" bson:"supplier_id"`
	ImageId        string  `json:"image_id" bson:"image_id"`
}
